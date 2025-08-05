// Package services 路径生成服务实现
package services

import (
	"context"
	"fmt"
	"math"
	"sort"

	"robot-path-editor/internal/domain"
)

// PathGenerationService 路径生成服务接口
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

// NewPathGenerationService 创建路径生成服务
func NewPathGenerationService(nodeService NodeService, pathService PathService) PathGenerationService {
	return &pathGenerationService{
		nodeService: nodeService,
		pathService: pathService,
	}
}

// GenerateShortestPaths 生成从指定节点到所有其他节点的最短路�?(基于Dijkstra算法的简化版)
func (s *pathGenerationService) GenerateShortestPaths(ctx context.Context, startNodeID domain.NodeID) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	// 找到起始节点
	var startNode *domain.Node
	nodeMap := make(map[domain.NodeID]*domain.Node)
	for i := range nodes {
		nodeMap[nodes[i].ID] = nodes[i]
		if nodes[i].ID == startNodeID {
			startNode = nodes[i]
		}
	}

	if startNode == nil {
		return nil, fmt.Errorf("起始节点不存�? %s", startNodeID)
	}

	// 获取现有路径以构建邻接图
	existingPaths, err := s.pathService.ListPaths(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取路径列表失败: %v", err)
	}

	// 构建邻接�?
	adjacencyList := make(map[domain.NodeID][]domain.NodeID)
	for _, path := range existingPaths {
		adjacencyList[path.StartNodeID] = append(adjacencyList[path.StartNodeID], path.EndNodeID)
		adjacencyList[path.EndNodeID] = append(adjacencyList[path.EndNodeID], path.StartNodeID) // 双向
	}

	// 使用简化的Dijkstra算法
	distances := make(map[domain.NodeID]float64)
	previous := make(map[domain.NodeID]domain.NodeID)
	visited := make(map[domain.NodeID]bool)

	// 初始化距�?
	for _, node := range nodes {
		distances[node.ID] = math.Inf(1)
	}
	distances[startNodeID] = 0

	// Dijkstra主循�?
	for len(visited) < len(nodes) {
		// 找到未访问节点中距离最小的
		minDist := math.Inf(1)
		var currentNode domain.NodeID
		for nodeID, dist := range distances {
			if !visited[nodeID] && dist < minDist {
				minDist = dist
				currentNode = nodeID
			}
		}

		if minDist == math.Inf(1) {
			break // 无法到达的节�?
		}

		visited[currentNode] = true

		// 更新邻居节点的距�?
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

	// 生成最短路�?
	var paths []domain.Path
	for nodeID, prevNodeID := range previous {
		if nodeID != startNodeID && prevNodeID != "" {
			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("shortest_%s_%s", prevNodeID, nodeID)),
				Name:        fmt.Sprintf("最短路�? %s -> %s", prevNodeID, nodeID),
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

// GenerateFullConnectivity 生成完全连通图（所有节点两两相连）
func (s *pathGenerationService) GenerateFullConnectivity(ctx context.Context) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	var paths []domain.Path
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			node1, node2 := nodes[i], nodes[j]
			distance := calculateDistance(node1.Position, node2.Position)

			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("full_%s_%s", node1.ID, node2.ID)),
				Name:        fmt.Sprintf("连接: %s <-> %s", node1.Name, node2.Name),
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

// GenerateTreeStructure 生成树状结构（最小生成树�?
func (s *pathGenerationService) GenerateTreeStructure(ctx context.Context, rootNodeID domain.NodeID) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	// 找到根节�?
	var rootNode *domain.Node
	nodeMap := make(map[domain.NodeID]*domain.Node)
	for i := range nodes {
		nodeMap[nodes[i].ID] = nodes[i]
		if nodes[i].ID == rootNodeID {
			rootNode = nodes[i]
		}
	}

	if rootNode == nil {
		return nil, fmt.Errorf("根节点不存在: %s", rootNodeID)
	}

	// 使用Prim算法生成最小生成树
	visited := make(map[domain.NodeID]bool)
	visited[rootNodeID] = true

	type edge struct {
		from, to domain.NodeID
		weight   float64
	}

	var paths []domain.Path

	for len(visited) < len(nodes) {
		var minEdge *edge

		// 找到连接已访问节点和未访问节点的最短边
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
			break // 无法连接更多节点
		}

		// 添加边到MST
		visited[minEdge.to] = true
		path := domain.Path{
			ID:          domain.PathID(fmt.Sprintf("tree_%s_%s", minEdge.from, minEdge.to)),
			Name:        fmt.Sprintf("树连�? %s -> %s", minEdge.from, minEdge.to),
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

// GenerateNearestNeighborPaths 生成最近邻路径（每个节点连接到最近的N个邻居）
func (s *pathGenerationService) GenerateNearestNeighborPaths(ctx context.Context, maxDistance float64) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	type neighbor struct {
		nodeID   domain.NodeID
		distance float64
	}

	var paths []domain.Path
	pathSet := make(map[string]bool) // 防止重复路径

	for _, node := range nodes {
		var neighbors []neighbor

		// 计算到所有其他节点的距离
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

		// 按距离排�?
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].distance < neighbors[j].distance
		})

		// 连接到最近的邻居（最�?个）
		maxNeighbors := min(3, len(neighbors))
		for i := 0; i < maxNeighbors; i++ {
			neighbor := neighbors[i]

			// 创建唯一的路径标识符（防止重复）
			pathKey := fmt.Sprintf("%s_%s", min(string(node.ID), string(neighbor.nodeID)), max(string(node.ID), string(neighbor.nodeID)))
			if pathSet[pathKey] {
				continue
			}
			pathSet[pathKey] = true

			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("neighbor_%s_%s", node.ID, neighbor.nodeID)),
				Name:        fmt.Sprintf("最近邻: %s <-> %s", node.Name, neighbor.nodeID),
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

// GenerateGridPaths 生成网格状路径（适用于规则排列的节点�?
func (s *pathGenerationService) GenerateGridPaths(ctx context.Context, connectDiagonal bool) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	if len(nodes) == 0 {
		return []domain.Path{}, nil
	}

	// 按位置排序节点，创建网格结构
	sort.Slice(nodes, func(i, j int) bool {
		if math.Abs(nodes[i].Position.Y-nodes[j].Position.Y) < 10 { // 同一�?
			return nodes[i].Position.X < nodes[j].Position.X
		}
		return nodes[i].Position.Y < nodes[j].Position.Y
	})

	var paths []domain.Path
	tolerance := 50.0 // 位置容差

	// 水平连接（同一行的相邻节点�?
	for i := 0; i < len(nodes)-1; i++ {
		current := nodes[i]
		next := nodes[i+1]

		// 检查是否在同一行且相邻
		if math.Abs(current.Position.Y-next.Position.Y) < tolerance {
			distance := calculateDistance(current.Position, next.Position)
			if distance < tolerance*3 { // 相邻判断
				path := domain.Path{
					ID:          domain.PathID(fmt.Sprintf("grid_h_%s_%s", current.ID, next.ID)),
					Name:        fmt.Sprintf("网格水平: %s -> %s", current.Name, next.Name),
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

	// 垂直连接（同一列的相邻节点�?
	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i >= j {
				continue
			}

			// 检查是否在同一�?
			if math.Abs(node1.Position.X-node2.Position.X) < tolerance {
				distance := calculateDistance(node1.Position, node2.Position)
				if distance < tolerance*3 {
					path := domain.Path{
						ID:          domain.PathID(fmt.Sprintf("grid_v_%s_%s", node1.ID, node2.ID)),
						Name:        fmt.Sprintf("网格垂直: %s -> %s", node1.Name, node2.Name),
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

	// 对角线连接（如果启用�?
	if connectDiagonal {
		for i, node1 := range nodes {
			for j, node2 := range nodes {
				if i >= j {
					continue
				}

				distance := calculateDistance(node1.Position, node2.Position)
				dx := math.Abs(node1.Position.X - node2.Position.X)
				dy := math.Abs(node1.Position.Y - node2.Position.Y)

				// 检查是否为对角线（45度角�?
				if math.Abs(dx-dy) < tolerance && distance < tolerance*2 {
					path := domain.Path{
						ID:          domain.PathID(fmt.Sprintf("grid_d_%s_%s", node1.ID, node2.ID)),
						Name:        fmt.Sprintf("网格对角: %s -> %s", node1.Name, node2.Name),
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

// 工具函数

// calculateDistance 计算两点之间的欧几里得距�?
func calculateDistance(pos1, pos2 domain.Position) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	dz := pos1.Z - pos2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// min 返回两个可比较值的最小�?
func min[T ~int | ~string](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// max 返回两个可比较值的最大�?
func max[T ~int | ~string](a, b T) T {
	if a > b {
		return a
	}
	return b
}
