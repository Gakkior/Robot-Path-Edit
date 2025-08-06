// Package repositories 表映射仓储实现
package repositories

import (
	"context"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// TableMappingRepository 表映射仓储接口
type TableMappingRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, mapping *domain.TableMapping) error
	GetByID(ctx context.Context, id string) (*domain.TableMapping, error)
	Update(ctx context.Context, mapping *domain.TableMapping) error
	Delete(ctx context.Context, id string) error

	// 查询操作
	List(ctx context.Context) ([]*domain.TableMapping, error)
	GetByTableName(ctx context.Context, tableName string) (*domain.TableMapping, error)
	GetByConnectionID(ctx context.Context, connectionID string) ([]*domain.TableMapping, error)
}

// tableMappingRepository GORM实现
type tableMappingRepository struct {
	db database.Database
}

// NewTableMappingRepository 创建新的表映射仓储实例
func NewTableMappingRepository(db database.Database) TableMappingRepository {
	return &tableMappingRepository{db: db}
}

// Create 创建表映射
func (r *tableMappingRepository) Create(ctx context.Context, mapping *domain.TableMapping) error {
	return r.db.GORMDB().WithContext(ctx).Create(mapping).Error
}

// GetByID 根据ID获取表映射
func (r *tableMappingRepository) GetByID(ctx context.Context, id string) (*domain.TableMapping, error) {
	var mapping domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&mapping).Error
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// Update 更新表映射
func (r *tableMappingRepository) Update(ctx context.Context, mapping *domain.TableMapping) error {
	return r.db.GORMDB().WithContext(ctx).Save(mapping).Error
}

// Delete 删除表映射
func (r *tableMappingRepository) Delete(ctx context.Context, id string) error {
	return r.db.GORMDB().WithContext(ctx).Delete(&domain.TableMapping{}, "id = ?", id).Error
}

// List 列出所有表映射
func (r *tableMappingRepository) List(ctx context.Context) ([]*domain.TableMapping, error) {
	var mappings []*domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Find(&mappings).Error
	return mappings, err
}

// GetByTableName 根据表名获取映射
func (r *tableMappingRepository) GetByTableName(ctx context.Context, tableName string) (*domain.TableMapping, error) {
	var mapping domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Where("table_name = ?", tableName).First(&mapping).Error
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// GetByConnectionID 根据连接ID获取所有映射
func (r *tableMappingRepository) GetByConnectionID(ctx context.Context, connectionID string) ([]*domain.TableMapping, error) {
	var mappings []*domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Where("connection_id = ?", connectionID).Find(&mappings).Error
	return mappings, err
}
