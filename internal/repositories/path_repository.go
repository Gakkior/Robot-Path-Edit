// Package repositories 璺緞浠撳偍瀹炵幇
package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// PathRepository 璺緞浠撳偍鎺ュ彛
type PathRepository interface {
	// 鍩虹CRUD鎿嶄綔
	Create(ctx context.Context, path *domain.Path) error
	GetByID(ctx context.Context, id domain.PathID) (*domain.Path, error)
	Update(ctx context.Context, path *domain.Path) error
	Delete(ctx context.Context, id domain.PathID) error

	// 鎵归噺鎿嶄綔
	CreateBatch(ctx context.Context, paths []*domain.Path) error
	GetByIDs(ctx context.Context, ids []domain.PathID) ([]*domain.Path, error)

	// 鏌ヨ鎿嶄綔
	List(ctx context.Context, options PathListOptions) ([]*domain.Path, error)
	Count(ctx context.Context, filter PathFilter) (int64, error)

	// 鍏崇郴鏌ヨ
	GetByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
	GetByNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error)
	GetConnectedPaths(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
}

// PathFilter 璺緞鏌ヨ杩囨护鍣?
type PathFilter struct {
	IDs         []domain.PathID      `json:"ids,omitempty"`
	Name        string               `json:"name,omitempty"`
	Type        domain.PathType      `json:"type,omitempty"`
	Status      domain.PathStatus    `json:"status,omitempty"`
	StartNodeID domain.NodeID        `json:"start_node_id,omitempty"`
	EndNodeID   domain.NodeID        `json:"end_node_id,omitempty"`
	Direction   domain.PathDirection `json:"direction,omitempty"`
}

// PathListOptions 璺緞鍒楄〃鏌ヨ閫夐」
type PathListOptions struct {
	Filter   PathFilter `json:"filter"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	OrderBy  string     `json:"order_by"`
	Order    string     `json:"order"`
}

// pathRepository 璺緞浠撳偍瀹炵幇
type pathRepository struct {
	db database.Database
}

// NewPathRepository 鍒涘缓璺緞浠撳偍瀹炰緥
func NewPathRepository(db database.Database) PathRepository {
	return &pathRepository{
		db: db,
	}
}

// Create 鍒涘缓璺緞
func (r *pathRepository) Create(ctx context.Context, path *domain.Path) error {
	if err := path.IsValid(); err != nil {
		return fmt.Errorf("璺緞楠岃瘉澶辫触: %w", err)
	}

	return r.db.GORMDB().WithContext(ctx).Create(path).Error
}

// GetByID 鏍规嵁ID鑾峰彇璺緞
func (r *pathRepository) GetByID(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	var path domain.Path
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&path).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("璺緞涓嶅瓨鍦? %s", id)
		}
		return nil, err
	}
	return &path, nil
}

// Update 鏇存柊璺緞
func (r *pathRepository) Update(ctx context.Context, path *domain.Path) error {
	if err := path.IsValid(); err != nil {
		return fmt.Errorf("璺緞楠岃瘉澶辫触: %w", err)
	}

	result := r.db.GORMDB().WithContext(ctx).Save(path)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("璺緞涓嶅瓨鍦? %s", path.ID)
	}

	return nil
}

// Delete 鍒犻櫎璺緞
func (r *pathRepository) Delete(ctx context.Context, id domain.PathID) error {
	result := r.db.GORMDB().WithContext(ctx).Delete(&domain.Path{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("璺緞涓嶅瓨鍦? %s", id)
	}

	return nil
}

// CreateBatch 鎵归噺鍒涘缓璺緞
func (r *pathRepository) CreateBatch(ctx context.Context, paths []*domain.Path) error {
	for _, path := range paths {
		if err := path.IsValid(); err != nil {
			return fmt.Errorf("璺緞楠岃瘉澶辫触: %w", err)
		}
	}

	return r.db.Transaction(ctx, func(tx interface{}) error {
		gormTx := tx.(*gorm.DB)
		return gormTx.WithContext(ctx).CreateInBatches(paths, 100).Error
	})
}

// GetByIDs 鏍规嵁ID鍒楄〃鑾峰彇璺緞
func (r *pathRepository) GetByIDs(ctx context.Context, ids []domain.PathID) ([]*domain.Path, error) {
	var paths []*domain.Path

	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.GORMDB().WithContext(ctx).Where("id IN ?", stringIDs).Find(&paths).Error
	return paths, err
}

// List 鍒楄〃鏌ヨ璺緞
func (r *pathRepository) List(ctx context.Context, options PathListOptions) ([]*domain.Path, error) {
	var paths []*domain.Path

	query := r.db.GORMDB().WithContext(ctx)
	query = r.applyPathFilter(query, options.Filter)

	// 搴旂敤鎺掑簭
	if options.OrderBy != "" {
		order := "asc"
		if options.Order == "desc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", options.OrderBy, order))
	} else {
		query = query.Order("created_at desc")
	}

	// 搴旂敤鍒嗛〉
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

// Count 缁熻璺緞鏁伴噺
func (r *pathRepository) Count(ctx context.Context, filter PathFilter) (int64, error) {
	var count int64

	query := r.db.GORMDB().WithContext(ctx).Model(&domain.Path{})
	query = r.applyPathFilter(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// GetByNode 鑾峰彇涓庢寚瀹氳妭鐐圭浉鍏崇殑鎵€鏈夎矾寰?
func (r *pathRepository) GetByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	var paths []*domain.Path

	err := r.db.GORMDB().WithContext(ctx).
		Where("start_node_id = ? OR end_node_id = ?", nodeID, nodeID).
		Where("status = ?", domain.PathStatusActive).
		Find(&paths).Error

	return paths, err
}

// GetByNodes 鑾峰彇杩炴帴涓や釜鑺傜偣鐨勮矾寰?
func (r *pathRepository) GetByNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error) {
	var paths []*domain.Path

	err := r.db.GORMDB().WithContext(ctx).
		Where("(start_node_id = ? AND end_node_id = ?) OR (start_node_id = ? AND end_node_id = ? AND direction = ?)",
			startNodeID, endNodeID, endNodeID, startNodeID, domain.PathDirectionBidirectional).
		Where("status = ?", domain.PathStatusActive).
		Find(&paths).Error

	return paths, err
}

// GetConnectedPaths 鑾峰彇涓庢寚瀹氳妭鐐硅繛鎺ョ殑璺緞
func (r *pathRepository) GetConnectedPaths(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	return r.GetByNode(ctx, nodeID)
}

// applyPathFilter 搴旂敤璺緞杩囨护鍣?
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
