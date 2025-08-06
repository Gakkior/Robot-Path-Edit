// Package repositories 路径仓储实现
//
// 设计参考：
// - DDD的仓储模式
// - 图数据库的路径查询模式
// - Neo4j的关系查询设计
//
// 特点：
// 1. 路径数据访问抽象
// 2. 图关系查询优化
// 3. 复杂路径算法支持
package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// PathRepository 路径仓储接口
// 定义路径数据访问的所有操作
type PathRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, path *domain.Path) error
	GetByID(ctx context.Context, id domain.PathID) (*domain.Path, error)
	Update(ctx context.Context, path *domain.Path) error
	Delete(ctx context.Context, id domain.PathID) error

	// 批量操作
	GetByIDs(ctx context.Context, ids []domain.PathID) ([]*domain.Path, error)
	CreateBatch(ctx context.Context, paths []*domain.Path) error
	DeleteBatch(ctx context.Context, ids []domain.PathID) error

	// 查询操作
	List(ctx context.Context, filter PathFilter) ([]*domain.Path, error)
	Count(ctx context.Context, filter PathFilter) (int64, error)

	// 关系查询
	GetByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
	GetByNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error)
	GetConnectedPaths(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
}

// PathFilter 路径查询过滤器
type PathFilter struct {
	IDs         []domain.PathID   `json:"ids,omitempty"`
	Name        string            `json:"name,omitempty"`
	Type        domain.PathType   `json:"type,omitempty"`
	Status      domain.PathStatus `json:"status,omitempty"`
	StartNodeID domain.NodeID     `json:"start_node_id,omitempty"`
	EndNodeID   domain.NodeID     `json:"end_node_id,omitempty"`

	// 权重范围过滤
	MinWeight *float64 `json:"min_weight,omitempty"`
	MaxWeight *float64 `json:"max_weight,omitempty"`

	// 分页参数
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`

	// 排序参数
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// pathRepository GORM实现
type pathRepository struct {
	db database.Database
}

// NewPathRepository 创建新的路径仓储实例
func NewPathRepository(db database.Database) PathRepository {
	return &pathRepository{db: db}
}

// Create 创建路径
func (r *pathRepository) Create(ctx context.Context, path *domain.Path) error {
	if err := path.IsValid(); err != nil {
		return fmt.Errorf("路径验证失败: %w", err)
	}

	err := r.db.GORMDB().WithContext(ctx).Create(path).Error
	if err != nil {
		return fmt.Errorf("创建路径失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取路径
func (r *pathRepository) GetByID(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	var path domain.Path
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&path).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("路径不存在: %s", id)
		}
		return nil, err
	}
	return &path, nil
}

// Update 更新路径
func (r *pathRepository) Update(ctx context.Context, path *domain.Path) error {
	if err := path.IsValid(); err != nil {
		return fmt.Errorf("路径验证失败: %w", err)
	}

	result := r.db.GORMDB().WithContext(ctx).Save(path)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("路径不存在: %s", path.ID)
	}

	return nil
}

// Delete 删除路径
func (r *pathRepository) Delete(ctx context.Context, id domain.PathID) error {
	result := r.db.GORMDB().WithContext(ctx).Delete(&domain.Path{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("路径不存在: %s", id)
	}

	return nil
}

// GetByIDs 根据ID列表获取路径
func (r *pathRepository) GetByIDs(ctx context.Context, ids []domain.PathID) ([]*domain.Path, error) {
	if len(ids) == 0 {
		return []*domain.Path{}, nil
	}

	var paths []*domain.Path

	// 将PathID转换为字符串
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.GORMDB().WithContext(ctx).Where("id IN ?", stringIDs).Find(&paths).Error
	return paths, err
}

// CreateBatch 批量创建路径
func (r *pathRepository) CreateBatch(ctx context.Context, paths []*domain.Path) error {
	return r.db.Transaction(ctx, func(tx interface{}) error {
		gormTx := tx.(*gorm.DB)
		for _, path := range paths {
			if err := path.IsValid(); err != nil {
				return fmt.Errorf("路径验证失败: %w", err)
			}
			if err := gormTx.Create(path).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteBatch 批量删除路径
func (r *pathRepository) DeleteBatch(ctx context.Context, ids []domain.PathID) error {
	if len(ids) == 0 {
		return nil
	}

	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.GORMDB().WithContext(ctx).Where("id IN ?", stringIDs).Delete(&domain.Path{}).Error
	return err
}

// List 列出路径
func (r *pathRepository) List(ctx context.Context, filter PathFilter) ([]*domain.Path, error) {
	var paths []*domain.Path

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

	if filter.StartNodeID != "" {
		query = query.Where("start_node_id = ?", filter.StartNodeID)
	}

	if filter.EndNodeID != "" {
		query = query.Where("end_node_id = ?", filter.EndNodeID)
	}

	if filter.MinWeight != nil {
		query = query.Where("weight >= ?", *filter.MinWeight)
	}

	if filter.MaxWeight != nil {
		query = query.Where("weight <= ?", *filter.MaxWeight)
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

	err := query.Find(&paths).Error
	return paths, err
}

// Count 统计路径数量
func (r *pathRepository) Count(ctx context.Context, filter PathFilter) (int64, error) {
	var count int64

	query := r.db.GORMDB().WithContext(ctx).Model(&domain.Path{})

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

	if filter.StartNodeID != "" {
		query = query.Where("start_node_id = ?", filter.StartNodeID)
	}

	if filter.EndNodeID != "" {
		query = query.Where("end_node_id = ?", filter.EndNodeID)
	}

	if filter.MinWeight != nil {
		query = query.Where("weight >= ?", *filter.MinWeight)
	}

	if filter.MaxWeight != nil {
		query = query.Where("weight <= ?", *filter.MaxWeight)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetByNode 获取与指定节点相关的所有路径
func (r *pathRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	var paths []*domain.Path

	err := r.db.GORMDB().WithContext(ctx).
		Where("start_node_id = ? OR end_node_id = ?", nodeID, nodeID).
		Find(&paths).Error

	return paths, err
}

// GetByNodes 获取连接两个特定节点的路径
func (r *pathRepository) GetByNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error) {
	var paths []*domain.Path

	err := r.db.GORMDB().WithContext(ctx).
		Where("(start_node_id = ? AND end_node_id = ?) OR (start_node_id = ? AND end_node_id = ?)",
			startNodeID, endNodeID, endNodeID, startNodeID).
		Find(&paths).Error

	return paths, err
}

// GetConnectedPaths 获取与指定节点连接的所有路径
func (r *pathRepository) GetConnectedPaths(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	// 这与GetByNode相同，但可以根据需要添加额外的业务逻辑
	return r.GetByNode(ctx, nodeID)
}
