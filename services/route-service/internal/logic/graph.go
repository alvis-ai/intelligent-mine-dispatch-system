package logic

import (
	"container/heap"
	"math"
	"sort"

	routev1 "github.com/aicong/mine-dispatch/proto/route/v1"
	"github.com/aicong/mine-dispatch/services/route-service/internal/model"
)

// Edge represents a directed edge in the routing graph.
type Edge struct {
	ToNodeID    uint64
	DistanceM   float64
	MaxSpeedKMH int32
	EdgeID      uint64
}

// Graph is an adjacency-list road network graph.
type Graph struct {
	Nodes map[uint64]*model.RoadNode
	Edges map[uint64][]Edge // adjacency list keyed by from_node_id
}

// BuildGraph builds a Graph from database models.
func BuildGraph(nodes []model.RoadNode, edges []model.RoadEdge) *Graph {
	g := &Graph{
		Nodes: make(map[uint64]*model.RoadNode),
		Edges: make(map[uint64][]Edge),
	}
	for i := range nodes {
		g.Nodes[nodes[i].ID] = &nodes[i]
	}
	for _, e := range edges {
		g.Edges[e.FromNodeID] = append(g.Edges[e.FromNodeID], Edge{
			ToNodeID:    e.ToNodeID,
			DistanceM:   e.DistanceM,
			MaxSpeedKMH: e.MaxSpeedKMH,
			EdgeID:      e.ID,
		})
		if !e.IsOneway {
			g.Edges[e.ToNodeID] = append(g.Edges[e.ToNodeID], Edge{
				ToNodeID:    e.FromNodeID,
				DistanceM:   e.DistanceM,
				MaxSpeedKMH: e.MaxSpeedKMH,
				EdgeID:      e.ID,
			})
		}
	}
	return g
}

// FindNearestNode finds the graph node closest to a given lat/lon point.
func (g *Graph) FindNearestNode(lat, lon float64) (uint64, float64) {
	var bestID uint64
	bestDist := math.MaxFloat64
	for _, node := range g.Nodes {
		d := haversine(lat, lon, node.Latitude, node.Longitude)
		if d < bestDist {
			bestDist = d
			bestID = node.ID
		}
	}
	return bestID, bestDist
}

// haversine computes distance in meters between two lat/lon points.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// ── Priority queue for Dijkstra / A* ──

type pqItem struct {
	nodeID   uint64
	priority float64 // f(n) = g(n) for Dijkstra, g(n) + h(n) for A*
	index    int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].priority < pq[j].priority }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *priorityQueue) Push(x interface{}) { n := len(*pq); item := x.(*pqItem); item.index = n; *pq = append(*pq, item) }
func (pq *priorityQueue) Pop() interface{}   { old := *pq; n := len(old); item := old[n-1]; old[n-1] = nil; item.index = -1; *pq = old[:n-1]; return item }

// ── Path finding result ──

type PathResult struct {
	Found      bool
	TotalDistM float64
	TotalDurS  float64
	NodeIDs    []uint64
	EdgeIDs    []uint64
}

// Dijkstra finds the shortest path between two nodes.
func (g *Graph) Dijkstra(fromID, toID uint64) *PathResult {
	return g.findPath(fromID, toID, false)
}

// AStar finds the shortest path using A* with Haversine heuristic.
func (g *Graph) AStar(fromID, toID uint64) *PathResult {
	return g.findPath(fromID, toID, true)
}

func (g *Graph) findPath(fromID, toID uint64, useAStar bool) *PathResult {
	if _, ok := g.Nodes[fromID]; !ok {
		return &PathResult{}
	}
	if _, ok := g.Nodes[toID]; !ok {
		return &PathResult{}
	}

	// gScore: actual shortest distance from start to each node
	gScore := make(map[uint64]float64)
	prev := make(map[uint64]uint64)
	usedEdge := make(map[uint64]uint64)

	for id := range g.Nodes {
		gScore[id] = math.MaxFloat64
	}
	gScore[fromID] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{nodeID: fromID, priority: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem)
		if cur.nodeID == toID {
			break
		}
		// Skip stale entries: if g(cur) has improved since this entry was pushed
		if useAStar {
			hCur := haversine(
				g.Nodes[cur.nodeID].Latitude, g.Nodes[cur.nodeID].Longitude,
				g.Nodes[toID].Latitude, g.Nodes[toID].Longitude,
			)
			if cur.priority > gScore[cur.nodeID]+hCur*0.8 {
				continue
			}
		} else if cur.priority > gScore[cur.nodeID] {
			continue
		}

		for _, edge := range g.Edges[cur.nodeID] {
			alt := gScore[cur.nodeID] + edge.DistanceM
			if alt < gScore[edge.ToNodeID] {
				gScore[edge.ToNodeID] = alt
				prev[edge.ToNodeID] = cur.nodeID
				usedEdge[edge.ToNodeID] = edge.EdgeID

				priority := alt
				if useAStar {
					h := haversine(
						g.Nodes[edge.ToNodeID].Latitude, g.Nodes[edge.ToNodeID].Longitude,
						g.Nodes[toID].Latitude, g.Nodes[toID].Longitude,
					)
					priority = alt + h*0.8
				}
				heap.Push(pq, &pqItem{nodeID: edge.ToNodeID, priority: priority})
			}
		}
	}

	if gScore[toID] >= math.MaxFloat64 {
		return &PathResult{}
	}

	// Reconstruct path
	var nodeIDs []uint64
	for at := toID; at != fromID; at = prev[at] {
		nodeIDs = append(nodeIDs, at)
	}
	nodeIDs = append(nodeIDs, fromID)
	reverse(nodeIDs)

	var edgeIDs []uint64
	for i := 1; i < len(nodeIDs); i++ {
		if eid, ok := usedEdge[nodeIDs[i]]; ok {
			edgeIDs = append(edgeIDs, eid)
		}
	}

	// Compute duration
	totalDur := 0.0
	for i := 0; i < len(edgeIDs); i++ {
		for _, edges := range g.Edges {
			for _, e := range edges {
				if e.EdgeID == edgeIDs[i] {
					speed := float64(e.MaxSpeedKMH)
					if speed <= 0 {
						speed = 30
					}
					totalDur += (e.DistanceM / 1000) / speed * 3600
					goto nextEdge
				}
			}
		}
	nextEdge:
	}

	return &PathResult{
		Found:      true,
		TotalDistM: gScore[toID],
		TotalDurS:  totalDur,
		NodeIDs:    nodeIDs,
		EdgeIDs:    edgeIDs,
	}
}

func reverse(s []uint64) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// PathToProto converts a PathResult to a RoutePath proto message.
func (g *Graph) PathToProto(res *PathResult) *routev1.RoutePath {
	if res == nil || !res.Found {
		return nil
	}
	path := &routev1.RoutePath{
		TotalDistanceM: res.TotalDistM,
		TotalDurationS: res.TotalDurS,
		NodeIds:        res.NodeIDs,
		EdgeIds:        res.EdgeIDs,
	}
	for _, nid := range res.NodeIDs {
		if node, ok := g.Nodes[nid]; ok {
			path.Points = append(path.Points, &routev1.Point{
				Latitude:  node.Latitude,
				Longitude: node.Longitude,
			})
		}
	}
	return path
}

// SortNodesByDistance sorts node IDs by Haversine distance from a point.
func SortNodesByDistance(nodes []model.RoadNode, lat, lon float64) []uint64 {
	type distNode struct {
		id   uint64
		dist float64
	}
	var withDist []distNode
	for _, n := range nodes {
		d := haversine(lat, lon, n.Latitude, n.Longitude)
		withDist = append(withDist, distNode{id: n.ID, dist: d})
	}
	sort.Slice(withDist, func(i, j int) bool {
		return withDist[i].dist < withDist[j].dist
	})
	var ids []uint64
	for _, dn := range withDist {
		ids = append(ids, dn.id)
	}
	return ids
}
