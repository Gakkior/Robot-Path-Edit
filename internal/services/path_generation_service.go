// Package services 璺緞鐢熸垚鏈嶅姟瀹炵幇
package services

import (
	"context"
	"fmt"
	"math"
	"sort"

	"robot-path-editor/internal/domain"
)

// PathGenerationService 璺緞鐢熸垚鏈嶅姟鎺ュ彛
type PathGenerationService interface {
	GenerateShortestPaths(ctx context.Context, startNodeID domain.NodeID) ([]domain.Path, error)
	GenerateFullConnectivity(ctx context.Context) ([]domain.Path, error)
	GenerateTreeStructure(ctx context.Context, rootNodeID domain.NodeID) ([]domain.Path, error)
	GenerateNearestNeighborPaths(ctx context.Context, maxDistance float64) ([]domain.Path, error)
	GenerateGridPaths(ctx context.Context, connectDiagonal bool) ([]domain.Path, error)
}

type pathGenerationService struct {
	nodeService NodeService
	pathService PathService
}

// NewPathGenerationService 鍒涘缓璺緞鐢熸垚鏈嶅姟
func NewPathGenerationService(nodeService NodeService, pathService PathService) PathGenerationService {
	return &pathGenerationService{
		nodeService: nodeService,
		pathService: pathService,
	}
}

// GenerateShortestPaths 鐢熸垚浠庢寚瀹氳妭鐐瑰埌鎵€鏈夊叾浠栬妭鐐圭殑鏈€鐭矾寰?(鍩轰簬Dijkstra绠楁硶鐨勭畝鍖栫増)
func (s *pathGenerationService) GenerateShortestPaths(ctx context.Context, startNodeID domain.NodeID) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣鍒楄〃澶辫触: %v", err)
	}

	// 鎵惧埌璧峰鑺傜偣
	var startNode *domain.Node
	nodeMap := make(map[domain.NodeID]*domain.Node)
	for i := range nodes {
		nodeMap[nodes[i].ID] = nodes[i]
		if nodes[i].ID == startNodeID {
			startNode = nodes[i]
		}
	}

	if startNode == nil {
		return nil, fmt.Errorf("璧峰鑺傜偣涓嶅瓨鍦? %s", startNodeID)
	}

	// 鑾峰彇鐜版湁璺緞浠ユ瀯寤洪偦鎺ュ浘
	existingPaths, err := s.pathService.ListPaths(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇璺緞鍒楄〃澶辫触: %v", err)
	}

	// 鏋勫缓閭绘帴鍥?
	adjacencyList := make(map[domain.NodeID][]domain.NodeID)
	for _, path := range existingPaths {
		adjacencyList[path.StartNodeID] = append(adjacencyList[path.StartNodeID], path.EndNodeID)
		adjacencyList[path.EndNodeID] = append(adjacencyList[path.EndNodeID], path.StartNodeID) // 鍙屽悜
	}

	// 浣跨敤绠€鍖栫殑Dijkstra绠楁硶
	distances := make(map[domain.NodeID]float64)
	previous := make(map[domain.NodeID]domain.NodeID)
	visited := make(map[domain.NodeID]bool)

	// 鍒濆鍖栬窛绂?
	for _, node := range nodes {
		distances[node.ID] = math.Inf(1)
	}
	distances[startNodeID] = 0

	// Dijkstra涓诲惊鐜?
	for len(visited) < len(nodes) {
		// 鎵惧埌鏈闂妭鐐逛腑璺濈鏈€灏忕殑
		minDist := math.Inf(1)
		var currentNode domain.NodeID
		for nodeID, dist := range distances {
			if !visited[nodeID] && dist < minDist {
				minDist = dist
				currentNode = nodeID
			}
		}

		if minDist == math.Inf(1) {
			break // 鏃犳硶鍒拌揪鐨勮妭鐐?
		}

		visited[currentNode] = true

		// 鏇存柊閭诲眳鑺傜偣鐨勮窛绂?
		for _, neighborID := range adjacencyList[currentNode] {
			if visited[neighborID] {
				continue
			}

			currentNodeData := nodeMap[currentNode]
			neighborNodeData := nodeMap[neighborID]
			edgeWeight := calculateDistance(currentNodeData.Position, neighborNodeData.Position)

			newDist := distances[currentNode] + edgeWeight
			if newDist < distances[neighborID] {
				distances[neighborID] = newDist
				previous[neighborID] = currentNode
			}
		}
	}

	// 鐢熸垚鏈€鐭矾寰?
	var paths []domain.Path
	for nodeID, prevNodeID := range previous {
		if nodeID != startNodeID && prevNodeID != "" {
			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("shortest_%s_%s", prevNodeID, nodeID)),
				Name:        fmt.Sprintf("鏈€鐭矾寰? %s -> %s", prevNodeID, nodeID),
				Type:        "shortest",
				Status:      "active",
				StartNodeID: prevNodeID,
				EndNodeID:   nodeID,
				Metadata: domain.ObjectMeta{
					Annotations: map[string]string{
						"distance":  fmt.Sprintf("%.2f", distances[nodeID]),
						"algorithm": "dijkstra",
					},
				},
			}
			paths = append(paths, path)
		}
	}

	return paths, nil
}

// GenerateFullConnectivity 鐢熸垚瀹屽叏杩為€氬浘锛堟墍鏈夎妭鐐逛袱涓ょ浉杩烇級
func (s *pathGenerationService) GenerateFullConnectivity(ctx context.Context) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣鍒楄〃澶辫触: %v", err)
	}

	var paths []domain.Path
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			node1, node2 := nodes[i], nodes[j]
			distance := calculateDistance(node1.Position, node2.Position)

			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("full_%s_%s", node1.ID, node2.ID)),
				Name:        fmt.Sprintf("杩炴帴: %s <-> %s", node1.Name, node2.Name),
				Type:        "full_connectivity",
				Status:      "active",
				StartNodeID: node1.ID,
				EndNodeID:   node2.ID,
				Metadata: domain.ObjectMeta{
					Annotations: map[string]string{
						"distance":  fmt.Sprintf("%.2f", distance),
						"algorithm": "full_connectivity",
					},
				},
			}
			paths = append(paths, path)
		}
	}

	return paths, nil
}

// GenerateTreeStructure 鐢熸垚鏍戠姸缁撴瀯锛堟渶灏忕敓鎴愭爲锛?
func (s *pathGenerationService) GenerateTreeStructure(ctx context.Context, rootNodeID domain.NodeID) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣鍒楄〃澶辫触: %v", err)
	}

	// 鎵惧埌鏍硅妭鐐?
	var rootNode *domain.Node
	nodeMap := make(map[domain.NodeID]*domain.Node)
	for i := range nodes {
		nodeMap[nodes[i].ID] = nodes[i]
		if nodes[i].ID == rootNodeID {
			rootNode = nodes[i]
		}
	}

	if rootNode == nil {
		return nil, fmt.Errorf("鏍硅妭鐐逛笉瀛樺湪: %s", rootNodeID)
	}

	// 浣跨敤Prim绠楁硶鐢熸垚鏈€灏忕敓鎴愭爲
	visited := make(map[domain.NodeID]bool)
	visited[rootNodeID] = true

	type edge struct {
		from, to domain.NodeID
		weight   float64
	}

	var paths []domain.Path

	for len(visited) < len(nodes) {
		var minEdge *edge

		// 鎵惧埌杩炴帴宸茶闂妭鐐瑰拰鏈闂妭鐐圭殑鏈€鐭竟
		for visitedNodeID := range visited {
			visitedNode := nodeMap[visitedNodeID]
			for _, node := range nodes {
				if !visited[node.ID] {
					weight := calculateDistance(visitedNode.Position, node.Position)
					if minEdge == nil || weight < minEdge.weight {
						minEdge = &edge{
							from:   visitedNodeID,
							to:     node.ID,
							weight: weight,
						}
					}
				}
			}
		}

		if minEdge == nil {
			break // 鏃犳硶杩炴帴鏇村鑺傜偣
		}

		// 娣诲姞杈瑰埌MST
		visited[minEdge.to] = true
		path := domain.Path{
			ID:          domain.PathID(fmt.Sprintf("tree_%s_%s", minEdge.from, minEdge.to)),
			Name:        fmt.Sprintf("鏍戣繛鎺? %s -> %s", minEdge.from, minEdge.to),
			Type:        "tree",
			Status:      "active",
			StartNodeID: minEdge.from,
			EndNodeID:   minEdge.to,
			Metadata: domain.ObjectMeta{
				Annotations: map[string]string{
					"weight":    fmt.Sprintf("%.2f", minEdge.weight),
					"algorithm": "prim_mst",
				},
			},
		}
		paths = append(paths, path)
	}

	return paths, nil
}

// GenerateNearestNeighborPaths 鐢熸垚鏈€杩戦偦璺緞锛堟瘡涓妭鐐硅繛鎺ュ埌鏈€杩戠殑N涓偦灞咃級
func (s *pathGenerationService) GenerateNearestNeighborPaths(ctx context.Context, maxDistance float64) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣鍒楄〃澶辫触: %v", err)
	}

	type neighbor struct {
		nodeID   domain.NodeID
		distance float64
	}

	var paths []domain.Path
	pathSet := make(map[string]bool) // 闃叉閲嶅璺緞

	for _, node := range nodes {
		var neighbors []neighbor

		// 璁＄畻鍒版墍鏈夊叾浠栬妭鐐圭殑璺濈
		for _, otherNode := range nodes {
			if node.ID != otherNode.ID {
				distance := calculateDistance(node.Position, otherNode.Position)
				if distance <= maxDistance {
					neighbors = append(neighbors, neighbor{
						nodeID:   otherNode.ID,
						distance: distance,
					})
				}
			}
		}

		// 鎸夎窛绂绘帓搴?
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].distance < neighbors[j].distance
		})

		// 杩炴帴鍒版渶杩戠殑閭诲眳锛堟渶澶?涓級
		maxNeighbors := min(3, len(neighbors))
		for i := 0; i < maxNeighbors; i++ {
			neighbor := neighbors[i]

			// 鍒涘缓鍞竴鐨勮矾寰勬爣璇嗙锛堥槻姝㈤噸澶嶏級
			pathKey := fmt.Sprintf("%s_%s", min(string(node.ID), string(neighbor.nodeID)), max(string(node.ID), string(neighbor.nodeID)))
			if pathSet[pathKey] {
				continue
			}
			pathSet[pathKey] = true

			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("neighbor_%s_%s", node.ID, neighbor.nodeID)),
				Name:        fmt.Sprintf("鏈€杩戦偦: %s <-> %s", node.Name, neighbor.nodeID),
				Type:        "nearest_neighbor",
				Status:      "active",
				StartNodeID: node.ID,
				EndNodeID:   neighbor.nodeID,
				Metadata: domain.ObjectMeta{
					Annotations: map[string]string{
						"distance":  fmt.Sprintf("%.2f", neighbor.distance),
						"algorithm": "nearest_neighbor",
					},
				},
			}
			paths = append(paths, path)
		}
	}

	return paths, nil
}

// GenerateGridPaths 鐢熸垚缃戞牸鐘惰矾寰勶紙閫傜敤浜庤鍒欐帓鍒楃殑鑺傜偣锛?
func (s *pathGenerationService) GenerateGridPaths(ctx context.Context, connectDiagonal bool) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣鍒楄〃澶辫触: %v", err)
	}

	if len(nodes) == 0 {
		return []domain.Path{}, nil
	}

	// 鎸変綅缃帓搴忚妭鐐癸紝鍒涘缓缃戞牸缁撴瀯
	sort.Slice(nodes, func(i, j int) bool {
		if math.Abs(nodes[i].Position.Y-nodes[j].Position.Y) < 10 { // 鍚屼竴琛?
			return nodes[i].Position.X < nodes[j].Position.X
		}
		return nodes[i].Position.Y < nodes[j].Position.Y
	})

	var paths []domain.Path
	tolerance := 50.0 // 浣嶇疆瀹瑰樊

	// 姘村钩杩炴帴锛堝悓涓€琛岀殑鐩搁偦鑺傜偣锛?
	for i := 0; i < len(nodes)-1; i++ {
		current := nodes[i]
		next := nodes[i+1]

		// 妫€鏌ユ槸鍚﹀湪鍚屼竴琛屼笖鐩搁偦
		if math.Abs(current.Position.Y-next.Position.Y) < tolerance {
			distance := calculateDistance(current.Position, next.Position)
			if distance < tolerance*3 { // 鐩搁偦鍒ゆ柇
				path := domain.Path{
					ID:          domain.PathID(fmt.Sprintf("grid_h_%s_%s", current.ID, next.ID)),
					Name:        fmt.Sprintf("缃戞牸姘村钩: %s -> %s", current.Name, next.Name),
					Type:        "grid_horizontal",
					Status:      "active",
					StartNodeID: current.ID,
					EndNodeID:   next.ID,
					Metadata: domain.ObjectMeta{
						Annotations: map[string]string{
							"distance":  fmt.Sprintf("%.2f", distance),
							"algorithm": "grid",
						},
					},
				}
				paths = append(paths, path)
			}
		}
	}

	// 鍨傜洿杩炴帴锛堝悓涓€鍒楃殑鐩搁偦鑺傜偣锛?
	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i >= j {
				continue
			}

			// 妫€鏌ユ槸鍚﹀湪鍚屼竴鍒?
			if math.Abs(node1.Position.X-node2.Position.X) < tolerance {
				distance := calculateDistance(node1.Position, node2.Position)
				if distance < tolerance*3 {
					path := domain.Path{
						ID:          domain.PathID(fmt.Sprintf("grid_v_%s_%s", node1.ID, node2.ID)),
						Name:        fmt.Sprintf("缃戞牸鍨傜洿: %s -> %s", node1.Name, node2.Name),
						Type:        "grid_vertical",
						Status:      "active",
						StartNodeID: node1.ID,
						EndNodeID:   node2.ID,
						Metadata: domain.ObjectMeta{
							Annotations: map[string]string{
								"distance":  fmt.Sprintf("%.2f", distance),
								"algorithm": "grid",
							},
						},
					}
					paths = append(paths, path)
				}
			}
		}
	}

	// 瀵硅绾胯繛鎺ワ紙濡傛灉鍚敤锛?
	if connectDiagonal {
		for i, node1 := range nodes {
			for j, node2 := range nodes {
				if i >= j {
					continue
				}

				distance := calculateDistance(node1.Position, node2.Position)
				dx := math.Abs(node1.Position.X - node2.Position.X)
				dy := math.Abs(node1.Position.Y - node2.Position.Y)

				// 妫€鏌ユ槸鍚︿负瀵硅绾匡紙45搴﹁锛?
				if math.Abs(dx-dy) < tolerance && distance < tolerance*2 {
					path := domain.Path{
						ID:          domain.PathID(fmt.Sprintf("grid_d_%s_%s", node1.ID, node2.ID)),
						Name:        fmt.Sprintf("缃戞牸瀵硅: %s -> %s", node1.Name, node2.Name),
						Type:        "grid_diagonal",
						Status:      "active",
						StartNodeID: node1.ID,
						EndNodeID:   node2.ID,
						Metadata: domain.ObjectMeta{
							Annotations: map[string]string{
								"distance":  fmt.Sprintf("%.2f", distance),
								"algorithm": "grid",
							},
						},
					}
					paths = append(paths, path)
				}
			}
		}
	}

	return paths, nil
}

// 宸ュ叿鍑芥暟

// calculateDistance 璁＄畻涓ょ偣涔嬮棿鐨勬鍑犻噷寰楄窛绂?
func calculateDistance(pos1, pos2 domain.Position) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	dz := pos1.Z - pos2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// min 杩斿洖涓や釜鍙瘮杈冨€肩殑鏈€灏忓€?
func min[T ~int | ~string](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// max 杩斿洖涓や釜鍙瘮杈冨€肩殑鏈€澶у€?
func max[T ~int | ~string](a, b T) T {
	if a > b {
		return a
	}
	return b
}
