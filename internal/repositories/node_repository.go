// Package repositories 实现数据访问层
//
// 设计参考：
// - DDD的仓储模式
// - Kubernetes的存储抽象
// - GitHub的仓储实现模式
//
// 特点：
// 1. 接口抽象：定义清晰的数据访问接口
// 2. 实现分离：支持不同的存储后端
// 3. 查询优化：支持复杂查询和分页
// 4. 缓存友好：设计便于缓存集成
package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// NodeRepository 节点仓储接口
// 定义节点数据访问的所有操作
type NodeRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, node *domain.Node) error
	GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error)
	Update(ctx context.Context, node *domain.Node) error
	Delete(ctx context.Context, id domain.NodeID) error

	// 批量操作
	GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error)
	UpdateBatch(ctx context.Context, nodes []*domain.Node) error
	DeleteBatch(ctx context.Context, ids []domain.NodeID) error

	// 查询操作
	List(ctx context.Context, filter NodeFilter) ([]*domain.Node, error)
	Count(ctx context.Context, filter NodeFilter) (int64, error)
	Search(ctx context.Context, query string, filter NodeFilter) ([]*domain.Node, error)

	// 关系查询
	GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error)
	GetNodesByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error)
	GetNodesByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error)
}

// NodeFilter 节点查询过滤器
type NodeFilter struct {
	IDs    []domain.NodeID   `json:"ids,omitempty"`
	Name   string            `json:"name,omitempty"`
	Type   domain.NodeType   `json:"type,omitempty"`
	Status domain.NodeStatus `json:"status,omitempty"`

	// 分页参数
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`

	// 排序参数
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// nodeRepository GORM实现
type nodeRepository struct {
	db database.Database
}

// NewNodeRepository 创建新的节点仓储实例
func NewNodeRepository(db database.Database) NodeRepository {
	return &nodeRepository{db: db}
}

// Create 创建节点
func (r *nodeRepository) Create(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("节点验证失败: %w", err)
	}

	// 检查是否支持内存数据库的直接接口
	if memDB, ok := r.db.(interface {
		CreateNode(*domain.Node) error
	}); ok {
		return memDB.CreateNode(node)
	}

	return fmt.Errorf("不支持的数据库类型")
}

// GetByID 根据ID获取节点
func (r *nodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	var node domain.Node
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&node).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("节点不存在: %s", id)
		}
		return nil, err
	}
	return &node, nil
}

// Update 更新节点
func (r *nodeRepository) Update(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("节点验证失败: %w", err)
	}

	result := r.db.GORMDB().WithContext(ctx).Save(node)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("节点不存在: %s", node.ID)
	}

	return nil
}

// Delete 删除节点
func (r *nodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	result := r.db.GORMDB().WithContext(ctx).Delete(&domain.Node{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("节点不存在: %s", id)
	}

	return nil
}

// GetByIDs 根据ID列表获取节点
func (r *nodeRepository) GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error) {
	if len(ids) == 0 {
		return []*domain.Node{}, nil
	}

	var nodes []*domain.Node

	// 将NodeID转换为字符串，因为GORM需要基础类型
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.GORMDB().WithContext(ctx).Where("id IN ?", stringIDs).Find(&nodes).Error
	return nodes, err
}

// UpdateBatch 批量更新节点
func (r *nodeRepository) UpdateBatch(ctx context.Context, nodes []*domain.Node) error {
	return r.db.Transaction(ctx, func(tx interface{}) error {
		gormTx := tx.(*gorm.DB)
		for _, node := range nodes {
			if err := node.IsValid(); err != nil {
				return fmt.Errorf("节点验证失败: %w", err)
			}
			if err := gormTx.Save(node).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteBatch 批量删除节点
func (r *nodeRepository) DeleteBatch(ctx context.Context, ids []domain.NodeID) error {
	if len(ids) == 0 {
		return nil
	}

	return r.db.Transaction(ctx, func(tx interface{}) error {
		gormTx := tx.(*gorm.DB)
		for _, id := range ids {
			if err := gormTx.Delete(&domain.Node{}, "id = ?", id).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// List 列出节点
func (r *nodeRepository) List(ctx context.Context, filter NodeFilter) ([]*domain.Node, error) {
	var nodes []*domain.Node

	query := r.db.GORMDB().WithContext(ctx)

	// 应用过滤条件
	if len(filter.IDs) > 0 {
		stringIDs := make([]string, len(filter.IDs))
		for i, id := range filter.IDs {
			stringIDs[i] = string(id)
		}
		query = query.Where("id IN ?", stringIDs)
	}

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 应用排序
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	}

	// 应用分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.Page > 0 {
			offset = (filter.Page - 1) * filter.PageSize
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Find(&nodes).Error
	return nodes, err
}

// Count 统计节点数量
func (r *nodeRepository) Count(ctx context.Context, filter NodeFilter) (int64, error) {
	var count int64

	query := r.db.GORMDB().WithContext(ctx).Model(&domain.Node{})

	// 应用过滤条件（复用List方法的逻辑）
	if len(filter.IDs) > 0 {
		stringIDs := make([]string, len(filter.IDs))
		for i, id := range filter.IDs {
			stringIDs[i] = string(id)
		}
		query = query.Where("id IN ?", stringIDs)
	}

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Count(&count).Error
	return count, err
}

// Search 搜索节点
func (r *nodeRepository) Search(ctx context.Context, query string, filter NodeFilter) ([]*domain.Node, error) {
	// 在名称和属性中搜索
	searchFilter := filter
	searchFilter.Name = query // 简化实现，实际项目中可能需要更复杂的全文搜索

	return r.List(ctx, searchFilter)
}

// GetConnectedNodes 获取连接的节点
func (r *nodeRepository) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 通过路径表查找连接的节点
	// 这里需要联合查询paths表
	err := r.db.GORMDB().WithContext(ctx).
		Table("nodes").
		Joins("JOIN paths ON (nodes.id = paths.start_node_id OR nodes.id = paths.end_node_id)").
		Where("(paths.start_node_id = ? OR paths.end_node_id = ?) AND nodes.id != ?",
			nodeID, nodeID, nodeID).
		Find(&nodes).Error

	return nodes, err
}

// GetNodesByType 根据类型获取节点
func (r *nodeRepository) GetNodesByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error) {
	var nodes []*domain.Node
	err := r.db.GORMDB().WithContext(ctx).Where("type = ?", nodeType).Find(&nodes).Error
	return nodes, err
}

// GetNodesByStatus 根据状态获取节点
func (r *nodeRepository) GetNodesByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error) {
	var nodes []*domain.Node
	err := r.db.GORMDB().WithContext(ctx).Where("status = ?", status).Find(&nodes).Error
	return nodes, err
}
