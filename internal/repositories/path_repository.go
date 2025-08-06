// Package repositories 路径仓储实现
package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// PathRepository 路径仓储接口
type PathRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, path *domain.Path) error
	GetByID(ctx context.Context, id domain.PathID) (*domain.Path, error)
	Update(ctx context.Context, path *domain.Path) error
	Delete(ctx context.Context, id domain.PathID) error

	// 批量操作
	CreateBatch(ctx context.Context, paths []*domain.Path) error
	GetByIDs(ctx context.Context, ids []domain.PathID) ([]*domain.Path, error)

	// 查询操作
	List(ctx context.Context, options PathListOptions) ([]*domain.Path, error)
	Count(ctx context.Context, filter PathFilter) (int64, error)

	// 关系查询
	GetByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
	GetByNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error)
	GetConnectedPaths(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
}

// PathFilter 路径查询过滤器
type PathFilter struct {
	IDs         []domain.PathID      `json:"ids,omitempty"`
	Name        string               `json:"name,omitempty"`
	Type        domain.PathType      `json:"type,omitempty"`
	Status      domain.PathStatus    `json:"status,omitempty"`
	StartNodeID domain.NodeID        `json:"start_node_id,omitempty"`
	EndNodeID   domain.NodeID        `json:"end_node_id,omitempty"`
	Direction   domain.PathDirection `json:"direction,omitempty"`
}

// PathListOptions 路径列表查询选项
type PathListOptions struct {
	Filter   PathFilter `json:"filter"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	OrderBy  string     `json:"order_by"`
	Order    string     `json:"order"`
}

// pathRepository 路径仓储实现
type pathRepository struct {
	db database.Database
}

// NewPathRepository 创建路径仓储实例
func NewPathRepository(db database.Database) PathRepository {
	return &pathRepository{
		db: db,
	}
}

// Create 创建路径
func (r *pathRepository) Create(ctx context.Context, path *domain.Path) error {
	if err := path.IsValid(); err != nil {
		return fmt.Errorf("路径验证失败: %w", err)
	}

	return r.db.GORMDB().WithContext(ctx).Create(path).Error
}

// GetByID 根据ID获取路径
func (r *pathRepository) GetByID(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	var path domain.Path
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&path).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("路径不存�? %s", id)
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
		return fmt.Errorf("路径不存�? %s", path.ID)
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
		return fmt.Errorf("路径不存�? %s", id)
	}

	return nil
}

// CreateBatch 批量创建路径
func (r *pathRepository) CreateBatch(ctx context.Context, paths []*domain.Path) error {
	for _, path := range paths {
		if err := path.IsValid(); err != nil {
			return fmt.Errorf("路径验证失败: %w", err)
		}
	}

	return r.db.Transaction(ctx, func(tx interface{}) error {
		gormTx := tx.(*gorm.DB)
		return gormTx.WithContext(ctx).CreateInBatches(paths, 100).Error
	})
}

// GetByIDs 根据ID列表获取路径
func (r *pathRepository) GetByIDs(ctx context.Context, ids []domain.PathID) ([]*domain.Path, error) {
	var paths []*domain.Path

	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.GORMDB().WithContext(ctx).Where("id IN ?", stringIDs).Find(&paths).Error
	return paths, err
}

// List 列表查询路径
func (r *pathRepository) List(ctx context.Context, options PathListOptions) ([]*domain.Path, error) {
	var paths []*domain.Path

	query := r.db.GORMDB().WithContext(ctx)
	query = r.applyPathFilter(query, options.Filter)

	// 应用排序
	if options.OrderBy != "" {
		order := "asc"
		if options.Order == "desc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", options.OrderBy, order))
	} else {
		query = query.Order("created_at desc")
	}

	// 应用分页
	if options.PageSize > 0 {
		offset := 0
		if options.Page > 1 {
			offset = (options.Page - 1) * options.PageSize
		}
		query = query.Offset(offset).Limit(options.PageSize)
	}

	err := query.Find(&paths).Error
	return paths, err
}

// Count 统计路径数量
func (r *pathRepository) Count(ctx context.Context, filter PathFilter) (int64, error) {
	var count int64

	query := r.db.GORMDB().WithContext(ctx).Model(&domain.Path{})
	query = r.applyPathFilter(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// GetByNode 获取与指定节点相关的所有路�?
func (r *pathRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	var paths []*domain.Path

	err := r.db.GORMDB().WithContext(ctx).
		Where("start_node_id = ? OR end_node_id = ?", nodeID, nodeID).
		Where("status = ?", domain.PathStatusActive).
		Find(&paths).Error

	return paths, err
}

// GetByNodes 获取连接两个节点的路�?
func (r *pathRepository) GetByNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error) {
	var paths []*domain.Path

	err := r.db.GORMDB().WithContext(ctx).
		Where("(start_node_id = ? AND end_node_id = ?) OR (start_node_id = ? AND end_node_id = ? AND direction = ?)",
			startNodeID, endNodeID, endNodeID, startNodeID, domain.PathDirectionBidirectional).
		Where("status = ?", domain.PathStatusActive).
		Find(&paths).Error

	return paths, err
}

// GetConnectedPaths 获取与指定节点连接的路径
func (r *pathRepository) GetConnectedPaths(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	return r.GetByNode(ctx, nodeID)
}

// applyPathFilter 应用路径过滤�?
func (r *pathRepository) applyPathFilter(query *gorm.DB, filter PathFilter) *gorm.DB {
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

	if filter.Direction != "" {
		query = query.Where("direction = ?", filter.Direction)
	}

	return query
}
