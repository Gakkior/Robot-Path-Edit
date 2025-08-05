// Package repositories 表映射仓储实�?
package repositories

import (
	"context"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// TableMappingRepository 表映射仓储接�?
type TableMappingRepository interface {
	Create(ctx context.Context, mapping *domain.TableMapping) error
	GetByID(ctx context.Context, id string) (*domain.TableMapping, error)
	Update(ctx context.Context, mapping *domain.TableMapping) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.TableMapping, error)
	GetByConnectionID(ctx context.Context, connectionID string) ([]*domain.TableMapping, error)
}

type tableMappingRepository struct {
	db database.Database
}

func NewTableMappingRepository(db database.Database) TableMappingRepository {
	return &tableMappingRepository{db: db}
}

func (r *tableMappingRepository) Create(ctx context.Context, mapping *domain.TableMapping) error {
	return r.db.GORMDB().WithContext(ctx).Create(mapping).Error
}

func (r *tableMappingRepository) GetByID(ctx context.Context, id string) (*domain.TableMapping, error) {
	var mapping domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&mapping).Error
	return &mapping, err
}

func (r *tableMappingRepository) Update(ctx context.Context, mapping *domain.TableMapping) error {
	return r.db.GORMDB().WithContext(ctx).Save(mapping).Error
}

func (r *tableMappingRepository) Delete(ctx context.Context, id string) error {
	return r.db.GORMDB().WithContext(ctx).Delete(&domain.TableMapping{}, "id = ?", id).Error
}

func (r *tableMappingRepository) List(ctx context.Context) ([]*domain.TableMapping, error) {
	var mappings []*domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Find(&mappings).Error
	return mappings, err
}

func (r *tableMappingRepository) GetByConnectionID(ctx context.Context, connectionID string) ([]*domain.TableMapping, error) {
	var mappings []*domain.TableMapping
	err := r.db.GORMDB().WithContext(ctx).Where("connection_id = ?", connectionID).Find(&mappings).Error
	return mappings, err
}
