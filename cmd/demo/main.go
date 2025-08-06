// 机器人路径编辑器演示版本
// 使用内存存储，无需数据库依赖
package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"robot-path-editor/internal/domain"
)

// 简化的内存存储
type MemoryStore struct {
	nodes map[string]*domain.Node
	paths map[string]*domain.Path
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		nodes: make(map[string]*domain.Node),
		paths: make(map[string]*domain.Path),
	}

	// 添加一些示例数据
	store.addSampleData()
	return store
}

func (s *MemoryStore) addSampleData() {
	// 创建示例节点
	node1 := domain.NewNode("起始点", "point")
	node1.Position = domain.Position{X: 100, Y: 100, Z: 0}

	node2 := domain.NewNode("中转点", "waypoint")
	node2.Position = domain.Position{X: 300, Y: 200, Z: 0}

	node3 := domain.NewNode("目标点", "point")
	node3.Position = domain.Position{X: 500, Y: 300, Z: 0}

	s.nodes[string(node1.ID)] = node1
	s.nodes[string(node2.ID)] = node2
	s.nodes[string(node3.ID)] = node3

	// 创建示例路径
	path1 := domain.NewPath("路径1", node1.ID, node2.ID)
	path2 := domain.NewPath("路径2", node2.ID, node3.ID)

	s.paths[string(path1.ID)] = path1
	s.paths[string(path2.ID)] = path2

	logrus.Info("已加载示例数据：3个节点，2条路径")
}

// API处理器
type DemoHandlers struct {
	store *MemoryStore
}

func NewDemoHandlers(store *MemoryStore) *DemoHandlers {
	return &DemoHandlers{store: store}
}

// 节点相关API
func (h *DemoHandlers) ListNodes(c *gin.Context) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range h.store.nodes {
		nodes = append(nodes, node)
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
		"count": len(nodes),
	})
}

func (h *DemoHandlers) CreateNode(c *gin.Context) {
	var req struct {
		Name     string          `json:"name" binding:"required"`
		Type     string          `json:"type,omitempty"`
		Position domain.Position `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 默认节点类型
	nodeType := req.Type
	if nodeType == "" {
		nodeType = "point"
	}

	node := domain.NewNode(req.Name, nodeType)
	node.Position = req.Position

	h.store.mu.Lock()
	h.store.nodes[string(node.ID)] = node
	h.store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "节点创建成功",
		"node":    node,
	})
}

func (h *DemoHandlers) GetNode(c *gin.Context) {
	id := c.Param("id")

	h.store.mu.RLock()
	node, exists := h.store.nodes[id]
	h.store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node})
}

// GetPath 获取单个路径
func (h *DemoHandlers) GetPath(c *gin.Context) {
	id := c.Param("id")

	h.store.mu.RLock()
	path, exists := h.store.paths[id]
	h.store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "路径不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": path})
}

func (h *DemoHandlers) UpdateNode(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name     *string          `json:"name"`
		Position *domain.Position `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	node, exists := h.store.nodes[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}

	if req.Name != nil {
		node.Name = *req.Name
	}
	if req.Position != nil {
		node.Position = *req.Position
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "节点更新成功",
		"node":    node,
	})
}

func (h *DemoHandlers) DeleteNode(c *gin.Context) {
	id := c.Param("id")

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, exists := h.store.nodes[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}

	delete(h.store.nodes, id)

	c.JSON(http.StatusOK, gin.H{"message": "节点删除成功"})
}

// UpdateNodePosition 更新节点位置（仅坐标）
func (h *DemoHandlers) UpdateNodePosition(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		X float64 `json:"x" binding:"required"`
		Y float64 `json:"y" binding:"required"`
		Z float64 `json:"z"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.store.mu.Lock()
	node, exists := h.store.nodes[id]
	if !exists {
		h.store.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}
	node.Position.X = req.X
	node.Position.Y = req.Y
	node.Position.Z = req.Z
	h.store.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "位置已更新", "node": node})
}

// 路径相关API
func (h *DemoHandlers) ListPaths(c *gin.Context) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	var paths []*domain.Path
	for _, path := range h.store.paths {
		paths = append(paths, path)
	}

	c.JSON(http.StatusOK, gin.H{
		"paths": paths,
		"count": len(paths),
	})
}

func (h *DemoHandlers) CreatePath(c *gin.Context) {
	var req struct {
		Name        string        `json:"name" binding:"required"`
		StartNodeID domain.NodeID `json:"start_node_id" binding:"required"`
		EndNodeID   domain.NodeID `json:"end_node_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查节点是否存在
	h.store.mu.RLock()
	_, startExists := h.store.nodes[string(req.StartNodeID)]
	_, endExists := h.store.nodes[string(req.EndNodeID)]
	h.store.mu.RUnlock()

	if !startExists || !endExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "起始节点或结束节点不存在"})
		return
	}

	path := domain.NewPath(req.Name, req.StartNodeID, req.EndNodeID)

	h.store.mu.Lock()
	h.store.paths[string(path.ID)] = path
	h.store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "路径创建成功",
		"path":    path,
	})
}

// DeletePath 删除路径
func (h *DemoHandlers) DeletePath(c *gin.Context) {
	id := c.Param("id")
	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	if _, ok := h.store.paths[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "路径不存在"})
		return
	}
	delete(h.store.paths, id)
	c.JSON(http.StatusOK, gin.H{"message": "路径已删除"})
}

// UpdatePath 更新路径
func (h *DemoHandlers) UpdatePath(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Status      string `json:"status"`
		StartNodeID string `json:"start_node_id"`
		EndNodeID   string `json:"end_node_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误: " + err.Error()})
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	path, exists := h.store.paths[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "路径不存在"})
		return
	}

	// 验证起始和结束节点存在
	if req.StartNodeID != "" {
		if _, exists := h.store.nodes[req.StartNodeID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "起始节点不存在"})
			return
		}
	}

	if req.EndNodeID != "" {
		if _, exists := h.store.nodes[req.EndNodeID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "结束节点不存在"})
			return
		}
	}

	// 更新路径属性
	if req.Name != "" {
		path.Name = req.Name
	}
	if req.Type != "" {
		path.Type = domain.PathType(req.Type)
	}
	if req.Status != "" {
		path.Status = domain.PathStatus(req.Status)
	}
	if req.StartNodeID != "" {
		path.StartNodeID = domain.NodeID(req.StartNodeID)
	}
	if req.EndNodeID != "" {
		path.EndNodeID = domain.NodeID(req.EndNodeID)
	}

	path.Metadata.UpdatedAt = time.Now()
	h.store.paths[id] = path

	c.JSON(http.StatusOK, gin.H{
		"message": "路径更新成功",
		"path":    path,
	})
}

// ApplyLayout 应用布局算法
func (h *DemoHandlers) ApplyLayout(c *gin.Context) {
	var req struct {
		Algorithm string `json:"algorithm" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误: " + err.Error()})
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	// 获取所有节点和路径
	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := make([]domain.Path, 0, len(h.store.paths))
	for _, path := range h.store.paths {
		paths = append(paths, *path)
	}

	// 应用布局算法
	var updatedNodes []domain.Node
	switch req.Algorithm {
	case "grid":
		updatedNodes = applyGridLayout(nodes, 120.0)
	case "force-directed":
		updatedNodes = applyForceDirectedLayout(nodes, paths, 50)
	case "circular":
		updatedNodes = applyCircularLayout(nodes, 250.0, 500.0, 400.0)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的布局算法: " + req.Algorithm})
		return
	}

	// 更新存储中的节点位置
	for _, node := range updatedNodes {
		if existingNode, ok := h.store.nodes[string(node.ID)]; ok {
			existingNode.Position = node.Position
			h.store.nodes[string(node.ID)] = existingNode
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "布局应用成功",
		"algorithm":      req.Algorithm,
		"affected_nodes": len(updatedNodes),
	})
}

// GenerateNearestNeighborPaths 生成最近邻路径
func (h *DemoHandlers) GenerateNearestNeighborPaths(c *gin.Context) {
	var req struct {
		MaxDistance float64 `json:"max_distance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.MaxDistance = 200.0 // 默认值
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := generateNearestNeighborPaths(nodes, req.MaxDistance)

	// 添加到存储
	createdCount := 0
	for _, path := range paths {
		if _, exists := h.store.paths[string(path.ID)]; !exists {
			h.store.paths[string(path.ID)] = &path
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "最近邻路径生成成功",
		"created_paths": createdCount,
		"max_distance":  req.MaxDistance,
	})
}

// GenerateFullConnectivity 生成完全连通路径
func (h *DemoHandlers) GenerateFullConnectivity(c *gin.Context) {
	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := generateFullConnectivityPaths(nodes)

	// 添加到存储
	createdCount := 0
	for _, path := range paths {
		if _, exists := h.store.paths[string(path.ID)]; !exists {
			h.store.paths[string(path.ID)] = &path
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "完全连通路径生成成功",
		"created_paths": createdCount,
		"total_nodes":   len(nodes),
	})
}

// GenerateGridPaths 生成网格路径
func (h *DemoHandlers) GenerateGridPaths(c *gin.Context) {
	var req struct {
		ConnectDiagonal bool `json:"connect_diagonal"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.ConnectDiagonal = false // 默认值
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := generateGridPaths(nodes, req.ConnectDiagonal)

	// 添加到存储
	createdCount := 0
	for _, path := range paths {
		if _, exists := h.store.paths[string(path.ID)]; !exists {
			h.store.paths[string(path.ID)] = &path
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "网格路径生成成功",
		"created_paths":    createdCount,
		"connect_diagonal": req.ConnectDiagonal,
	})
}

// 布局算法实现

// applyGridLayout 网格布局
func applyGridLayout(nodes []domain.Node, spacing float64) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	cols := int(math.Ceil(math.Sqrt(float64(len(nodes)))))
	updatedNodes := make([]domain.Node, len(nodes))

	for i, node := range nodes {
		row := i / cols
		col := i % cols

		updatedNode := node
		updatedNode.Position.X = float64(col)*spacing + 100
		updatedNode.Position.Y = float64(row)*spacing + 100
		updatedNodes[i] = updatedNode
	}

	return updatedNodes
}

// applyForceDirectedLayout 力导向布局 (简化版)
func applyForceDirectedLayout(nodes []domain.Node, paths []domain.Path, iterations int) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	// 初始化随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 参数设置
	width, height := 1000.0, 800.0
	k := math.Sqrt((width * height) / float64(len(nodes)))

	// 初始化节点位置
	updatedNodes := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		updatedNode := node
		if node.Position.X == 0 && node.Position.Y == 0 {
			updatedNode.Position.X = r.Float64() * width
			updatedNode.Position.Y = r.Float64() * height
		}
		updatedNodes[i] = updatedNode
	}

	// 迭代计算
	for iter := 0; iter < iterations; iter++ {
		forces := make(map[string]struct{ fx, fy float64 })

		// 初始化力
		for i := range updatedNodes {
			forces[string(updatedNodes[i].ID)] = struct{ fx, fy float64 }{0, 0}
		}

		// 计算排斥力
		for i := 0; i < len(updatedNodes); i++ {
			for j := i + 1; j < len(updatedNodes); j++ {
				node1, node2 := &updatedNodes[i], &updatedNodes[j]
				dx := node1.Position.X - node2.Position.X
				dy := node1.Position.Y - node2.Position.Y
				distance := math.Max(math.Sqrt(dx*dx+dy*dy), 0.01)

				repulsiveForce := k * k / distance
				fx := repulsiveForce * dx / distance
				fy := repulsiveForce * dy / distance

				force1 := forces[string(node1.ID)]
				force1.fx += fx
				force1.fy += fy
				forces[string(node1.ID)] = force1

				force2 := forces[string(node2.ID)]
				force2.fx -= fx
				force2.fy -= fy
				forces[string(node2.ID)] = force2
			}
		}

		// 计算吸引力
		for _, path := range paths {
			var node1, node2 *domain.Node
			for i := range updatedNodes {
				if updatedNodes[i].ID == path.StartNodeID {
					node1 = &updatedNodes[i]
				}
				if updatedNodes[i].ID == path.EndNodeID {
					node2 = &updatedNodes[i]
				}
			}

			if node1 != nil && node2 != nil {
				dx := node2.Position.X - node1.Position.X
				dy := node2.Position.Y - node1.Position.Y
				distance := math.Max(math.Sqrt(dx*dx+dy*dy), 0.01)

				attractiveForce := distance * distance / k
				fx := attractiveForce * dx / distance
				fy := attractiveForce * dy / distance

				force1 := forces[string(node1.ID)]
				force1.fx += fx
				force1.fy += fy
				forces[string(node1.ID)] = force1

				force2 := forces[string(node2.ID)]
				force2.fx -= fx
				force2.fy -= fy
				forces[string(node2.ID)] = force2
			}
		}

		// 应用力
		temperature := 10.0 * (1.0 - float64(iter)/float64(iterations))
		for i := range updatedNodes {
			force := forces[string(updatedNodes[i].ID)]
			displacement := math.Min(math.Sqrt(force.fx*force.fx+force.fy*force.fy), temperature)

			if displacement > 0.01 {
				updatedNodes[i].Position.X += force.fx / displacement * temperature
				updatedNodes[i].Position.Y += force.fy / displacement * temperature
			}

			// 保持在画布范围内
			updatedNodes[i].Position.X = math.Max(50, math.Min(width-50, updatedNodes[i].Position.X))
			updatedNodes[i].Position.Y = math.Max(50, math.Min(height-50, updatedNodes[i].Position.Y))
		}
	}

	return updatedNodes
}

// applyCircularLayout 圆形布局
func applyCircularLayout(nodes []domain.Node, radius, centerX, centerY float64) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	updatedNodes := make([]domain.Node, len(nodes))
	angleStep := 2 * math.Pi / float64(len(nodes))

	for i, node := range nodes {
		angle := float64(i) * angleStep
		updatedNode := node
		updatedNode.Position.X = centerX + radius*math.Cos(angle)
		updatedNode.Position.Y = centerY + radius*math.Sin(angle)
		updatedNodes[i] = updatedNode
	}

	return updatedNodes
}

// 路径生成算法实现

// generateNearestNeighborPaths 生成最近邻路径
func generateNearestNeighborPaths(nodes []domain.Node, maxDistance float64) []domain.Path {
	if len(nodes) < 2 {
		return []domain.Path{}
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

		// 按距离排序
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].distance < neighbors[j].distance
		})

		// 连接到最近的邻居（最多3个）
		maxNeighbors := minInt(3, len(neighbors))
		for i := 0; i < maxNeighbors; i++ {
			neighbor := neighbors[i]

			// 创建唯一的路径标识符（防止重复）
			pathKey := fmt.Sprintf("%s_%s", minString(string(node.ID), string(neighbor.nodeID)), maxString(string(node.ID), string(neighbor.nodeID)))
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

	return paths
}

// generateFullConnectivityPaths 生成完全连通路径
func generateFullConnectivityPaths(nodes []domain.Node) []domain.Path {
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

	return paths
}

// generateGridPaths 生成网格路径
func generateGridPaths(nodes []domain.Node, connectDiagonal bool) []domain.Path {
	if len(nodes) == 0 {
		return []domain.Path{}
	}

	// 按位置排序节点，创建网格结构
	sort.Slice(nodes, func(i, j int) bool {
		if math.Abs(nodes[i].Position.Y-nodes[j].Position.Y) < 10 { // 同一行
			return nodes[i].Position.X < nodes[j].Position.X
		}
		return nodes[i].Position.Y < nodes[j].Position.Y
	})

	var paths []domain.Path
	tolerance := 50.0 // 位置容差

	// 水平连接（同一行的相邻节点）
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

	// 垂直连接（同一列的相邻节点）
	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if i >= j {
				continue
			}

			// 检查是否在同一列
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

	// 对角线连接（如果启用）
	if connectDiagonal {
		for i, node1 := range nodes {
			for j, node2 := range nodes {
				if i >= j {
					continue
				}

				distance := calculateDistance(node1.Position, node2.Position)
				dx := math.Abs(node1.Position.X - node2.Position.X)
				dy := math.Abs(node1.Position.Y - node2.Position.Y)

				// 检查是否为对角线（45度角）
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

	return paths
}

// 工具函数

// calculateDistance 计算两点之间的欧几里得距离
func calculateDistance(pos1, pos2 domain.Position) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	dz := pos1.Z - pos2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// minInt 返回两个整数的最小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// minString 返回两个字符串的字典序最小值
func minString(a, b string) string {
	if a < b {
		return a
	}
	return b
}

// maxString 返回两个字符串的字典序最大值
func maxString(a, b string) string {
	if a > b {
		return a
	}
	return b
}

// 健康检查
func (h *DemoHandlers) HealthCheck(c *gin.Context) {
	h.store.mu.RLock()
	nodeCount := len(h.store.nodes)
	pathCount := len(h.store.paths)
	h.store.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "robot-path-editor-demo",
		"storage": "memory",
		"data": gin.H{
			"nodes": nodeCount,
			"paths": pathCount,
		},
	})
}

// 获取画布数据
func (h *DemoHandlers) GetCanvasData(c *gin.Context) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	// 准备画布数据
	canvasData := gin.H{
		"nodes": h.store.nodes,
		"paths": h.store.paths,
		"canvas": gin.H{
			"width":  1920,
			"height": 1080,
			"zoom":   1.0,
		},
	}

	c.JSON(http.StatusOK, canvasData)
}

// 表格视图HTML
const tableHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>表格视图 - 机器人路径编辑器</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f8f9fa;
            min-height: 100vh;
        }
        .header {
            background: #2c3e50;
            color: white;
            padding: 1rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .header h1 {
            font-size: 1.5rem;
            font-weight: 500;
        }
        .nav-links a {
            color: white;
            text-decoration: none;
            margin-left: 1rem;
            padding: 0.5rem 1rem;
            border-radius: 4px;
            transition: background 0.2s;
        }
        .nav-links a:hover {
            background: rgba(255,255,255,0.1);
        }
        .controls {
            padding: 1rem 2rem;
            background: white;
            border-bottom: 1px solid #dee2e6;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .view-toggle {
            display: flex;
            gap: 0.5rem;
        }
        .btn {
            padding: 0.5rem 1rem;
            border: 1px solid #dee2e6;
            background: white;
            color: #495057;
            border-radius: 4px;
            cursor: pointer;
            transition: all 0.2s;
        }
        .btn:hover {
            background: #f8f9fa;
        }
        .btn.active {
            background: #007bff;
            color: white;
            border-color: #007bff;
        }
        .btn-primary {
            background: #007bff;
            color: white;
            border-color: #007bff;
        }
        .btn-primary:hover {
            background: #0056b3;
        }
        .btn-small {
            padding: 0.25rem 0.5rem;
            font-size: 0.875rem;
            margin: 0 0.125rem;
        }
        .btn-save {
            background: #28a745;
            color: white;
            border-color: #28a745;
        }
        .btn-delete {
            background: #dc3545;
            color: white;
            border-color: #dc3545;
        }
        .main-content {
            padding: 2rem;
        }
        .table-wrapper {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow-x: auto;
        }
        .table-header {
            padding: 1rem 1.5rem;
            border-bottom: 1px solid #dee2e6;
        }
        .table-header h3 {
            color: #495057;
            font-weight: 500;
        }
        .data-table {
            width: 100%;
            border-collapse: collapse;
        }
        .data-table th,
        .data-table td {
            padding: 0.75rem;
            text-align: left;
            border-bottom: 1px solid #dee2e6;
        }
        .data-table th {
            background: #f8f9fa;
            font-weight: 500;
            color: #495057;
        }
        .data-table tr:hover {
            background: #f8f9fa;
        }
        .editable-input {
            border: 1px solid #ced4da;
            border-radius: 4px;
            padding: 0.25rem 0.5rem;
            width: 100%;
            font-size: 0.875rem;
        }
        .editable-input:focus {
            outline: none;
            border-color: #80bdff;
            box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
        }
        .message {
            animation: slideIn 0.3s ease-out;
        }
        @keyframes slideIn {
            from { transform: translateX(100%); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>📊 表格视图</h1>
        <div class="nav-links">
            <a href="/">画布视图</a>
            <a href="/table">表格视图</a>
        </div>
    </div>
    
    <div class="controls">
        <div class="view-toggle">
            <button id="nodeViewBtn" class="btn active">节点</button>
            <button id="pathViewBtn" class="btn">路径</button>
        </div>
        <div class="actions">
            <button id="refreshBtn" class="btn">🔄 刷新</button>
            			<button id="addBtn" class="btn btn-primary">+ 添加</button>
        </div>
    </div>
    
    <div class="main-content">
        <div id="tableContainer">
            <!-- 表格内容将在此处动态生成 -->
        </div>
    </div>
    
    <script src="/static/table.js"></script>
</body>
</html>`

// 主页面HTML
const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>机器人路径编辑器 - 演示版</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            text-align: center;
            padding: 2rem;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            backdrop-filter: blur(10px);
            max-width: 800px;
        }
        .logo { font-size: 4rem; margin-bottom: 1rem; }
        h1 { color: #2c3e50; margin-bottom: 1rem; font-size: 2.5rem; font-weight: 300; }
        .subtitle { color: #7f8c8d; margin-bottom: 2rem; font-size: 1.2rem; }
        .demo-badge {
            display: inline-block;
            padding: 0.5rem 1rem;
            background: #e74c3c;
            color: white;
            border-radius: 25px;
            font-weight: 500;
            margin: 1rem 0;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.05); }
            100% { transform: scale(1); }
        }
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin: 2rem 0;
        }
        .feature {
            padding: 1rem;
            background: rgba(52, 152, 219, 0.1);
            border-radius: 10px;
            border: 1px solid rgba(52, 152, 219, 0.2);
        }
        .feature h3 { color: #3498db; margin-bottom: 0.5rem; }
        .api-section {
            margin-top: 2rem;
            padding: 1rem;
            background: rgba(241, 196, 15, 0.1);
            border-radius: 10px;
            border: 1px solid rgba(241, 196, 15, 0.2);
        }
        .api-section h3 { color: #f39c12; margin-bottom: 1rem; }
        .api-endpoints {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 0.5rem;
            text-align: left;
        }
        .api-endpoint {
            font-family: 'Monaco', 'Menlo', monospace;
            background: rgba(0, 0, 0, 0.05);
            padding: 0.5rem;
            border-radius: 4px;
            font-size: 0.9rem;
        }
        .stats {
            margin: 2rem 0;
            padding: 1rem;
            background: rgba(46, 204, 113, 0.1);
            border-radius: 10px;
            border: 1px solid rgba(46, 204, 113, 0.2);
        }
        .stats h3 { color: #27ae60; margin-bottom: 0.5rem; }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
            gap: 1rem;
        }
        .stat-item {
            text-align: center;
        }
        .stat-number {
            font-size: 2rem;
            font-weight: bold;
            color: #27ae60;
        }
        .stat-label {
            color: #7f8c8d;
            font-size: 0.9rem;
        }
            .sidebar {
            position: fixed;
            right: 0;
            top: 0;
            width: 300px;
            height: 100vh;
            background: #ecf0f1;
            box-shadow: -2px 0 8px rgba(0,0,0,0.1);
            padding: 1rem;
            overflow-y: auto;
        }
        .sidebar h2 { font-size: 1.2rem; margin-bottom: 1rem; color:#2c3e50; }
        .form-group { margin-bottom: 1rem; text-align:left; }
        .form-group label { display:block; font-size:0.9rem; color:#7f8c8d; margin-bottom:0.25rem; }
        .form-group input { width:100%; padding:0.4rem; border:1px solid #bdc3c7; border-radius:4px; }
        .btn { display:inline-block; background:#3498db; color:#fff; padding:0.4rem 0.8rem; border:none; border-radius:4px; cursor:pointer; }
        .btn:hover { background:#2980b9; }
    </style>
</head>
<body>
    <div id="toolbar" style="position:fixed;top:10px;left:10px;z-index:1000;display:flex;gap:10px;">
        <button id="undoBtn" style="padding:8px 12px;background:#3498db;color:white;border:none;border-radius:4px;cursor:pointer;" disabled title="撤销 (Ctrl+Z)">↶ 撤销</button>
        <button id="redoBtn" style="padding:8px 12px;background:#3498db;color:white;border:none;border-radius:4px;cursor:pointer;" disabled title="重做 (Ctrl+Y)">↷ 重做</button>
        <div style="border-left:1px solid rgba(255,255,255,0.3);margin:0 10px;"></div>
        <button id="gridLayoutBtn" style="padding:8px 12px;background:#27ae60;color:white;border:none;border-radius:4px;cursor:pointer;" title="网格布局">🔳 网格</button>
        <button id="forceLayoutBtn" style="padding:8px 12px;background:#e67e22;color:white;border:none;border-radius:4px;cursor:pointer;" title="力导向布局">⚡ 力导向</button>
        <button id="circularLayoutBtn" style="padding:8px 12px;background:#9b59b6;color:white;border:none;border-radius:4px;cursor:pointer;" title="圆形布局">⭕ 圆形</button>
        <div style="border-left:1px solid rgba(255,255,255,0.3);margin:0 10px;"></div>
        <button id="nearestPathBtn" style="padding:8px 12px;background:#f39c12;color:white;border:none;border-radius:4px;cursor:pointer;" title="生成最近邻路径">🔗 最近邻</button>
        <button id="fullConnectBtn" style="padding:8px 12px;background:#e74c3c;color:white;border:none;border-radius:4px;cursor:pointer;" title="生成完全连通">🕸️ 全连通</button>
        <button id="gridPathBtn" style="padding:8px 12px;background:#8e44ad;color:white;border:none;border-radius:4px;cursor:pointer;" title="生成网格路径">📐 网格路径</button>
    </div>
    <div id="canvas-container" style="position:fixed;left:0;top:0;width:calc(100% - 300px);height:100vh;"></div>
    <div class="container">
        <div class="logo">🤖</div>
        <h1>机器人路径编辑器</h1>
        <p class="subtitle">现代化的三端兼容路径管理工具</p>
        
        <div class="demo-badge">🚀 演示版本 - 内存存储模式</div>
        
        <div class="stats">
            <h3>📊 实时数据统计</h3>
            <div class="stats-grid">
                <div class="stat-item">
                    <div class="stat-number" id="nodeCount">-</div>
                    <div class="stat-label">节点数量</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number" id="pathCount">-</div>
                    <div class="stat-label">路径数量</div>
                </div>
            </div>
        </div>
        
        <div class="features">
            <div class="feature">
                <h3>🗄️ 内存存储</h3>
                <p>无需数据库，快速启动演示</p>
            </div>
            <div class="feature">
                <h3>🎨 RESTful API</h3>
                <p>完整的节点和路径管理接口</p>
            </div>
            <div class="feature">
                <h3>📱 响应式设计</h3>
                <p>自适应不同设备屏幕</p>
            </div>
            <div class="feature">
                <h3>⚡ 实时更新</h3>
                <p>数据变化实时同步显示</p>
            </div>
        </div>
        
        <div class="api-section">
            <h3>🔗 快速导航</h3>
            <div class="nav-buttons" style="display:flex;gap:1rem;margin:1rem 0;">
                <a href="/" style="text-decoration:none;">
                    <button class="btn" style="width:100%;">🎨 画布视图</button>
                </a>
                <a href="/table" style="text-decoration:none;">
                    <button class="btn" style="width:100%;">📊 表格视图</button>
                </a>
            </div>
            
            <h3>🔌 API 端点</h3>
            <div class="api-endpoints">
                <div class="api-endpoint">GET /api/v1/nodes</div>
                <div class="api-endpoint">POST /api/v1/nodes</div>
                <div class="api-endpoint">GET /api/v1/paths</div>
                <div class="api-endpoint">POST /api/v1/paths</div>
                <div class="api-endpoint">GET /canvas-data</div>
                <div class="api-endpoint">GET /health</div>
            </div>
        </div>
    </div>
    
    <script src="https://unpkg.com/konva@9.3.3/konva.min.js"></script>
    <script src="/static/canvas.js"></script>
    <script>
        // 实时更新统计数据
        function updateStats() {
            fetch('/health')
                .then(response => response.json())
                .then(data => {
                    if (data.data) {
                        document.getElementById('nodeCount').textContent = data.data.nodes || 0;
                        document.getElementById('pathCount').textContent = data.data.paths || 0;
                    }
                })
                .catch(error => {
                    console.error('获取统计数据失败:', error);
                });
        }
        
        // 页面加载时更新一次
        updateStats();
        
        // 5秒更新一次
        setInterval(updateStats, 5000);
        
        // 简单的API测试
        console.log('🤖 机器人路径编辑器演示版已启动');
        console.log('📡 API测试:');
        
        fetch('/api/v1/nodes')
            .then(response => response.json())
            .then(data => {
                console.log('✅ 节点API测试成功:', data);
            })
            .catch(error => {
                console.error('❌ 节点API测试失败:', error);
            });
    </script>
</body>
</html>`

func main() {
	// 设置日志
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	fmt.Println("🤖 机器人路径编辑器演示版启动中...")

	// 初始化内存存储
	store := NewMemoryStore()
	handlers := NewDemoHandlers(store)

	// 设置Gin路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Static("/static", "./web/static")
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 主页面
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})

	// 表格视图页面
	r.GET("/table", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(tableHTML))
	})

	// 健康检查
	r.GET("/health", handlers.HealthCheck)

	// 画布数据
	r.GET("/canvas-data", handlers.GetCanvasData)

	// API路由
	api := r.Group("/api/v1")
	{
		// 节点管理
		nodes := api.Group("/nodes")
		{
			nodes.GET("", handlers.ListNodes)
			nodes.POST("", handlers.CreateNode)
			nodes.GET("/:id", handlers.GetNode)
			nodes.PUT("/:id", handlers.UpdateNode)
			nodes.PUT("/:id/position", handlers.UpdateNodePosition)
			nodes.DELETE("/:id", handlers.DeleteNode)
		}

		// 路径管理
		paths := api.Group("/paths")
		{
			paths.GET("", handlers.ListPaths)
			paths.GET("/:id", handlers.GetPath)
			paths.POST("", handlers.CreatePath)
			paths.PUT("/:id", handlers.UpdatePath)
			paths.DELETE("/:id", handlers.DeletePath)
		}

		// 布局算法端点
		layout := api.Group("/layout")
		{
			layout.POST("/apply", handlers.ApplyLayout)
		}

		// 路径生成端点
		pathGen := api.Group("/path-generation")
		{
			pathGen.POST("/nearest-neighbor", handlers.GenerateNearestNeighborPaths)
			pathGen.POST("/full-connectivity", handlers.GenerateFullConnectivity)
			pathGen.POST("/grid", handlers.GenerateGridPaths)
		}
	}

	// 启动服务器
	port := ":8080"
	fmt.Printf("🚀 演示服务器启动成功！访问地址: http://localhost%s\n", port)
	fmt.Println("📊 API端点:")
	fmt.Println("  - GET  /health        健康检查")
	fmt.Println("  - GET  /canvas-data   画布数据")
	fmt.Println("  - GET  /api/v1/nodes  节点列表")
	fmt.Println("  - POST /api/v1/nodes  创建节点")
	fmt.Println("  - GET  /api/v1/paths  路径列表")
	fmt.Println("  - POST /api/v1/paths  创建路径")

	if err := r.Run(port); err != nil {
		logrus.WithError(err).Fatal("服务器启动失败")
	}
}
