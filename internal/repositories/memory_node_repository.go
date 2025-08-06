// Package repositories 内存节点仓储实现
// 用于演示，不依赖外部数据库
package repositories

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"robot-path-editor/internal/domain"
)

// memoryNodeRepository 内存节点仓储实现
type memoryNodeRepository struct {
	nodes map[string]*domain.Node
	mu    sync.RWMutex
}

// NewMemoryNodeRepository 创建内存节点仓储实例
func NewMemoryNodeRepository() NodeRepository {
	return &memoryNodeRepository{
		nodes: make(map[string]*domain.Node),
	}
}

// Create 创建节点
func (r *memoryNodeRepository) Create(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("节点验证失败: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[string(node.ID)]; exists {
		return fmt.Errorf("节点已存�? %s", node.ID)
	}

	r.nodes[string(node.ID)] = node
	return nil
}

// GetByID 根据ID获取节点
func (r *memoryNodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	node, exists := r.nodes[string(id)]
	if !exists {
		return nil, fmt.Errorf("节点不存�? %s", id)
	}

	// 返回副本以避免并发修�?
	nodeCopy := *node
	return &nodeCopy, nil
}

// Update 更新节点
func (r *memoryNodeRepository) Update(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("节点验证失败: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[string(node.ID)]; !exists {
		return fmt.Errorf("节点不存�? %s", node.ID)
	}

	r.nodes[string(node.ID)] = node
	return nil
}

// Delete 删除节点
func (r *memoryNodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[string(id)]; !exists {
		return fmt.Errorf("节点不存�? %s", id)
	}

	delete(r.nodes, string(id))
	return nil
}

// CreateBatch 批量创建节点
func (r *memoryNodeRepository) CreateBatch(ctx context.Context, nodes []*domain.Node) error {
	for _, node := range nodes {
		if err := r.Create(ctx, node); err != nil {
			return err
		}
	}
	return nil
}

// GetByIDs 根据ID列表获取节点
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

// UpdateBatch 批量更新节点
func (r *memoryNodeRepository) UpdateBatch(ctx context.Context, nodes []*domain.Node) error {
	for _, node := range nodes {
		if err := r.Update(ctx, node); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch 批量删除节点
func (r *memoryNodeRepository) DeleteBatch(ctx context.Context, ids []domain.NodeID) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// List 列表查询节点
func (r *memoryNodeRepository) List(ctx context.Context, options ListOptions) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 首先应用过滤�?
	var filtered []*domain.Node
	for _, node := range r.nodes {
		if r.matchesFilter(node, options.Filter) {
			nodeCopy := *node
			filtered = append(filtered, &nodeCopy)
		}
	}

	// 应用分页
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

// Count 统计节点数量
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

// GetByArea 获取指定区域内的节点
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

// GetNearby 获取指定位置附近的节�?
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

// GetConnectedNodes 获取与指定节点连接的所有节�?
func (r *memoryNodeRepository) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	// 内存实现中，这需要路径信息，暂时返回空列�?
	return []*domain.Node{}, nil
}

// GetIsolatedNodes 获取孤立节点
func (r *memoryNodeRepository) GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error) {
	// 内存实现中，这需要路径信息，暂时返回所有节�?
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, node := range r.nodes {
		nodeCopy := *node
		nodes = append(nodes, &nodeCopy)
	}

	return nodes, nil
}

// GetByLabels 根据标签查询节点
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

// GetByType 根据类型查询节点
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

// GetByStatus 根据状态查询节�?
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

// matchesFilter 检查节点是否匹配过滤器
func (r *memoryNodeRepository) matchesFilter(node *domain.Node, filter NodeFilter) bool {
	// ID过滤
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

	// 名称过滤（模糊查询）
	if filter.Name != "" && !strings.Contains(strings.ToLower(node.Name), strings.ToLower(filter.Name)) {
		return false
	}

	// 类型过滤
	if filter.Type != "" && node.Type != filter.Type {
		return false
	}

	// 状态过�?
	if filter.Status != "" && node.Status != filter.Status {
		return false
	}

	// 位置过滤
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

	// 标签过滤
	for key, value := range filter.Labels {
		if nodeValue, exists := node.Metadata.Labels[key]; !exists || nodeValue != value {
			return false
		}
	}

	return true
}
