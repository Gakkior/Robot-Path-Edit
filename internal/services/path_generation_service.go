// Package services 路径生成服务实现
//
// 设计参考：
// - Dijkstra最短路径算法
// - Kruskal最小生成树算法
// - A*寻路算法思想
//
// 特点：
// 1. 多种路径生成算法
// 2. 图论算法实现
// 3. 性能优化
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
	// 最短路径生成
	GenerateShortestPaths(ctx context.Context, startNodeID domain.NodeID) ([]domain.Path, error)

	// 完全连通图生成
	GenerateFullConnectivity(ctx context.Context) ([]domain.Path, error)

	// 树状结构生成
	GenerateTreeStructure(ctx context.Context, rootNodeID domain.NodeID) ([]domain.Path, error)

	// 最近邻路径生成
	GenerateNearestNeighborPaths(ctx context.Context, maxNeighbors int) ([]domain.Path, error)

	// 网格路径生成
	GenerateGridPaths(ctx context.Context, enableDiagonal bool) ([]domain.Path, error)
}

// pathGenerationService 路径生成服务实现
type pathGenerationService struct {
	nodeService NodeService
	pathService PathService
}

// NewPathGenerationService 创建新的路径生成服务实例
func NewPathGenerationService(nodeService NodeService, pathService PathService) PathGenerationService {
	return &pathGenerationService{
		nodeService: nodeService,
		pathService: pathService,
	}
}

// GenerateShortestPaths 生成从指定节点到所有其他节点的最短路径(基于Dijkstra算法的简化版)
func (s *pathGenerationService) GenerateShortestPaths(ctx context.Context, startNodeID domain.NodeID) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	if len(nodes) < 2 {
		return []domain.Path{}, nil
	}

	// 检查起始节点是否存在
	var startNode *domain.Node
	for _, node := range nodes {
		if node.ID == startNodeID {
			startNode = node
			break
		}
	}
	if startNode == nil {
		return nil, fmt.Errorf("起始节点不存在: %s", startNodeID)
	}

	// 获取现有路径用于构建图
	req := ListPathsRequest{}
	existingPathsResp, err := s.pathService.ListPaths(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("获取路径列表失败: %v", err)
	}

	// 构建邻接表
	adjacency := make(map[domain.NodeID]map[domain.NodeID]float64)
	for _, node := range nodes {
		adjacency[node.ID] = make(map[domain.NodeID]float64)
	}

	for _, path := range existingPathsResp.Paths {
		adjacency[path.StartNodeID][path.EndNodeID] = path.Weight
		adjacency[path.EndNodeID][path.StartNodeID] = path.Weight // 无向图
	}

	// 初始化距离
	distances := make(map[domain.NodeID]float64)
	previous := make(map[domain.NodeID]domain.NodeID)
	unvisited := make(map[domain.NodeID]bool)

	for _, node := range nodes {
		distances[node.ID] = math.Inf(1)
		unvisited[node.ID] = true
	}
	distances[startNodeID] = 0

	// Dijkstra主循环
	for len(unvisited) > 0 {
		// 找到距离最小的未访问节点
		var currentNodeID domain.NodeID
		minDistance := math.Inf(1)
		for nodeID := range unvisited {
			if distances[nodeID] < minDistance {
				minDistance = distances[nodeID]
				currentNodeID = nodeID
			}
		}

		if math.IsInf(minDistance, 1) {
			break // 无法到达的节点
		}

		delete(unvisited, currentNodeID)

		// 更新邻居节点的距离
		for neighborID, weight := range adjacency[currentNodeID] {
			if _, exists := unvisited[neighborID]; exists {
				altDistance := distances[currentNodeID] + weight
				if altDistance < distances[neighborID] {
					distances[neighborID] = altDistance
					previous[neighborID] = currentNodeID
				}
			}
		}
	}

	// 生成最短路径
	var paths []domain.Path
	for nodeID, prevNodeID := range previous {
		if prevNodeID != "" {
			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("shortest_%s_%s", prevNodeID, nodeID)),
				Name:        fmt.Sprintf("最短路径 %s -> %s", prevNodeID, nodeID),
				StartNodeID: prevNodeID,
				EndNodeID:   nodeID,
				Weight:      distances[nodeID] - distances[prevNodeID],
				Type:        "normal",
				Status:      "active",
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

	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i < j { // 避免重复连接
				distance := s.calculateDistance(node1.Position, node2.Position)
				path := domain.Path{
					ID:          domain.PathID(fmt.Sprintf("full_%s_%s", node1.ID, node2.ID)),
					Name:        fmt.Sprintf("连接: %s <-> %s", node1.Name, node2.Name),
					StartNodeID: node1.ID,
					EndNodeID:   node2.ID,
					Weight:      distance,
					Type:        "normal",
					Status:      "active",
				}
				paths = append(paths, path)
			}
		}
	}

	return paths, nil
}

// GenerateTreeStructure 生成树状结构（最小生成树）
func (s *pathGenerationService) GenerateTreeStructure(ctx context.Context, rootNodeID domain.NodeID) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	// 找到根节点
	var rootNode *domain.Node
	for _, node := range nodes {
		if node.ID == rootNodeID {
			rootNode = node
			break
		}
	}
	if rootNode == nil {
		return nil, fmt.Errorf("根节点不存在: %s", rootNodeID)
	}

	// 创建所有可能的边
	type edge struct {
		from, to domain.NodeID
		weight   float64
	}

	var edges []edge
	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i < j {
				distance := s.calculateDistance(node1.Position, node2.Position)
				edges = append(edges, edge{
					from:   node1.ID,
					to:     node2.ID,
					weight: distance,
				})
			}
		}
	}

	// 按权重排序
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
	})

	// Kruskal算法生成最小生成树
	var paths []domain.Path
	connected := make(map[domain.NodeID]bool)
	connected[rootNodeID] = true

	for _, edge := range edges {
		// 如果这条边连接了已连接和未连接的节点
		fromConnected := connected[edge.from]
		toConnected := connected[edge.to]

		if fromConnected != toConnected { // 一个连接，一个未连接
			path := domain.Path{
				ID:          domain.PathID(fmt.Sprintf("tree_%s_%s", edge.from, edge.to)),
				Name:        fmt.Sprintf("树连接 %s -> %s", edge.from, edge.to),
				StartNodeID: edge.from,
				EndNodeID:   edge.to,
				Weight:      edge.weight,
				Type:        "normal",
				Status:      "active",
			}
			paths = append(paths, path)

			// 标记为已连接
			connected[edge.from] = true
			connected[edge.to] = true

			// 如果所有节点都已连接，结束
			if len(connected) == len(nodes) {
				break
			}
		}
	}

	return paths, nil
}

// GenerateNearestNeighborPaths 生成最近邻路径（每个节点连接到最近的N个邻居）
func (s *pathGenerationService) GenerateNearestNeighborPaths(ctx context.Context, maxNeighbors int) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	if maxNeighbors <= 0 {
		maxNeighbors = 3 // 默认连接到最近的3个邻居
	}

	var paths []domain.Path
	pathSet := make(map[string]bool) // 防止重复路径

	for _, node := range nodes {
		// 计算到其他所有节点的距离
		type neighbor struct {
			nodeID   domain.NodeID
			distance float64
		}

		var neighbors []neighbor
		for _, otherNode := range nodes {
			if node.ID != otherNode.ID {
				distance := s.calculateDistance(node.Position, otherNode.Position)
				neighbors = append(neighbors, neighbor{
					nodeID:   otherNode.ID,
					distance: distance,
				})
			}
		}

		// 按距离排序
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].distance < neighbors[j].distance
		})

		// 连接到最近的邻居（最多maxNeighbors个）
		maxConnections := len(neighbors)
		if maxConnections > maxNeighbors {
			maxConnections = maxNeighbors
		}

		for i := 0; i < maxConnections; i++ {
			neighbor := neighbors[i]

			// 创建唯一的路径标识符（防止重复）
			pathKey1 := fmt.Sprintf("%s_%s", node.ID, neighbor.nodeID)
			pathKey2 := fmt.Sprintf("%s_%s", neighbor.nodeID, node.ID)

			if !pathSet[pathKey1] && !pathSet[pathKey2] {
				path := domain.Path{
					ID:          domain.PathID(fmt.Sprintf("neighbor_%s_%s", node.ID, neighbor.nodeID)),
					Name:        fmt.Sprintf("最近邻: %s <-> %s", node.Name, neighbor.nodeID),
					StartNodeID: node.ID,
					EndNodeID:   neighbor.nodeID,
					Weight:      neighbor.distance,
					Type:        "normal",
					Status:      "active",
				}
				paths = append(paths, path)
				pathSet[pathKey1] = true
				pathSet[pathKey2] = true
			}
		}
	}

	return paths, nil
}

// GenerateGridPaths 生成网格状路径（适用于规则排列的节点）
func (s *pathGenerationService) GenerateGridPaths(ctx context.Context, enableDiagonal bool) ([]domain.Path, error) {
	nodes, err := s.nodeService.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	if len(nodes) < 2 {
		return []domain.Path{}, nil
	}

	// 按位置排序节点，创建网格结构
	// 简化实现：假设节点按网格排列
	// 网格路径生成逻辑
	// 实际实现需要更复杂的逻辑来检测网格结构

	var paths []domain.Path
	tolerance := 50.0 // 位置容差

	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i >= j {
				continue
			}

			dx := math.Abs(node1.Position.X - node2.Position.X)
			dy := math.Abs(node1.Position.Y - node2.Position.Y)

			isHorizontal := dy < tolerance && dx > tolerance
			isVertical := dx < tolerance && dy > tolerance
			isDiagonal := enableDiagonal && math.Abs(dx-dy) < tolerance

			if isHorizontal || isVertical || isDiagonal {
				distance := s.calculateDistance(node1.Position, node2.Position)

				pathType := "网格"
				if isHorizontal {
					pathType = "水平"
				} else if isVertical {
					pathType = "垂直"
				} else if isDiagonal {
					pathType = "对角"
				}

				path := domain.Path{
					ID:          domain.PathID(fmt.Sprintf("grid_%s_%s", node1.ID, node2.ID)),
					Name:        fmt.Sprintf("%s连接: %s <-> %s", pathType, node1.Name, node2.Name),
					StartNodeID: node1.ID,
					EndNodeID:   node2.ID,
					Weight:      distance,
					Type:        "normal",
					Status:      "active",
				}
				paths = append(paths, path)
			}
		}
	}

	return paths, nil
}

// calculateDistance 计算两点之间的欧几里得距离
func (s *pathGenerationService) calculateDistance(pos1, pos2 domain.Position) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	dz := pos1.Z - pos2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
