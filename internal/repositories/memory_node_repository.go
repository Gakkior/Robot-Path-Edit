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
	nodes map[domain.NodeID]*domain.Node
	mu    sync.RWMutex
}

// NewMemoryNodeRepository 创建内存节点仓储实例
func NewMemoryNodeRepository() NodeRepository {
	return &memoryNodeRepository{
		nodes: make(map[domain.NodeID]*domain.Node),
	}
}

// Create 创建节点
func (r *memoryNodeRepository) Create(ctx context.Context, node *domain.Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[node.ID]; exists {
		return fmt.Errorf("节点已存在: %s", node.ID)
	}

	// 创建副本以避免外部修改
	nodeCopy := *node
	r.nodes[node.ID] = &nodeCopy
	return nil
}

// GetByID 根据ID获取节点
func (r *memoryNodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	node, exists := r.nodes[id]
	if !exists {
		return nil, fmt.Errorf("节点不存在: %s", id)
	}

	// 返回副本以避免并发修改
	nodeCopy := *node
	return &nodeCopy, nil
}

// Update 更新节点
func (r *memoryNodeRepository) Update(ctx context.Context, node *domain.Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[node.ID]; !exists {
		return fmt.Errorf("节点不存在: %s", node.ID)
	}

	// 创建副本
	nodeCopy := *node
	r.nodes[node.ID] = &nodeCopy
	return nil
}

// Delete 删除节点
func (r *memoryNodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[id]; !exists {
		return fmt.Errorf("节点不存在: %s", id)
	}

	delete(r.nodes, id)
	return nil
}

// GetByIDs 根据ID列表获取节点
func (r *memoryNodeRepository) GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var nodes []*domain.Node
	for _, id := range ids {
		if node, exists := r.nodes[id]; exists {
			nodeCopy := *node
			nodes = append(nodes, &nodeCopy)
		}
	}

	return nodes, nil
}

// UpdateBatch 批量更新节点
func (r *memoryNodeRepository) UpdateBatch(ctx context.Context, nodes []*domain.Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, node := range nodes {
		if _, exists := r.nodes[node.ID]; !exists {
			return fmt.Errorf("节点不存在: %s", node.ID)
		}
		nodeCopy := *node
		r.nodes[node.ID] = &nodeCopy
	}

	return nil
}

// DeleteBatch 批量删除节点
func (r *memoryNodeRepository) DeleteBatch(ctx context.Context, ids []domain.NodeID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		if _, exists := r.nodes[id]; !exists {
			return fmt.Errorf("节点不存在: %s", id)
		}
		delete(r.nodes, id)
	}

	return nil
}

// List 列出节点
func (r *memoryNodeRepository) List(ctx context.Context, filter NodeFilter) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allNodes []*domain.Node
	for _, node := range r.nodes {
		nodeCopy := *node
		allNodes = append(allNodes, &nodeCopy)
	}

	// 首先应用过滤器
	var filteredNodes []*domain.Node
	for _, node := range allNodes {
		if r.matchesFilter(node, filter) {
			filteredNodes = append(filteredNodes, node)
		}
	}

	// 应用分页
	if filter.PageSize > 0 {
		start := 0
		if filter.Page > 0 {
			start = (filter.Page - 1) * filter.PageSize
		}

		end := start + filter.PageSize
		if start >= len(filteredNodes) {
			return []*domain.Node{}, nil
		}
		if end > len(filteredNodes) {
			end = len(filteredNodes)
		}

		filteredNodes = filteredNodes[start:end]
	}

	return filteredNodes, nil
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

// Search 搜索节点
func (r *memoryNodeRepository) Search(ctx context.Context, query string, filter NodeFilter) ([]*domain.Node, error) {
	searchFilter := filter
	searchFilter.Name = query
	return r.List(ctx, searchFilter)
}

// GetConnectedNodes 获取与指定节点连接的所有节点
func (r *memoryNodeRepository) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	// 内存实现中，这需要路径信息，暂时返回空列表
	return []*domain.Node{}, nil
}

// GetNodesByType 根据类型获取节点
func (r *memoryNodeRepository) GetNodesByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error) {
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

// GetNodesByStatus 根据状态获取节点
func (r *memoryNodeRepository) GetNodesByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error) {
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

// matchesFilter 检查节点是否匹配过滤条件
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

	// 名称过滤
	if filter.Name != "" {
		if !strings.Contains(strings.ToLower(node.Name), strings.ToLower(filter.Name)) {
			return false
		}
	}

	// 类型过滤
	if filter.Type != "" && node.Type != filter.Type {
		return false
	}

	// 状态过滤
	if filter.Status != "" && node.Status != filter.Status {
		return false
	}

	return true
}
