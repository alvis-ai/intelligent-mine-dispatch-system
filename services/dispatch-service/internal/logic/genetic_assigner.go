package logic

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/pkg/utils"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/model"
	"github.com/aicong/mine-dispatch/services/dispatch-service/internal/svc"
)

// ── GA Parameters ──

type GAParams struct {
	PopulationSize int
	Generations    int
	MutationRate   float64
	ElitismCount   int
	TournamentSize int
}

func DefaultGAParams() GAParams {
	return GAParams{
		PopulationSize: 50,
		Generations:    100,
		MutationRate:   0.1,
		ElitismCount:   2,
		TournamentSize: 3,
	}
}

// ── Core types ──

// VehicleInfo holds vehicle state for GA optimization.
type VehicleInfo struct {
	ID          uint64
	Plate       string
	Latitude    float64
	Longitude   float64
	ActiveTasks int // number of currently active/in-progress tasks
}

// TaskInfo holds a dispatch task candidate for GA.
type TaskInfo struct {
	Index       int // index into the task list
	LoadPointID uint64
	DumpPointID uint64
	LoadLat     float64
	LoadLon     float64
	DumpLat     float64
	DumpLon     float64
}

// Chromosome encodes a vehicle→task assignment.
// Genes[i] = vehicle index (into VehicleInfo slice) assigned to TaskInfo i.
type Chromosome struct {
	Genes   []int
	Fitness float64
}

// DistanceCache caches road distances between coordinate pairs to avoid
// redundant route-service calls during GA evaluation.
type DistanceCache struct {
	mu   sync.RWMutex
	data map[distKey]distValue
}

type distKey struct {
	fromLat, fromLon, toLat, toLon float64
}

type distValue struct {
	distanceM float64
	durationS float64
}

func NewDistanceCache() *DistanceCache {
	return &DistanceCache{data: make(map[distKey]distValue)}
}

func (dc *DistanceCache) Get(fromLat, fromLon, toLat, toLon float64) (float64, float64, bool) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	v, ok := dc.data[distKey{fromLat, fromLon, toLat, toLon}]
	return v.distanceM, v.durationS, ok
}

func (dc *DistanceCache) Set(fromLat, fromLon, toLat, toLon, distanceM, durationS float64) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.data[distKey{fromLat, fromLon, toLat, toLon}] = distValue{distanceM, durationS}
}

// ── LoadingPoint model (mirrors DB table) ──

type LoadingPoint struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement:false"`
	Name      string  `gorm:"size:128"`
	Type      string  `gorm:"size:32"`
	Latitude  float64 `gorm:"default:0"`
	Longitude float64 `gorm:"default:0"`
	Material  string  `gorm:"size:64"`
	Status    int     `gorm:"default:1"`
	MineID    uint64  `gorm:"index"`
}

func (LoadingPoint) TableName() string { return "loading_points" }

// ── GeneticAlgorithmAssigner ──

type GeneticAlgorithmAssigner struct {
	svc    *svc.ServiceContext
	ctx    context.Context
	params GAParams
	cache  *DistanceCache
}

func NewGeneticAlgorithmAssigner(ctx context.Context, svc *svc.ServiceContext) *GeneticAlgorithmAssigner {
	return &GeneticAlgorithmAssigner{
		svc:    svc,
		ctx:    ctx,
		params: DefaultGAParams(),
		cache:  NewDistanceCache(),
	}
}

// Assign implements TaskAssigner for a single task.
// It uses route-service for road distance and considers existing vehicle load.
func (a *GeneticAlgorithmAssigner) Assign(vehicleID, loadPointID, dumpPointID uint64) (uint64, error) {
	// Load coordinates for the load/dump points
	loadLat, loadLon := a.getPointCoords(loadPointID)
	dumpLat, dumpLon := a.getPointCoords(dumpPointID)

	// Get vehicle's position and existing active task load
	var activeCount int64
	a.svc.DB.Model(&model.DispatchTask{}).Where("vehicle_id = ? AND status IN ?", vehicleID, []string{model.StatusActive, model.StatusPending}).Count(&activeCount)

	var vhc struct{ Latitude, Longitude float64 }
	a.svc.DB.Table("vehicles").Where("id = ?", vehicleID).Select("latitude, longitude").Scan(&vhc)
	vinf := VehicleInfo{
		ID:          vehicleID,
		Latitude:    vhc.Latitude,
		Longitude:   vhc.Longitude,
		ActiveTasks: int(activeCount),
	}

	// Get road distances
	emptyDist, emptyDur := a.getRoadDistance(vinf.Latitude, vinf.Longitude, loadLat, loadLon)
	loadDist, loadDur := a.getRoadDistance(loadLat, loadLon, dumpLat, dumpLon)

	totalDist := emptyDist + loadDist
	totalDur := emptyDur + loadDur

	// Load penalty: bias against vehicles already carrying many tasks
	loadPenalty := float64(vinf.ActiveTasks) * 300.0 // 300m equivalent per extra task
	score := totalDist + loadPenalty

	_ = score // Track score for future optimizations

	// Create the task
	task := model.DispatchTask{
		ID:          utils.NextID(),
		VehicleID:   vehicleID,
		LoadPointID: loadPointID,
		DumpPointID: dumpPointID,
		LoadLat:     loadLat,
		LoadLon:     loadLon,
		DumpLat:     dumpLat,
		DumpLon:     dumpLon,
		Status:      model.StatusActive,
		Algorithm:   "genetic_algorithm",
	}

	// Check if there are pending tasks that could benefit from batch optimization
	var pendingCount int64
	a.svc.DB.Model(&model.DispatchTask{}).Where("status = ?", model.StatusPending).Count(&pendingCount)
	if pendingCount > 0 {
		// Run batch optimization with pending tasks
		go a.asyncBatchOptimize()
	}

	if err := a.svc.DB.Create(&task).Error; err != nil {
		return 0, err
	}

	msg := fmt.Sprintf(`{"task_id":%d,"vehicle_id":%d,"action":"assign","distance_m":%.0f,"duration_s":%.0f}`,
		task.ID, vehicleID, totalDist, totalDur)
	a.svc.Redis.Publish(a.ctx, "dispatch:events", msg)
	return task.ID, nil
}

// BatchOptimize runs GA to find optimal assignments for pending tasks.
// Returns the best chromosome found.
func (a *GeneticAlgorithmAssigner) BatchOptimize(vehicles []VehicleInfo, tasks []TaskInfo) *Chromosome {
	if len(tasks) == 0 || len(vehicles) == 0 {
		return nil
	}

	pop := a.initializePopulation(vehicles, tasks)
	a.evaluatePopulation(pop, vehicles, tasks)

	for gen := 0; gen < a.params.Generations; gen++ {
		pop = a.evolve(pop, vehicles, tasks)
	}

	// Return the best chromosome
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fitness > pop[j].Fitness
	})
	return &pop[0]
}

// asyncBatchOptimize runs batch optimization asynchronously for pending tasks.
func (a *GeneticAlgorithmAssigner) asyncBatchOptimize() {
	// Load pending tasks
	var pendingTasks []model.DispatchTask
	a.svc.DB.Where("status = ?", model.StatusPending).Find(&pendingTasks)
	if len(pendingTasks) == 0 {
		return
	}

	// Load all vehicles
	var vhcList []struct {
		ID        uint64
		Plate     string
		Latitude  float64
		Longitude float64
	}
	a.svc.DB.Table("vehicles").Select("id, plate, latitude, longitude").Find(&vhcList)
	if len(vhcList) == 0 {
		return
	}

	vehicles := make([]VehicleInfo, len(vhcList))
	taskCounts := make(map[uint64]int64)
	for _, t := range pendingTasks {
		taskCounts[t.VehicleID]++
	}
	for i, v := range vhcList {
		vehicles[i] = VehicleInfo{
			ID:          v.ID,
			Plate:       v.Plate,
			Latitude:    v.Latitude,
			Longitude:   v.Longitude,
			ActiveTasks: int(taskCounts[v.ID]),
		}
	}

	tasks := make([]TaskInfo, len(pendingTasks))
	for i, t := range pendingTasks {
		tasks[i] = TaskInfo{
			Index:       i,
			LoadPointID: t.LoadPointID,
			DumpPointID: t.DumpPointID,
			LoadLat:     t.LoadLat,
			LoadLon:     t.LoadLon,
			DumpLat:     t.DumpLat,
			DumpLon:     t.DumpLon,
		}
	}

	result := a.BatchOptimize(vehicles, tasks)
	if result == nil {
		return
	}

	// Apply the best assignments: update tasks with optimal vehicle assignments
	for i, task := range pendingTasks {
		vi := result.Genes[i]
		if vi >= 0 && vi < len(vehicles) {
			a.svc.DB.Model(&model.DispatchTask{}).Where("id = ?", task.ID).Update("vehicle_id", vehicles[vi].ID)
		}
	}
}

// ── Population initialization ──

func (a *GeneticAlgorithmAssigner) initializePopulation(vehicles []VehicleInfo, tasks []TaskInfo) []Chromosome {
	pop := make([]Chromosome, a.params.PopulationSize)
	nTasks := len(tasks)
	nVeh := len(vehicles)

	// Heuristic: assign each task to the nearest available vehicle
	heuristicGenes := make([]int, nTasks)
	for i, task := range tasks {
		bestVeh := 0
		bestDist := math.MaxFloat64
		for vi, v := range vehicles {
			d := haversineGA(v.Latitude, v.Longitude, task.LoadLat, task.LoadLon)
			loadPenalty := float64(v.ActiveTasks) * 1000
			if d+loadPenalty < bestDist {
				bestDist = d + loadPenalty
				bestVeh = vi
			}
		}
		heuristicGenes[i] = bestVeh
	}

	for i := range pop {
		genes := make([]int, nTasks)
		if i == 0 {
			// Best chromosome: use heuristic
			copy(genes, heuristicGenes)
		} else {
			// Random chromosomes with heuristic bias
			for j := range genes {
				if rand.Float64() < 0.3 {
					genes[j] = heuristicGenes[j]
				} else {
					genes[j] = rand.Intn(nVeh)
				}
			}
		}
		pop[i] = Chromosome{Genes: genes}
	}
	return pop
}

// ── Evaluation ──

func (a *GeneticAlgorithmAssigner) evaluatePopulation(pop []Chromosome, vehicles []VehicleInfo, tasks []TaskInfo) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // limit concurrency

	for i := range pop {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			pop[idx].Fitness = a.evaluate(&pop[idx], vehicles, tasks)
		}(i)
	}
	wg.Wait()
}

func (a *GeneticAlgorithmAssigner) evaluate(ch *Chromosome, vehicles []VehicleInfo, tasks []TaskInfo) float64 {
	totalDist := 0.0
	vehicleTaskCount := make(map[int]int)

	for ti, vi := range ch.Genes {
		if vi < 0 || vi >= len(vehicles) {
			return math.Inf(-1)
		}
		v := vehicles[vi]
		t := tasks[ti]

		// Empty travel: vehicle → load point
		emptyDist, _ := a.getRoadDistanceCached(v.Latitude, v.Longitude, t.LoadLat, t.LoadLon)
		// Loaded travel: load point → dump point
		loadDist, _ := a.getRoadDistanceCached(t.LoadLat, t.LoadLon, t.DumpLat, t.DumpLon)

		totalDist += emptyDist + loadDist
		vehicleTaskCount[vi]++
	}

	// Load balance penalty: penalize stddev of tasks per vehicle
	nVeh := len(vehicles)
	nTasks := len(tasks)
	mean := float64(nTasks) / float64(nVeh)
	var varianceSum float64
	for vi := 0; vi < nVeh; vi++ {
		c := vehicleTaskCount[vi]
		varianceSum += (float64(c) - mean) * (float64(c) - mean)
	}
	stddev := math.Sqrt(varianceSum / float64(nVeh))
	loadPenalty := stddev * 2000.0 // weight: 2km per stddev

	// Vehicle utilization bonus: penalize unused vehicles
	unusedPenalty := 0.0
	for vi := 0; vi < nVeh; vi++ {
		if vehicleTaskCount[vi] == 0 {
			unusedPenalty += 5000.0 // 5km penalty per unused vehicle
		}
	}

	// Fitness: minimize total distance + penalties
	// Using negative so higher fitness = better
	return -(totalDist + loadPenalty + unusedPenalty)
}

// ── Evolution ──

func (a *GeneticAlgorithmAssigner) evolve(pop []Chromosome, vehicles []VehicleInfo, tasks []TaskInfo) []Chromosome {
	n := len(pop)
	next := make([]Chromosome, 0, n)

	// Sort by fitness (descending)
	sorted := make([]Chromosome, n)
	copy(sorted, pop)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Fitness > sorted[j].Fitness
	})

	// Elitism: keep best chromosomes
	for i := 0; i < a.params.ElitismCount && i < n; i++ {
		clone := make([]int, len(sorted[i].Genes))
		copy(clone, sorted[i].Genes)
		next = append(next, Chromosome{Genes: clone, Fitness: sorted[i].Fitness})
	}

	// Fill rest through selection, crossover, mutation
	for len(next) < n {
		p1 := a.tournamentSelect(sorted)
		p2 := a.tournamentSelect(sorted)
		c1, c2 := a.crossover(p1, p2, tasks)
		a.mutate(&c1, vehicles)
		a.mutate(&c2, vehicles)
		next = append(next, c1)
		if len(next) < n {
			next = append(next, c2)
		}
	}

	// Evaluate new population
	a.evaluatePopulation(next, vehicles, tasks)
	return next
}

func (a *GeneticAlgorithmAssigner) tournamentSelect(pop []Chromosome) Chromosome {
	best := pop[rand.Intn(len(pop))]
	for i := 1; i < a.params.TournamentSize; i++ {
		c := pop[rand.Intn(len(pop))]
		if c.Fitness > best.Fitness {
			best = c
		}
	}
	return best
}

func (a *GeneticAlgorithmAssigner) crossover(p1, p2 Chromosome, tasks []TaskInfo) (Chromosome, Chromosome) {
	n := len(p1.Genes)
	if n <= 1 {
		c1 := make([]int, n)
		c2 := make([]int, n)
		copy(c1, p1.Genes)
		copy(c2, p2.Genes)
		return Chromosome{Genes: c1}, Chromosome{Genes: c2}
	}

	point := rand.Intn(n - 1) + 1

	c1Genes := make([]int, n)
	c2Genes := make([]int, n)
	copy(c1Genes[:point], p1.Genes[:point])
	copy(c1Genes[point:], p2.Genes[point:])
	copy(c2Genes[:point], p2.Genes[:point])
	copy(c2Genes[point:], p1.Genes[point:])

	return Chromosome{Genes: c1Genes}, Chromosome{Genes: c2Genes}
}

func (a *GeneticAlgorithmAssigner) mutate(ch *Chromosome, vehicles []VehicleInfo) {
	nVeh := len(vehicles)
	for i := range ch.Genes {
		if rand.Float64() < a.params.MutationRate {
			ch.Genes[i] = rand.Intn(nVeh)
		}
	}
}

// ── Distance helpers ──

// getRoadDistance calls route-service GetDistance, falls back to haversine.
func (a *GeneticAlgorithmAssigner) getRoadDistance(fromLat, fromLon, toLat, toLon float64) (float64, float64) {
	return a.getRoadDistanceCached(fromLat, fromLon, toLat, toLon)
}

func (a *GeneticAlgorithmAssigner) getRoadDistanceCached(fromLat, fromLon, toLat, toLon float64) (float64, float64) {
	if a.cache != nil {
		if d, dur, ok := a.cache.Get(fromLat, fromLon, toLat, toLon); ok {
			return d, dur
		}
	}

	dist, dur := a.callRouteService(fromLat, fromLon, toLat, toLon)
	if a.cache != nil {
		a.cache.Set(fromLat, fromLon, toLat, toLon, dist, dur)
	}
	return dist, dur
}

func (a *GeneticAlgorithmAssigner) callRouteService(fromLat, fromLon, toLat, toLon float64) (float64, float64) {
	if a.svc == nil || a.svc.RouteClient == nil {
		// Fall back to haversine + 30 km/h
		d := haversineGA(fromLat, fromLon, toLat, toLon)
		return d, (d / 1000) / 30 * 3600
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	resp, err := a.svc.RouteClient.GetDistance(ctx, &routev1.GetDistanceRequest{
		FromLat: fromLat,
		FromLon: fromLon,
		ToLat:   toLat,
		ToLon:   toLon,
	})
	if err != nil || resp == nil || resp.Code != 0 {
		d := haversineGA(fromLat, fromLon, toLat, toLon)
		return d, (d / 1000) / 30 * 3600
	}
	return resp.DistanceM, resp.DurationS
}

// getPointCoords retrieves the coordinates of a load/dump point from the DB.
func (a *GeneticAlgorithmAssigner) getPointCoords(pointID uint64) (float64, float64) {
	if a.svc == nil || a.svc.DB == nil {
		return 0, 0
	}
	var pt LoadingPoint
	if err := a.svc.DB.First(&pt, pointID).Error; err != nil {
		return 0, 0
	}
	return pt.Latitude, pt.Longitude
}

// haversineGA computes the great-circle distance in meters between two lat/lon points.
func haversineGA(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// ── BatchAssignResult ──

type BatchAssignResult struct {
	Assignments []BatchAssignment
	TotalDistM  float64
	TotalDurS   float64
}

type BatchAssignment struct {
	TaskIdx     int
	VehicleID   uint64
	LoadPointID uint64
	DumpPointID uint64
	DistanceM   float64
	DurationS   float64
}

// BatchOptimizeAssignments runs GA optimization and returns best assignments.
func (a *GeneticAlgorithmAssigner) BatchOptimizeAssignments(vehicles []VehicleInfo, tasks []TaskInfo) *BatchAssignResult {
	best := a.BatchOptimize(vehicles, tasks)
	if best == nil || math.IsInf(best.Fitness, -1) {
		return nil
	}

	result := &BatchAssignResult{}
	totalDist := 0.0
	totalDur := 0.0

	for ti, vi := range best.Genes {
		v := vehicles[vi]
		t := tasks[ti]

		emptyDist, emptyDur := a.getRoadDistanceCached(v.Latitude, v.Longitude, t.LoadLat, t.LoadLon)
		loadDist, loadDur := a.getRoadDistanceCached(t.LoadLat, t.LoadLon, t.DumpLat, t.DumpLon)

		dist := emptyDist + loadDist
		dur := emptyDur + loadDur
		totalDist += dist
		totalDur += dur

		result.Assignments = append(result.Assignments, BatchAssignment{
			TaskIdx:     t.Index,
			VehicleID:   v.ID,
			LoadPointID: t.LoadPointID,
			DumpPointID: t.DumpPointID,
			DistanceM:   dist,
			DurationS:   dur,
		})
	}

	result.TotalDistM = totalDist
	result.TotalDurS = totalDur
	return result
}
