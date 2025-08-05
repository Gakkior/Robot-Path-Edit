// 鏈哄櫒浜鸿矾寰勭紪杈戝櫒婕旂ず鐗堟湰
// 浣跨敤鍐呭瓨瀛樺偍锛屾棤闇€鏁版嵁搴撲緷璧?
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

// 绠€鍖栫殑鍐呭瓨瀛樺偍
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

	// 娣诲姞涓€浜涚ず渚嬫暟鎹?
	store.addSampleData()
	return store
}

func (s *MemoryStore) addSampleData() {
	// 鍒涘缓绀轰緥鑺傜偣
	node1 := domain.NewNode("璧峰鐐?, domain.Position{X: 100, Y: 100, Z: 0})
	node2 := domain.NewNode("涓浆鐐?, domain.Position{X: 300, Y: 200, Z: 0})
	node3 := domain.NewNode("鐩爣鐐?, domain.Position{X: 500, Y: 300, Z: 0})

	s.nodes[string(node1.ID)] = node1
	s.nodes[string(node2.ID)] = node2
	s.nodes[string(node3.ID)] = node3

	// 鍒涘缓绀轰緥璺緞
	path1 := domain.NewPath("璺緞1", node1.ID, node2.ID)
	path2 := domain.NewPath("璺緞2", node2.ID, node3.ID)

	s.paths[string(path1.ID)] = path1
	s.paths[string(path2.ID)] = path2

	logrus.Info("宸插姞杞界ず渚嬫暟鎹細3涓妭鐐癸紝2鏉¤矾寰?)
}

// API澶勭悊鍣?
type DemoHandlers struct {
	store *MemoryStore
}

func NewDemoHandlers(store *MemoryStore) *DemoHandlers {
	return &DemoHandlers{store: store}
}

// 鑺傜偣鐩稿叧API
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
		Position domain.Position `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node := domain.NewNode(req.Name, req.Position)

	h.store.mu.Lock()
	h.store.nodes[string(node.ID)] = node
	h.store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "鑺傜偣鍒涘缓鎴愬姛",
		"node":    node,
	})
}

func (h *DemoHandlers) GetNode(c *gin.Context) {
	id := c.Param("id")

	h.store.mu.RLock()
	node, exists := h.store.nodes[id]
	h.store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "鑺傜偣涓嶅瓨鍦?})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node})
}

// GetPath 鑾峰彇鍗曚釜璺緞
func (h *DemoHandlers) GetPath(c *gin.Context) {
	id := c.Param("id")

	h.store.mu.RLock()
	path, exists := h.store.paths[id]
	h.store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "璺緞涓嶅瓨鍦?})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "鑺傜偣涓嶅瓨鍦?})
		return
	}

	if req.Name != nil {
		node.Name = *req.Name
	}
	if req.Position != nil {
		node.Position = *req.Position
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "鑺傜偣鏇存柊鎴愬姛",
		"node":    node,
	})
}

func (h *DemoHandlers) DeleteNode(c *gin.Context) {
	id := c.Param("id")

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, exists := h.store.nodes[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "鑺傜偣涓嶅瓨鍦?})
		return
	}

	delete(h.store.nodes, id)

	c.JSON(http.StatusOK, gin.H{"message": "鑺傜偣鍒犻櫎鎴愬姛"})
}

// UpdateNodePosition 鏇存柊鑺傜偣浣嶇疆锛堜粎鍧愭爣锛?
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
		c.JSON(http.StatusNotFound, gin.H{"error": "鑺傜偣涓嶅瓨鍦?})
		return
	}
	node.Position.X = req.X
	node.Position.Y = req.Y
	node.Position.Z = req.Z
	h.store.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "浣嶇疆宸叉洿鏂?, "node": node})
}

// 璺緞鐩稿叧API
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

	// 妫€鏌ヨ妭鐐规槸鍚﹀瓨鍦?
	h.store.mu.RLock()
	_, startExists := h.store.nodes[string(req.StartNodeID)]
	_, endExists := h.store.nodes[string(req.EndNodeID)]
	h.store.mu.RUnlock()

	if !startExists || !endExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "璧峰鑺傜偣鎴栫粨鏉熻妭鐐逛笉瀛樺湪"})
		return
	}

	path := domain.NewPath(req.Name, req.StartNodeID, req.EndNodeID)

	h.store.mu.Lock()
	h.store.paths[string(path.ID)] = path
	h.store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "璺緞鍒涘缓鎴愬姛",
		"path":    path,
	})
}

// DeletePath 鍒犻櫎璺緞
func (h *DemoHandlers) DeletePath(c *gin.Context) {
	id := c.Param("id")
	h.store.mu.Lock()
	defer h.store.mu.Unlock()
	if _, ok := h.store.paths[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "璺緞涓嶅瓨鍦?})
		return
	}
	delete(h.store.paths, id)
	c.JSON(http.StatusOK, gin.H{"message": "璺緞宸插垹闄?})
}

// UpdatePath 鏇存柊璺緞
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "璇锋眰鏍煎紡閿欒: " + err.Error()})
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	path, exists := h.store.paths[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "璺緞涓嶅瓨鍦?})
		return
	}

	// 楠岃瘉璧峰鍜岀粨鏉熻妭鐐瑰瓨鍦?
	if req.StartNodeID != "" {
		if _, exists := h.store.nodes[req.StartNodeID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "璧峰鑺傜偣涓嶅瓨鍦?})
			return
		}
	}

	if req.EndNodeID != "" {
		if _, exists := h.store.nodes[req.EndNodeID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缁撴潫鑺傜偣涓嶅瓨鍦?})
			return
		}
	}

	// 鏇存柊璺緞灞炴€?
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
		"message": "璺緞鏇存柊鎴愬姛",
		"path":    path,
	})
}

// ApplyLayout 搴旂敤甯冨眬绠楁硶
func (h *DemoHandlers) ApplyLayout(c *gin.Context) {
	var req struct {
		Algorithm string `json:"algorithm" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "璇锋眰鏍煎紡閿欒: " + err.Error()})
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	// 鑾峰彇鎵€鏈夎妭鐐瑰拰璺緞
	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := make([]domain.Path, 0, len(h.store.paths))
	for _, path := range h.store.paths {
		paths = append(paths, *path)
	}

	// 搴旂敤甯冨眬绠楁硶
	var updatedNodes []domain.Node
	switch req.Algorithm {
	case "grid":
		updatedNodes = applyGridLayout(nodes, 120.0)
	case "force-directed":
		updatedNodes = applyForceDirectedLayout(nodes, paths, 50)
	case "circular":
		updatedNodes = applyCircularLayout(nodes, 250.0, 500.0, 400.0)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "涓嶆敮鎸佺殑甯冨眬绠楁硶: " + req.Algorithm})
		return
	}

	// 鏇存柊瀛樺偍涓殑鑺傜偣浣嶇疆
	for _, node := range updatedNodes {
		if existingNode, ok := h.store.nodes[string(node.ID)]; ok {
			existingNode.Position = node.Position
			h.store.nodes[string(node.ID)] = existingNode
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "甯冨眬搴旂敤鎴愬姛",
		"algorithm":      req.Algorithm,
		"affected_nodes": len(updatedNodes),
	})
}

// GenerateNearestNeighborPaths 鐢熸垚鏈€杩戦偦璺緞
func (h *DemoHandlers) GenerateNearestNeighborPaths(c *gin.Context) {
	var req struct {
		MaxDistance float64 `json:"max_distance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.MaxDistance = 200.0 // 榛樿鍊?
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := generateNearestNeighborPaths(nodes, req.MaxDistance)

	// 娣诲姞鍒板瓨鍌?
	createdCount := 0
	for _, path := range paths {
		if _, exists := h.store.paths[string(path.ID)]; !exists {
			h.store.paths[string(path.ID)] = &path
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "鏈€杩戦偦璺緞鐢熸垚鎴愬姛",
		"created_paths": createdCount,
		"max_distance":  req.MaxDistance,
	})
}

// GenerateFullConnectivity 鐢熸垚瀹屽叏杩為€氳矾寰?
func (h *DemoHandlers) GenerateFullConnectivity(c *gin.Context) {
	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := generateFullConnectivityPaths(nodes)

	// 娣诲姞鍒板瓨鍌?
	createdCount := 0
	for _, path := range paths {
		if _, exists := h.store.paths[string(path.ID)]; !exists {
			h.store.paths[string(path.ID)] = &path
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "瀹屽叏杩為€氳矾寰勭敓鎴愭垚鍔?,
		"created_paths": createdCount,
		"total_nodes":   len(nodes),
	})
}

// GenerateGridPaths 鐢熸垚缃戞牸璺緞
func (h *DemoHandlers) GenerateGridPaths(c *gin.Context) {
	var req struct {
		ConnectDiagonal bool `json:"connect_diagonal"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		req.ConnectDiagonal = false // 榛樿鍊?
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	nodes := make([]domain.Node, 0, len(h.store.nodes))
	for _, node := range h.store.nodes {
		nodes = append(nodes, *node)
	}

	paths := generateGridPaths(nodes, req.ConnectDiagonal)

	// 娣诲姞鍒板瓨鍌?
	createdCount := 0
	for _, path := range paths {
		if _, exists := h.store.paths[string(path.ID)]; !exists {
			h.store.paths[string(path.ID)] = &path
			createdCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "缃戞牸璺緞鐢熸垚鎴愬姛",
		"created_paths":    createdCount,
		"connect_diagonal": req.ConnectDiagonal,
	})
}

// 甯冨眬绠楁硶瀹炵幇

// applyGridLayout 缃戞牸甯冨眬
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

// applyForceDirectedLayout 鍔涘鍚戝竷灞€ (绠€鍖栫増)
func applyForceDirectedLayout(nodes []domain.Node, paths []domain.Path, iterations int) []domain.Node {
	if len(nodes) == 0 {
		return nodes
	}

	// 鍒濆鍖栭殢鏈虹瀛?
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 鍙傛暟璁剧疆
	width, height := 1000.0, 800.0
	k := math.Sqrt((width * height) / float64(len(nodes)))

	// 鍒濆鍖栬妭鐐逛綅缃?
	updatedNodes := make([]domain.Node, len(nodes))
	for i, node := range nodes {
		updatedNode := node
		if node.Position.X == 0 && node.Position.Y == 0 {
			updatedNode.Position.X = r.Float64() * width
			updatedNode.Position.Y = r.Float64() * height
		}
		updatedNodes[i] = updatedNode
	}

	// 杩唬璁＄畻
	for iter := 0; iter < iterations; iter++ {
		forces := make(map[string]struct{ fx, fy float64 })

		// 鍒濆鍖栧姏
		for i := range updatedNodes {
			forces[string(updatedNodes[i].ID)] = struct{ fx, fy float64 }{0, 0}
		}

		// 璁＄畻鎺掓枼鍔?
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

		// 璁＄畻鍚稿紩鍔?
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

		// 搴旂敤鍔?
		temperature := 10.0 * (1.0 - float64(iter)/float64(iterations))
		for i := range updatedNodes {
			force := forces[string(updatedNodes[i].ID)]
			displacement := math.Min(math.Sqrt(force.fx*force.fx+force.fy*force.fy), temperature)

			if displacement > 0.01 {
				updatedNodes[i].Position.X += force.fx / displacement * temperature
				updatedNodes[i].Position.Y += force.fy / displacement * temperature
			}

			// 淇濇寔鍦ㄧ敾甯冭寖鍥村唴
			updatedNodes[i].Position.X = math.Max(50, math.Min(width-50, updatedNodes[i].Position.X))
			updatedNodes[i].Position.Y = math.Max(50, math.Min(height-50, updatedNodes[i].Position.Y))
		}
	}

	return updatedNodes
}

// applyCircularLayout 鍦嗗舰甯冨眬
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

// 璺緞鐢熸垚绠楁硶瀹炵幇

// generateNearestNeighborPaths 鐢熸垚鏈€杩戦偦璺緞
func generateNearestNeighborPaths(nodes []domain.Node, maxDistance float64) []domain.Path {
	if len(nodes) < 2 {
		return []domain.Path{}
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
		maxNeighbors := minInt(3, len(neighbors))
		for i := 0; i < maxNeighbors; i++ {
			neighbor := neighbors[i]

			// 鍒涘缓鍞竴鐨勮矾寰勬爣璇嗙锛堥槻姝㈤噸澶嶏級
			pathKey := fmt.Sprintf("%s_%s", minString(string(node.ID), string(neighbor.nodeID)), maxString(string(node.ID), string(neighbor.nodeID)))
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
				Metadata:    domain.ObjectMeta{
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

// generateFullConnectivityPaths 鐢熸垚瀹屽叏杩為€氳矾寰?
func generateFullConnectivityPaths(nodes []domain.Node) []domain.Path {
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
				Metadata:    domain.ObjectMeta{
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

// generateGridPaths 鐢熸垚缃戞牸璺緞
func generateGridPaths(nodes []domain.Node, connectDiagonal bool) []domain.Path {
	if len(nodes) == 0 {
		return []domain.Path{}
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
					Metadata:    domain.ObjectMeta{
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
						Metadata:    domain.ObjectMeta{
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
						Metadata:    domain.ObjectMeta{
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

// 宸ュ叿鍑芥暟

// calculateDistance 璁＄畻涓ょ偣涔嬮棿鐨勬鍑犻噷寰楄窛绂?
func calculateDistance(pos1, pos2 domain.Position) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	dz := pos1.Z - pos2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// minInt 杩斿洖涓や釜鏁存暟鐨勬渶灏忓€?
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// minString 杩斿洖涓や釜瀛楃涓茬殑瀛楀吀搴忔渶灏忓€?
func minString(a, b string) string {
	if a < b {
		return a
	}
	return b
}

// maxString 杩斿洖涓や釜瀛楃涓茬殑瀛楀吀搴忔渶澶у€?
func maxString(a, b string) string {
	if a > b {
		return a
	}
	return b
}

// 鍋ュ悍妫€鏌?
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

// 鑾峰彇鐢诲竷鏁版嵁
func (h *DemoHandlers) GetCanvasData(c *gin.Context) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	// 鍑嗗鐢诲竷鏁版嵁
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

// 琛ㄦ牸瑙嗗浘HTML
const tableHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>琛ㄦ牸瑙嗗浘 - 鏈哄櫒浜鸿矾寰勭紪杈戝櫒</title>
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
        <h1>馃搳 琛ㄦ牸瑙嗗浘</h1>
        <div class="nav-links">
            <a href="/">鐢诲竷瑙嗗浘</a>
            <a href="/table">琛ㄦ牸瑙嗗浘</a>
        </div>
    </div>
    
    <div class="controls">
        <div class="view-toggle">
            <button id="nodeViewBtn" class="btn active">鑺傜偣</button>
            <button id="pathViewBtn" class="btn">璺緞</button>
        </div>
        <div class="actions">
            <button id="refreshBtn" class="btn">馃攧 鍒锋柊</button>
            <button id="addBtn" class="btn btn-primary">鉃?娣诲姞</button>
        </div>
    </div>
    
    <div class="main-content">
        <div id="tableContainer">
            <!-- 琛ㄦ牸鍐呭灏嗗湪姝ゅ鍔ㄦ€佺敓鎴?-->
        </div>
    </div>
    
    <script src="/static/table.js"></script>
</body>
</html>`

// 涓婚〉闈TML
const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>鏈哄櫒浜鸿矾寰勭紪杈戝櫒 - 婕旂ず鐗?/title>
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
        <button id="undoBtn" style="padding:8px 12px;background:#3498db;color:white;border:none;border-radius:4px;cursor:pointer;" disabled title="鎾ら攢 (Ctrl+Z)">鈫?鎾ら攢</button>
        <button id="redoBtn" style="padding:8px 12px;background:#3498db;color:white;border:none;border-radius:4px;cursor:pointer;" disabled title="閲嶅仛 (Ctrl+Y)">鈫?閲嶅仛</button>
        <div style="border-left:1px solid rgba(255,255,255,0.3);margin:0 10px;"></div>
        <button id="gridLayoutBtn" style="padding:8px 12px;background:#27ae60;color:white;border:none;border-radius:4px;cursor:pointer;" title="缃戞牸甯冨眬">馃敵 缃戞牸</button>
        <button id="forceLayoutBtn" style="padding:8px 12px;background:#e67e22;color:white;border:none;border-radius:4px;cursor:pointer;" title="鍔涘鍚戝竷灞€">鈿?鍔涘鍚?/button>
        <button id="circularLayoutBtn" style="padding:8px 12px;background:#9b59b6;color:white;border:none;border-radius:4px;cursor:pointer;" title="鍦嗗舰甯冨眬">猸?鍦嗗舰</button>
        <div style="border-left:1px solid rgba(255,255,255,0.3);margin:0 10px;"></div>
        <button id="nearestPathBtn" style="padding:8px 12px;background:#f39c12;color:white;border:none;border-radius:4px;cursor:pointer;" title="鐢熸垚鏈€杩戦偦璺緞">馃敆 鏈€杩戦偦</button>
        <button id="fullConnectBtn" style="padding:8px 12px;background:#e74c3c;color:white;border:none;border-radius:4px;cursor:pointer;" title="鐢熸垚瀹屽叏杩為€?>馃暩锔?鍏ㄨ繛閫?/button>
        <button id="gridPathBtn" style="padding:8px 12px;background:#8e44ad;color:white;border:none;border-radius:4px;cursor:pointer;" title="鐢熸垚缃戞牸璺緞">馃搻 缃戞牸璺緞</button>
    </div>
    <div id="canvas-container" style="position:fixed;left:0;top:0;width:calc(100% - 300px);height:100vh;"></div>
    <div class="container">
        <div class="logo">馃</div>
        <h1>鏈哄櫒浜鸿矾寰勭紪杈戝櫒</h1>
        <p class="subtitle">鐜颁唬鍖栫殑涓夌鍏煎璺緞绠＄悊宸ュ叿</p>
        
        <div class="demo-badge">馃殌 婕旂ず鐗堟湰 - 鍐呭瓨瀛樺偍妯″紡</div>
        
        <div class="stats">
            <h3>馃搳 瀹炴椂鏁版嵁缁熻</h3>
            <div class="stats-grid">
                <div class="stat-item">
                    <div class="stat-number" id="nodeCount">-</div>
                    <div class="stat-label">鑺傜偣鏁伴噺</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number" id="pathCount">-</div>
                    <div class="stat-label">璺緞鏁伴噺</div>
                </div>
            </div>
        </div>
        
        <div class="features">
            <div class="feature">
                <h3>馃梽锔?鍐呭瓨瀛樺偍</h3>
                <p>鏃犻渶鏁版嵁搴擄紝蹇€熷惎鍔ㄦ紨绀?/p>
            </div>
            <div class="feature">
                <h3>馃帹 RESTful API</h3>
                <p>瀹屾暣鐨勮妭鐐瑰拰璺緞绠＄悊鎺ュ彛</p>
            </div>
            <div class="feature">
                <h3>馃摫 鍝嶅簲寮忚璁?/h3>
                <p>鑷€傚簲涓嶅悓璁惧灞忓箷</p>
            </div>
            <div class="feature">
                <h3>鈿?瀹炴椂鏇存柊</h3>
                <p>鏁版嵁鍙樺寲瀹炴椂鍚屾鏄剧ず</p>
            </div>
        </div>
        
        <div class="api-section">
            <h3>馃敆 蹇€熷鑸?/h3>
            <div class="nav-buttons" style="display:flex;gap:1rem;margin:1rem 0;">
                <a href="/" style="text-decoration:none;">
                    <button class="btn" style="width:100%;">馃帹 鐢诲竷瑙嗗浘</button>
                </a>
                <a href="/table" style="text-decoration:none;">
                    <button class="btn" style="width:100%;">馃搳 琛ㄦ牸瑙嗗浘</button>
                </a>
            </div>
            
            <h3>馃攲 API 绔偣</h3>
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
        // 瀹炴椂鏇存柊缁熻鏁版嵁
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
                    console.error('鑾峰彇缁熻鏁版嵁澶辫触:', error);
                });
        }
        
        // 椤甸潰鍔犺浇鏃舵洿鏂颁竴娆?
        updateStats();
        
        // 姣?绉掓洿鏂颁竴娆?
        setInterval(updateStats, 5000);
        
        // 绠€鍗曠殑API娴嬭瘯
        console.log('馃 鏈哄櫒浜鸿矾寰勭紪杈戝櫒婕旂ず鐗堝凡鍚姩');
        console.log('馃摗 API娴嬭瘯:');
        
        fetch('/api/v1/nodes')
            .then(response => response.json())
            .then(data => {
                console.log('鉁?鑺傜偣API娴嬭瘯鎴愬姛:', data);
            })
            .catch(error => {
                console.error('鉂?鑺傜偣API娴嬭瘯澶辫触:', error);
            });
    </script>
</body>
</html>`

func main() {
	// 璁剧疆鏃ュ織
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	fmt.Println("馃 鏈哄櫒浜鸿矾寰勭紪杈戝櫒婕旂ず鐗堝惎鍔ㄤ腑...")

	// 鍒濆鍖栧唴瀛樺瓨鍌?
	store := NewMemoryStore()
	handlers := NewDemoHandlers(store)

	// 璁剧疆Gin璺敱
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

	// CORS涓棿浠?
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

	// 涓婚〉闈?
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})

	// 琛ㄦ牸瑙嗗浘椤甸潰
	r.GET("/table", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(tableHTML))
	})

	// 鍋ュ悍妫€鏌?
	r.GET("/health", handlers.HealthCheck)

	// 鐢诲竷鏁版嵁
	r.GET("/canvas-data", handlers.GetCanvasData)

	// API璺敱
	api := r.Group("/api/v1")
	{
		// 鑺傜偣绠＄悊
		nodes := api.Group("/nodes")
		{
			nodes.GET("", handlers.ListNodes)
			nodes.POST("", handlers.CreateNode)
			nodes.GET("/:id", handlers.GetNode)
			nodes.PUT("/:id", handlers.UpdateNode)
			nodes.PUT("/:id/position", handlers.UpdateNodePosition)
			nodes.DELETE("/:id", handlers.DeleteNode)
		}

		// 璺緞绠＄悊
		paths := api.Group("/paths")
		{
			paths.GET("", handlers.ListPaths)
			paths.GET("/:id", handlers.GetPath)
			paths.POST("", handlers.CreatePath)
			paths.PUT("/:id", handlers.UpdatePath)
			paths.DELETE("/:id", handlers.DeletePath)
		}

		// 甯冨眬绠楁硶绔偣
		layout := api.Group("/layout")
		{
			layout.POST("/apply", handlers.ApplyLayout)
		}

		// 璺緞鐢熸垚绔偣
		pathGen := api.Group("/path-generation")
		{
			pathGen.POST("/nearest-neighbor", handlers.GenerateNearestNeighborPaths)
			pathGen.POST("/full-connectivity", handlers.GenerateFullConnectivity)
			pathGen.POST("/grid", handlers.GenerateGridPaths)
		}
	}

	// 鍚姩鏈嶅姟鍣?
	port := ":8080"
	fmt.Printf("馃殌 婕旂ず鏈嶅姟鍣ㄥ惎鍔ㄦ垚鍔? 璁块棶鍦板潃: http://localhost%s\n", port)
	fmt.Println("馃搳 API绔偣:")
	fmt.Println("  - GET  /health        鍋ュ悍妫€鏌?)
	fmt.Println("  - GET  /canvas-data   鐢诲竷鏁版嵁")
	fmt.Println("  - GET  /api/v1/nodes  鑺傜偣鍒楄〃")
	fmt.Println("  - POST /api/v1/nodes  鍒涘缓鑺傜偣")
	fmt.Println("  - GET  /api/v1/paths  璺緞鍒楄〃")
	fmt.Println("  - POST /api/v1/paths  鍒涘缓璺緞")

	if err := r.Run(port); err != nil {
		logrus.WithError(err).Fatal("鏈嶅姟鍣ㄥ惎鍔ㄥけ璐?)
	}
}
