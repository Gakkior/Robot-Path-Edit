// Package repositories 数据库连接仓储实�?
package repositories

import (
	"context"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// DatabaseConnectionRepository 数据库连接仓储接�?
type DatabaseConnectionRepository interface {
	Create(ctx context.Context, conn *domain.DatabaseConnection) error
	GetByID(ctx context.Context, id string) (*domain.DatabaseConnection, error)
	Update(ctx context.Context, conn *domain.DatabaseConnection) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.DatabaseConnection, error)
}

type databaseConnectionRepository struct {
	db database.Database
}

func NewDatabaseConnectionRepository(db database.Database) DatabaseConnectionRepository {
	return &databaseConnectionRepository{db: db}
}

func (r *databaseConnectionRepository) Create(ctx context.Context, conn *domain.DatabaseConnection) error {
	return r.db.GORMDB().WithContext(ctx).Create(conn).Error
}

func (r *databaseConnectionRepository) GetByID(ctx context.Context, id string) (*domain.DatabaseConnection, error) {
	var conn domain.DatabaseConnection
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&conn).Error
	return &conn, err
}

func (r *databaseConnectionRepository) Update(ctx context.Context, conn *domain.DatabaseConnection) error {
	return r.db.GORMDB().WithContext(ctx).Save(conn).Error
}

func (r *databaseConnectionRepository) Delete(ctx context.Context, id string) error {
	return r.db.GORMDB().WithContext(ctx).Delete(&domain.DatabaseConnection{}, "id = ?", id).Error
}

func (r *databaseConnectionRepository) List(ctx context.Context) ([]*domain.DatabaseConnection, error) {
	var conns []*domain.DatabaseConnection
	err := r.db.GORMDB().WithContext(ctx).Find(&conns).Error
	return conns, err
}
