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
// 定义节点数据访问的所有操�?
type NodeRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, node *domain.Node) error
	GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error)
	Update(ctx context.Context, node *domain.Node) error
	Delete(ctx context.Context, id domain.NodeID) error

	// 批量操作
	CreateBatch(ctx context.Context, nodes []*domain.Node) error
	GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error)
	UpdateBatch(ctx context.Context, nodes []*domain.Node) error
	DeleteBatch(ctx context.Context, ids []domain.NodeID) error

	// 查询操作
	List(ctx context.Context, options ListOptions) ([]*domain.Node, error)
	Count(ctx context.Context, filter NodeFilter) (int64, error)

	// 空间查询 - 用于画布操作
	GetByArea(ctx context.Context, minX, minY, maxX, maxY float64) ([]*domain.Node, error)
	GetNearby(ctx context.Context, position domain.Position, radius float64) ([]*domain.Node, error)

	// 关系查询
	GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error)
	GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error)

	// 元数据查�?
	GetByLabels(ctx context.Context, labels map[string]string) ([]*domain.Node, error)
	GetByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error)
	GetByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error)
}

// NodeFilter 节点查询过滤�?
type NodeFilter struct {
	IDs    []domain.NodeID   `json:"ids,omitempty"`
	Name   string            `json:"name,omitempty"`
	Type   domain.NodeType   `json:"type,omitempty"`
	Status domain.NodeStatus `json:"status,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`

	// 位置过滤
	MinX *float64 `json:"min_x,omitempty"`
	MaxX *float64 `json:"max_x,omitempty"`
	MinY *float64 `json:"min_y,omitempty"`
	MaxY *float64 `json:"max_y,omitempty"`
	MinZ *float64 `json:"min_z,omitempty"`
	MaxZ *float64 `json:"max_z,omitempty"`
}

// ListOptions 列表查询选项
type ListOptions struct {
	Filter   NodeFilter `json:"filter"`
	Page     int        `json:"page"`      // 页码，从1开�?
	PageSize int        `json:"page_size"` // 页大�?
	OrderBy  string     `json:"order_by"`  // 排序字段
	Order    string     `json:"order"`     // 排序方向: asc, desc
}

// nodeRepository 节点仓储实现
type nodeRepository struct {
	db database.Database
}

// NewNodeRepository 创建节点仓储实例
func NewNodeRepository(db database.Database) NodeRepository {
	return &nodeRepository{
		db: db,
	}
}

// Create 创建节点
func (r *nodeRepository) Create(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("节点验证失败: %w", err)
	}

	// 检查是否为GORM数据�?
	if gormDB := r.db.GORMDB(); gormDB != nil {
		return gormDB.WithContext(ctx).Create(node).Error
	}

	// 如果是内存数据库，需要类型断言
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
		return fmt.Errorf("节点不存�? %s", node.ID)
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
		return fmt.Errorf("节点不存�? %s", id)
	}

	return nil
}

// CreateBatch 批量创建节点
func (r *nodeRepository) CreateBatch(ctx context.Context, nodes []*domain.Node) error {
	// 验证所有节�?
	for _, node := range nodes {
		if err := node.IsValid(); err != nil {
			return fmt.Errorf("节点验证失败: %w", err)
		}
	}

	// 批量插入 - 使用事务确保一致�?
	return r.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).CreateInBatches(nodes, 100).Error
	})
}

// GetByIDs 根据ID列表获取节点
func (r *nodeRepository) GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 转换为字符串切片
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.GORMDB().WithContext(ctx).Where("id IN ?", stringIDs).Find(&nodes).Error
	return nodes, err
}

// UpdateBatch 批量更新节点
func (r *nodeRepository) UpdateBatch(ctx context.Context, nodes []*domain.Node) error {
	return r.db.Transaction(ctx, func(tx *gorm.DB) error {
		for _, node := range nodes {
			if err := node.IsValid(); err != nil {
				return fmt.Errorf("节点验证失败: %w", err)
			}

			if err := tx.WithContext(ctx).Save(node).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteBatch 批量删除节点
func (r *nodeRepository) DeleteBatch(ctx context.Context, ids []domain.NodeID) error {
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	return r.db.GORMDB().WithContext(ctx).Delete(&domain.Node{}, "id IN ?", stringIDs).Error
}

// List 列表查询节点
func (r *nodeRepository) List(ctx context.Context, options ListOptions) ([]*domain.Node, error) {
	var nodes []*domain.Node

	query := r.db.GORMDB().WithContext(ctx)

	// 应用过滤�?
	query = r.applyFilter(query, options.Filter)

	// 应用排序
	if options.OrderBy != "" {
		order := "asc"
		if options.Order == "desc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", options.OrderBy, order))
	} else {
		query = query.Order("created_at desc") // 默认按创建时间降�?
	}

	// 应用分页
	if options.PageSize > 0 {
		offset := 0
		if options.Page > 1 {
			offset = (options.Page - 1) * options.PageSize
		}
		query = query.Offset(offset).Limit(options.PageSize)
	}

	err := query.Find(&nodes).Error
	return nodes, err
}

// Count 统计节点数量
func (r *nodeRepository) Count(ctx context.Context, filter NodeFilter) (int64, error) {
	var count int64

	query := r.db.GORMDB().WithContext(ctx).Model(&domain.Node{})
	query = r.applyFilter(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// GetByArea 获取指定区域内的节点
func (r *nodeRepository) GetByArea(ctx context.Context, minX, minY, maxX, maxY float64) ([]*domain.Node, error) {
	var nodes []*domain.Node

	err := r.db.GORMDB().WithContext(ctx).
		Where("pos_x BETWEEN ? AND ?", minX, maxX).
		Where("pos_y BETWEEN ? AND ?", minY, maxY).
		Find(&nodes).Error

	return nodes, err
}

// GetNearby 获取指定位置附近的节�?
func (r *nodeRepository) GetNearby(ctx context.Context, position domain.Position, radius float64) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 使用简单的矩形范围查询（可优化为真正的圆形范围�?
	minX := position.X - radius
	maxX := position.X + radius
	minY := position.Y - radius
	maxY := position.Y + radius

	err := r.db.GORMDB().WithContext(ctx).
		Where("pos_x BETWEEN ? AND ?", minX, maxX).
		Where("pos_y BETWEEN ? AND ?", minY, maxY).
		Find(&nodes).Error

	// TODO: 在应用层过滤出真正在圆形范围内的节点
	return nodes, err
}

// GetConnectedNodes 获取与指定节点连接的所有节�?
func (r *nodeRepository) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 通过路径表关联查�?
	err := r.db.GORMDB().WithContext(ctx).
		Joins("JOIN paths ON (nodes.id = paths.start_node_id OR nodes.id = paths.end_node_id)").
		Where("(paths.start_node_id = ? OR paths.end_node_id = ?) AND nodes.id != ?", nodeID, nodeID, nodeID).
		Where("paths.status = ?", domain.PathStatusActive).
		Distinct().
		Find(&nodes).Error

	return nodes, err
}

// GetIsolatedNodes 获取孤立节点（没有连接的节点�?
func (r *nodeRepository) GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 左连接路径表，查找没有路径的节点
	err := r.db.GORMDB().WithContext(ctx).
		Where("NOT EXISTS (SELECT 1 FROM paths WHERE nodes.id = paths.start_node_id OR nodes.id = paths.end_node_id)").
		Find(&nodes).Error

	return nodes, err
}

// GetByLabels 根据标签查询节点
func (r *nodeRepository) GetByLabels(ctx context.Context, labels map[string]string) ([]*domain.Node, error) {
	var nodes []*domain.Node

	query := r.db.GORMDB().WithContext(ctx)

	// 使用JSON查询（需要数据库支持�?
	for key, value := range labels {
		query = query.Where("JSON_EXTRACT(labels, ?) = ?", "$."+key, value)
	}

	err := query.Find(&nodes).Error
	return nodes, err
}

// GetByType 根据类型查询节点
func (r *nodeRepository) GetByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error) {
	var nodes []*domain.Node
	err := r.db.GORMDB().WithContext(ctx).Where("type = ?", nodeType).Find(&nodes).Error
	return nodes, err
}

// GetByStatus 根据状态查询节�?
func (r *nodeRepository) GetByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error) {
	var nodes []*domain.Node
	err := r.db.GORMDB().WithContext(ctx).Where("status = ?", status).Find(&nodes).Error
	return nodes, err
}

// applyFilter 应用查询过滤�?
func (r *nodeRepository) applyFilter(query *gorm.DB, filter NodeFilter) *gorm.DB {
	// ID过滤
	if len(filter.IDs) > 0 {
		stringIDs := make([]string, len(filter.IDs))
		for i, id := range filter.IDs {
			stringIDs[i] = string(id)
		}
		query = query.Where("id IN ?", stringIDs)
	}

	// 名称过滤（模糊查询）
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	// 类型过滤
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	// 状态过�?
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 位置过滤
	if filter.MinX != nil {
		query = query.Where("pos_x >= ?", *filter.MinX)
	}
	if filter.MaxX != nil {
		query = query.Where("pos_x <= ?", *filter.MaxX)
	}
	if filter.MinY != nil {
		query = query.Where("pos_y >= ?", *filter.MinY)
	}
	if filter.MaxY != nil {
		query = query.Where("pos_y <= ?", *filter.MaxY)
	}
	if filter.MinZ != nil {
		query = query.Where("pos_z >= ?", *filter.MinZ)
	}
	if filter.MaxZ != nil {
		query = query.Where("pos_z <= ?", *filter.MaxZ)
	}

	// 标签过滤
	for key, value := range filter.Labels {
		query = query.Where("JSON_EXTRACT(labels, ?) = ?", "$."+key, value)
	}

	return query
}
