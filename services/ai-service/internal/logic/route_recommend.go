package logic

import (
	"container/heap"
	"math"

	aiv1 "github.com/aicong/mine-dispatch/proto/ai/v1"
	"github.com/aicong/mine-dispatch/services/ai-service/internal/model"
)

// AI-weighted Dijkstra for congestion-aware route recommendation.

type aiEdge struct {
	ToNodeID  uint64
	BaseDistM float64
	MaxSpeed  int32
	EdgeID    uint64
}

type aiGraph struct {
	Nodes    map[uint64]*model.RoadNode
	EdgeList map[uint64][]aiEdge
}

func buildAiGraph(nodes []model.RoadNode, edges []model.RoadEdge) *aiGraph {
	g := &aiGraph{
		Nodes:    make(map[uint64]*model.RoadNode),
		EdgeList: make(map[uint64][]aiEdge),
	}
	for i := range nodes {
		g.Nodes[nodes[i].ID] = &nodes[i]
	}
	for _, e := range edges {
		g.EdgeList[e.FromNodeID] = append(g.EdgeList[e.FromNodeID], aiEdge{
			ToNodeID:  e.ToNodeID,
			BaseDistM: e.DistanceM,
			MaxSpeed:  e.MaxSpeedKMH,
			EdgeID:    e.ID,
		})
		if !e.IsOneway {
			g.EdgeList[e.ToNodeID] = append(g.EdgeList[e.ToNodeID], aiEdge{
				ToNodeID:  e.FromNodeID,
				BaseDistM: e.DistanceM,
				MaxSpeed:  e.MaxSpeedKMH,
				EdgeID:    e.ID,
			})
		}
	}
	return g
}

type aiPathResult struct {
	Found      bool
	TotalDistM float64
	TotalDurS  float64
	NodeIDs    []uint64
	EdgeIDs    []uint64
}

type aiPqItem struct {
	nodeID   uint64
	priority float64
	index    int
}

type aiPriorityQueue []*aiPqItem

func (pq aiPriorityQueue) Len() int           { return len(pq) }
func (pq aiPriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq aiPriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *aiPriorityQueue) Push(x any)        { n := len(*pq); item := x.(*aiPqItem); item.index = n; *pq = append(*pq, item) }
func (pq *aiPriorityQueue) Pop() any          { old := *pq; n := len(old); item := old[n-1]; old[n-1] = nil; item.index = -1; *pq = old[:n-1]; return item }

func (g *aiGraph) findPath(fromID, toID uint64, congestion map[uint64]float64) *aiPathResult {
	if _, ok := g.Nodes[fromID]; !ok {
		return &aiPathResult{}
	}
	if _, ok := g.Nodes[toID]; !ok {
		return &aiPathResult{}
	}

	gScore := make(map[uint64]float64)
	prev := make(map[uint64]uint64)
	usedEdge := make(map[uint64]uint64)

	for id := range g.Nodes {
		gScore[id] = math.MaxFloat64
	}
	gScore[fromID] = 0

	pq := &aiPriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &aiPqItem{nodeID: fromID, priority: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*aiPqItem)
		if cur.nodeID == toID {
			break
		}
		if cur.priority > gScore[cur.nodeID] {
			continue
		}

		for _, edge := range g.EdgeList[cur.nodeID] {
			// Apply congestion multiplier to effective distance
			cong := congestion[edge.EdgeID]
			effectiveDist := edge.BaseDistM * (1 + cong*2)
			alt := gScore[cur.nodeID] + effectiveDist
			if alt < gScore[edge.ToNodeID] {
				gScore[edge.ToNodeID] = alt
				prev[edge.ToNodeID] = cur.nodeID
				usedEdge[edge.ToNodeID] = edge.EdgeID
				heap.Push(pq, &aiPqItem{nodeID: edge.ToNodeID, priority: alt})
			}
		}
	}

	if gScore[toID] >= math.MaxFloat64 {
		return &aiPathResult{}
	}

	var nodeIDs []uint64
	for at := toID; at != fromID; at = prev[at] {
		nodeIDs = append(nodeIDs, at)
	}
	nodeIDs = append(nodeIDs, fromID)
	reverseAI(nodeIDs)

	var edgeIDs []uint64
	for i := 1; i < len(nodeIDs); i++ {
		if eid, ok := usedEdge[nodeIDs[i]]; ok {
			edgeIDs = append(edgeIDs, eid)
		}
	}

	totalDur := 0.0
	totalDist := 0.0
	for _, eid := range edgeIDs {
		for _, edges := range g.EdgeList {
			for _, e := range edges {
				if e.EdgeID == eid {
					speed := float64(e.MaxSpeed)
					if speed <= 0 {
						speed = 30
					}
					totalDur += (e.BaseDistM / 1000) / speed * 3600
					totalDist += e.BaseDistM
					goto nextEdgeAI
				}
			}
		}
	nextEdgeAI:
	}

	return &aiPathResult{
		Found:      true,
		TotalDistM: totalDist,
		TotalDurS:  totalDur,
		NodeIDs:    nodeIDs,
		EdgeIDs:    edgeIDs,
	}
}

func (g *aiGraph) findNearestNodeID(lat, lon float64) uint64 {
	var bestID uint64
	bestDist := math.MaxFloat64
	for _, n := range g.Nodes {
		d := haversineAI(lat, lon, n.Latitude, n.Longitude)
		if d < bestDist {
			bestDist = d
			bestID = n.ID
		}
	}
	return bestID
}

func reverseAI(s []uint64) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (l *AiLogic) RecommendRoute(in *aiv1.RecommendRouteRequest) (*aiv1.RecommendRouteResponse, error) {
	mineID := in.MineId
	if mineID == 0 {
		mineID = 1
	}

	var nodes []model.RoadNode
	var edges []model.RoadEdge

	l.svc.DB.Where("mine_id = ?", mineID).Find(&nodes)
	l.svc.DB.Where("mine_id = ?", mineID).Find(&edges)

	if len(nodes) == 0 {
		return &aiv1.RecommendRouteResponse{Code: 404, Message: "no road data for this mine"}, nil
	}

	g := buildAiGraph(nodes, edges)
	fromID := g.findNearestNodeID(in.FromLat, in.FromLon)
	toID := g.findNearestNodeID(in.ToLat, in.ToLon)

	if fromID == 0 || toID == 0 {
		return &aiv1.RecommendRouteResponse{Code: 404, Message: "no nearby road nodes"}, nil
	}
	if fromID == toID {
		return &aiv1.RecommendRouteResponse{
			Code:    0,
			Message: "already at destination",
			Data: &aiv1.AiRoutePath{
				Points: []*aiv1.RoutePoint{
					{Latitude: in.FromLat, Longitude: in.FromLon},
					{Latitude: in.ToLat, Longitude: in.ToLon},
				},
				TotalDistanceM: haversineAI(in.FromLat, in.FromLon, in.ToLat, in.ToLon),
			},
		}, nil
	}

	// Get congestion weights if avoid_congestion is enabled
	congestion := make(map[uint64]float64)
	if in.AvoidCongestion {
		congResp, err := l.PredictCongestion(&aiv1.PredictCongestionRequest{
			MineId:         mineID,
			LookbackMinutes: 30,
		})
		if err == nil && congResp != nil {
			for _, ec := range congResp.Data {
				congestion[ec.EdgeId] = ec.CongestionScore
			}
		}
	}

	result := g.findPath(fromID, toID, congestion)
	if result == nil || !result.Found {
		// Fallback: direct haversine
		d := haversineAI(in.FromLat, in.FromLon, in.ToLat, in.ToLon)
		return &aiv1.RecommendRouteResponse{
			Code:    0,
			Message: "direct (no road path found)",
			Data: &aiv1.AiRoutePath{
				Points: []*aiv1.RoutePoint{
					{Latitude: in.FromLat, Longitude: in.FromLon},
					{Latitude: in.ToLat, Longitude: in.ToLon},
				},
				TotalDistanceM: d,
				TotalDurationS: (d / 1000) / 30 * 3600,
			},
		}, nil
	}

	path := &aiv1.AiRoutePath{
		NodeIds:          result.NodeIDs,
		EdgeIds:          result.EdgeIDs,
		TotalDistanceM:   result.TotalDistM,
		TotalDurationS:   result.TotalDurS,
	}
	for _, nid := range result.NodeIDs {
		if n, ok := g.Nodes[nid]; ok {
			path.Points = append(path.Points, &aiv1.RoutePoint{
				Latitude:  n.Latitude,
				Longitude: n.Longitude,
			})
		}
	}

	return &aiv1.RecommendRouteResponse{
		Code:    0,
		Message: "success",
		Data:    path,
	}, nil
}
