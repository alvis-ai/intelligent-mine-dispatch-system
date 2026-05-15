package logic

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

func TestDefaultGAParams(t *testing.T) {
	p := DefaultGAParams()
	if p.PopulationSize != 50 {
		t.Errorf("PopulationSize = %d, want 50", p.PopulationSize)
	}
	if p.Generations != 100 {
		t.Errorf("Generations = %d, want 100", p.Generations)
	}
	if p.MutationRate != 0.1 {
		t.Errorf("MutationRate = %f, want 0.1", p.MutationRate)
	}
	if p.ElitismCount != 2 {
		t.Errorf("ElitismCount = %d, want 2", p.ElitismCount)
	}
}

func TestHaversineGA(t *testing.T) {
	// Same point
	d := haversineGA(39.9, 116.4, 39.9, 116.4)
	if d != 0 {
		t.Errorf("same point distance = %f, want 0", d)
	}

	// Approx 1 degree latitude ≈ 111km
	d = haversineGA(0, 0, 1, 0)
	if d < 110000 || d > 112000 {
		t.Errorf("1 degree lat = %f, want ~111000", d)
	}
}

func TestHaversineGA_Symmetry(t *testing.T) {
	d1 := haversineGA(39.9, 116.4, 31.2, 121.5)
	d2 := haversineGA(31.2, 121.5, 39.9, 116.4)
	if math.Abs(d1-d2) > 0.001 {
		t.Errorf("not symmetric: %f vs %f", d1, d2)
	}
}

func TestDistanceCache(t *testing.T) {
	cache := NewDistanceCache()

	// Initially empty
	_, _, ok := cache.Get(39.9, 116.4, 31.2, 121.5)
	if ok {
		t.Error("cache should be empty initially")
	}

	// Set and get
	cache.Set(39.9, 116.4, 31.2, 121.5, 1000.5, 120.0)
	d, dur, ok := cache.Get(39.9, 116.4, 31.2, 121.5)
	if !ok {
		t.Error("cache should have the entry")
	}
	if d != 1000.5 {
		t.Errorf("distance = %f, want 1000.5", d)
	}
	if dur != 120.0 {
		t.Errorf("duration = %f, want 120.0", dur)
	}
}

func TestDistanceCache_KeyUniqueness(t *testing.T) {
	cache := NewDistanceCache()
	cache.Set(1, 2, 3, 4, 100, 10)
	_, _, ok := cache.Get(1, 2, 3, 5) // different to_lon
	if ok {
		t.Error("cache should not return for different key")
	}
}

func TestInitializePopulation_HasHeuristic(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Latitude: 39.91, Longitude: 116.41},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadLat: 39.90, LoadLon: 116.40, DumpLat: 39.91, DumpLon: 116.41},
		{Index: 1, LoadLat: 39.91, LoadLon: 116.41, DumpLat: 39.90, DumpLon: 116.40},
	}

	pop := a.initializePopulation(vehicles, tasks)
	if len(pop) != a.params.PopulationSize {
		t.Errorf("population size = %d, want %d", len(pop), a.params.PopulationSize)
	}

	// First chromosome should be heuristic (nearest vehicle to each load)
	// Vehicle 1 at (39.90, 116.40) is nearest to task 0's load at (39.90, 116.40)
	if pop[0].Genes[0] != 0 {
		t.Errorf("heuristic: task 0 should be assigned to vehicle 0 (index), got %d", pop[0].Genes[0])
	}
}

func TestEvaluate_EqualDistribution(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Latitude: 39.91, Longitude: 116.41},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadLat: 39.90, LoadLon: 116.40, DumpLat: 39.91, DumpLon: 116.41},
		{Index: 1, LoadLat: 39.91, LoadLon: 116.41, DumpLat: 39.90, DumpLon: 116.40},
	}

	// Perfectly balanced: one task per vehicle
	ch := Chromosome{Genes: []int{0, 1}}
	fitness := a.evaluate(&ch, vehicles, tasks)

	if math.IsInf(fitness, -1) {
		t.Fatal("fitness should not be -inf")
	}
	if fitness >= 0 {
		t.Error("fitness should be negative (minimization problem)")
	}
}

func TestEvaluate_UnbalancedPenalty(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Latitude: 39.91, Longitude: 116.41},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadLat: 39.90, LoadLon: 116.40, DumpLat: 39.91, DumpLon: 116.41},
		{Index: 1, LoadLat: 39.91, LoadLon: 116.41, DumpLat: 39.90, DumpLon: 116.40},
	}

	// Balanced: 0,1 → vehicle 0 gets 1 task, vehicle 1 gets 1 task
	balanced := Chromosome{Genes: []int{0, 1}}
	// Unbalanced: 0,0 → vehicle 0 gets 2 tasks, vehicle 1 gets 0
	unbalanced := Chromosome{Genes: []int{0, 0}}

	fBalanced := a.evaluate(&balanced, vehicles, tasks)
	fUnbalanced := a.evaluate(&unbalanced, vehicles, tasks)

	// Balanced should have better (higher) fitness
	if fBalanced <= fUnbalanced {
		t.Errorf("balanced fitness %f should be > unbalanced fitness %f", fBalanced, fUnbalanced)
	}
}

func TestCrossover(t *testing.T) {
	// Use fixed seed for reproducibility
	rand.Seed(42)
	defer rand.Seed(0)

	a := &GeneticAlgorithmAssigner{}
	p1 := Chromosome{Genes: []int{0, 0, 1, 1, 2}}
	p2 := Chromosome{Genes: []int{2, 1, 0, 2, 1}}
	tasks := make([]TaskInfo, 5)

	c1, c2 := a.crossover(p1, p2, tasks)

	if len(c1.Genes) != 5 || len(c2.Genes) != 5 {
		t.Errorf("crossover children have wrong length: %d, %d", len(c1.Genes), len(c2.Genes))
	}

	// Verify genes come from both parents
	allP1 := true
	allP2 := true
	for i := range c1.Genes {
		if c1.Genes[i] != p1.Genes[i] {
			allP1 = false
		}
		if c1.Genes[i] != p2.Genes[i] {
			allP2 = false
		}
	}
	if allP1 || allP2 {
		t.Error("crossover child should be a mix of both parents")
	}
}

func TestMutate(t *testing.T) {
	rand.Seed(123)
	defer rand.Seed(0)

	a := &GeneticAlgorithmAssigner{params: GAParams{MutationRate: 1.0}} // always mutate
	vehicles := []VehicleInfo{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	ch := Chromosome{Genes: []int{0, 1, 2, 3, 4}}

	original := make([]int, len(ch.Genes))
	copy(original, ch.Genes)

	a.mutate(&ch, vehicles)

	// With 100% mutation rate and 5 vehicles, most genes should change
	changed := 0
	for i := range ch.Genes {
		if ch.Genes[i] != original[i] {
			changed++
		}
	}
	if changed == 0 {
		t.Error("with 1.0 mutation rate, at least some genes should have mutated")
	}

	// All genes should be valid vehicle indices
	for _, g := range ch.Genes {
		if g < 0 || g >= len(vehicles) {
			t.Errorf("mutated gene %d is out of range", g)
		}
	}
}

func TestTournamentSelect(t *testing.T) {
	rand.Seed(456)
	defer rand.Seed(0)

	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}
	pop := []Chromosome{
		{Genes: []int{0, 0}, Fitness: -1000},
		{Genes: []int{1, 1}, Fitness: -500},
	}

	winner := a.tournamentSelect(pop)
	_ = winner // tournament is probabilistic; just ensure it returns without panic
}

func TestBatchOptimize_Convergence(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Latitude: 39.92, Longitude: 116.42},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadLat: 39.90, LoadLon: 116.40, DumpLat: 39.91, DumpLon: 116.41},
		{Index: 1, LoadLat: 39.92, LoadLon: 116.42, DumpLat: 39.91, DumpLon: 116.41},
		{Index: 2, LoadLat: 39.90, LoadLon: 116.40, DumpLat: 39.92, DumpLon: 116.42},
	}

	result := a.BatchOptimize(vehicles, tasks)
	if result == nil {
		t.Fatal("BatchOptimize returned nil")
	}
	if len(result.Genes) != len(tasks) {
		t.Errorf("result has %d genes, want %d", len(result.Genes), len(tasks))
	}

	// All genes should be valid vehicle indices
	for i, g := range result.Genes {
		if g < 0 || g >= len(vehicles) {
			t.Errorf("gene[%d] = %d, out of range [0, %d)", i, g, len(vehicles))
		}
	}

	if math.IsInf(result.Fitness, -1) {
		t.Error("fitness should not be -inf after optimization")
	}
}

func TestBatchOptimize_SingleVehicle(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadLat: 39.91, LoadLon: 116.41, DumpLat: 39.92, DumpLon: 116.42},
	}

	result := a.BatchOptimize(vehicles, tasks)
	if result == nil {
		t.Fatal("BatchOptimize with 1 vehicle returned nil")
	}
	// All tasks must be assigned to vehicle 0 (only option)
	for i, g := range result.Genes {
		if g != 0 {
			t.Errorf("gene[%d] = %d, want 0 (only vehicle)", i, g)
		}
	}
}

func TestBatchOptimize_EmptyInput(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams()}

	result := a.BatchOptimize(nil, nil)
	if result != nil {
		t.Error("BatchOptimize with nil inputs should return nil")
	}

	result = a.BatchOptimize([]VehicleInfo{{ID: 1}}, nil)
	if result != nil {
		t.Error("BatchOptimize with nil tasks should return nil")
	}
}

func TestBatchOptimizeAssignments(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: DefaultGAParams(), cache: NewDistanceCache()}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadPointID: 1, DumpPointID: 2, LoadLat: 39.91, LoadLon: 116.41, DumpLat: 39.92, DumpLon: 116.42},
	}

	result := a.BatchOptimizeAssignments(vehicles, tasks)
	if result == nil {
		t.Fatal("BatchOptimizeAssignments returned nil")
	}
	if len(result.Assignments) != 1 {
		t.Errorf("got %d assignments, want 1", len(result.Assignments))
	}
	if result.Assignments[0].VehicleID != 1 {
		t.Errorf("VehicleID = %d, want 1", result.Assignments[0].VehicleID)
	}
	if result.TotalDistM <= 0 {
		t.Errorf("TotalDistM = %f, want > 0", result.TotalDistM)
	}
}

func TestEvolution_GenerationalImprovement(t *testing.T) {
	a := &GeneticAlgorithmAssigner{params: GAParams{
		PopulationSize: 20,
		Generations:    1,
		MutationRate:   0.2,
		ElitismCount:   2,
		TournamentSize: 3,
	}}
	vehicles := []VehicleInfo{
		{ID: 1, Latitude: 39.90, Longitude: 116.40},
		{ID: 2, Latitude: 39.91, Longitude: 116.41},
	}
	tasks := []TaskInfo{
		{Index: 0, LoadLat: 39.90, LoadLon: 116.40, DumpLat: 39.91, DumpLon: 116.41},
		{Index: 1, LoadLat: 39.91, LoadLon: 116.41, DumpLat: 39.90, DumpLon: 116.40},
	}

	// Run a single generation of evolution
	pop := a.initializePopulation(vehicles, tasks)
	a.evaluatePopulation(pop, vehicles, tasks)
	evolved := a.evolve(pop, vehicles, tasks)

	if len(evolved) != len(pop) {
		t.Errorf("evolved population size = %d, want %d", len(evolved), len(pop))
	}

	// Elitism should preserve the best chromosome
	sort.Slice(pop, func(i, j int) bool { return pop[i].Fitness > pop[j].Fitness })
	sort.Slice(evolved, func(i, j int) bool { return evolved[i].Fitness > evolved[j].Fitness })

	if evolved[0].Fitness < pop[0].Fitness {
		t.Errorf("best fitness decreased: before=%f, after=%f", pop[0].Fitness, evolved[0].Fitness)
	}
}

func TestGetPointCoords_Fallback(t *testing.T) {
	// Without a real DB, getPointCoords should return (0,0)
	a := &GeneticAlgorithmAssigner{}

	lat, lon := a.getPointCoords(999)
	if lat != 0 || lon != 0 {
		t.Errorf("without DB, getPointCoords should return (0,0), got (%f,%f)", lat, lon)
	}
}

func TestNewDistanceCache(t *testing.T) {
	cache := NewDistanceCache()
	if cache == nil {
		t.Fatal("NewDistanceCache returned nil")
	}
	if cache.data == nil {
		t.Error("NewDistanceCache should initialize the map")
	}
}
