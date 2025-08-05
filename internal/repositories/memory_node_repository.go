// Package repositories 鍐呭瓨鑺傜偣浠撳偍瀹炵幇
// 鐢ㄤ簬婕旂ず锛屼笉渚濊禆澶栭儴鏁版嵁搴?
package repositories

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"robot-path-editor/internal/domain"
)

// memoryNodeRepository 鍐呭瓨鑺傜偣浠撳偍瀹炵幇
type memoryNodeRepository struct {
	nodes map[string]*domain.Node
	mu    sync.RWMutex
}

// NewMemoryNodeRepository 鍒涘缓鍐呭瓨鑺傜偣浠撳偍瀹炰緥
func NewMemoryNodeRepository() NodeRepository {
	return &memoryNodeRepository{
		nodes: make(map[string]*domain.Node),
	}
}

// Create 鍒涘缓鑺傜偣
func (r *memoryNodeRepository) Create(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[string(node.ID)]; exists {
		return fmt.Errorf("鑺傜偣宸插瓨鍦? %s", node.ID)
	}

	r.nodes[string(node.ID)] = node
	return nil
}

// GetByID 鏍规嵁ID鑾峰彇鑺傜偣
func (r *memoryNodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	node, exists := r.nodes[string(id)]
	if !exists {
		return nil, fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", id)
	}

	// 杩斿洖鍓湰浠ラ伩鍏嶅苟鍙戜慨鏀?
	nodeCopy := *node
	return &nodeCopy, nil
}

// Update 鏇存柊鑺傜偣
func (r *memoryNodeRepository) Update(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[string(node.ID)]; !exists {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", node.ID)
	}

	r.nodes[string(node.ID)] = node
	return nil
}

// Delete 鍒犻櫎鑺傜偣
func (r *memoryNodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[string(id)]; !exists {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", id)
	}

	delete(r.nodes, string(id))
	return nil
}

// CreateBatch 鎵归噺鍒涘缓鑺傜偣
func (r *memoryNodeRepository) CreateBatch(ctx context.Context, nodes []*domain.Node) error {
	for _, node := range nodes {
		if err := r.Create(ctx, node); err != nil {
			return err
		}
	}
	return nil
}

// GetByIDs 鏍规嵁ID鍒楄〃鑾峰彇鑺傜偣
func (r *memoryNodeRepository) GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, id := range ids {
		if node, exists := r.nodes[string(id)]; exists {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// UpdateBatch 鎵归噺鏇存柊鑺傜偣
func (r *memoryNodeRepository) UpdateBatch(ctx context.Context, nodes []*domain.Node) error {
	for _, node := range nodes {
		if err := r.Update(ctx, node); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch 鎵归噺鍒犻櫎鑺傜偣
func (r *memoryNodeRepository) DeleteBatch(ctx context.Context, ids []domain.NodeID) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// List 鍒楄〃鏌ヨ鑺傜偣
func (r *memoryNodeRepository) List(ctx context.Context, options ListOptions) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 棣栧厛搴旂敤杩囨护鍣?
	var filtered []*domain.Node
	for _, node := range r.nodes {
		if r.matchesFilter(node, options.Filter) {
			nodeCopy := *node
			filtered = append(filtered, &nodeCopy)
		}
	}

	// 搴旂敤鍒嗛〉
	if options.PageSize > 0 {
		start := 0
		if options.Page > 1 {
			start = (options.Page - 1) * options.PageSize
		}

		end := start + options.PageSize
		if start >= len(filtered) {
			return []*domain.Node{}, nil
		}
		if end > len(filtered) {
			end = len(filtered)
		}

		filtered = filtered[start:end]
	}

	return filtered, nil
}

// Count 缁熻鑺傜偣鏁伴噺
func (r *memoryNodeRepository) Count(ctx context.Context, filter NodeFilter) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := int64(0)
	for _, node := range r.nodes {
		if r.matchesFilter(node, filter) {
			count++
		}
	}

	return count, nil
}

// GetByArea 鑾峰彇鎸囧畾鍖哄煙鍐呯殑鑺傜偣
func (r *memoryNodeRepository) GetByArea(ctx context.Context, minX, minY, maxX, maxY float64) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		if node.Position.X >= minX && node.Position.X <= maxX &&
			node.Position.Y >= minY && node.Position.Y <= maxY {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// GetNearby 鑾峰彇鎸囧畾浣嶇疆闄勮繎鐨勮妭鐐?
func (r *memoryNodeRepository) GetNearby(ctx context.Context, position domain.Position, radius float64) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		distance := position.Distance(node.Position)
		if distance <= radius {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// GetConnectedNodes 鑾峰彇涓庢寚瀹氳妭鐐硅繛鎺ョ殑鎵€鏈夎妭鐐?
func (r *memoryNodeRepository) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	// 鍐呭瓨瀹炵幇涓紝杩欓渶瑕佽矾寰勪俊鎭紝鏆傛椂杩斿洖绌哄垪琛?
	return []*domain.Node{}, nil
}

// GetIsolatedNodes 鑾峰彇瀛ょ珛鑺傜偣
func (r *memoryNodeRepository) GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error) {
	// 鍐呭瓨瀹炵幇涓紝杩欓渶瑕佽矾寰勪俊鎭紝鏆傛椂杩斿洖鎵€鏈夎妭鐐?
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		nodeCopy := *node
		nodes = append(nodes, &nodeCopy)
	}

	return nodes, nil
}

// GetByLabels 鏍规嵁鏍囩鏌ヨ鑺傜偣
func (r *memoryNodeRepository) GetByLabels(ctx context.Context, labels map[string]string) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		match := true
		for key, value := range labels {
			if nodeValue, exists := node.Metadata.Labels[key]; !exists || nodeValue != value {
				match = false
				break
			}
		}

		if match {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// GetByType 鏍规嵁绫诲瀷鏌ヨ鑺傜偣
func (r *memoryNodeRepository) GetByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		if node.Type == nodeType {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// GetByStatus 鏍规嵁鐘舵€佹煡璇㈣妭鐐?
func (r *memoryNodeRepository) GetByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		if node.Status == status {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// matchesFilter 妫€鏌ヨ妭鐐规槸鍚﹀尮閰嶈繃婊ゅ櫒
func (r *memoryNodeRepository) matchesFilter(node *domain.Node, filter NodeFilter) bool {
	// ID杩囨护
	if len(filter.IDs) > 0 {
		found := false
		for _, id := range filter.IDs {
			if node.ID == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 鍚嶇О杩囨护锛堟ā绯婃煡璇級
	if filter.Name != "" && !strings.Contains(strings.ToLower(node.Name), strings.ToLower(filter.Name)) {
		return false
	}

	// 绫诲瀷杩囨护
	if filter.Type != "" && node.Type != filter.Type {
		return false
	}

	// 鐘舵€佽繃婊?
	if filter.Status != "" && node.Status != filter.Status {
		return false
	}

	// 浣嶇疆杩囨护
	if filter.MinX != nil && node.Position.X < *filter.MinX {
		return false
	}
	if filter.MaxX != nil && node.Position.X > *filter.MaxX {
		return false
	}
	if filter.MinY != nil && node.Position.Y < *filter.MinY {
		return false
	}
	if filter.MaxY != nil && node.Position.Y > *filter.MaxY {
		return false
	}
	if filter.MinZ != nil && node.Position.Z < *filter.MinZ {
		return false
	}
	if filter.MaxZ != nil && node.Position.Z > *filter.MaxZ {
		return false
	}

	// 鏍囩杩囨护
	for key, value := range filter.Labels {
		if nodeValue, exists := node.Metadata.Labels[key]; !exists || nodeValue != value {
			return false
		}
	}

	return true
}
