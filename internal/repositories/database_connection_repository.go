// Package repositories 数据库连接仓储实现
package repositories

import (
	"context"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// DatabaseConnectionRepository 数据库连接仓储接口
type DatabaseConnectionRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, conn *domain.DatabaseConnection) error
	GetByID(ctx context.Context, id string) (*domain.DatabaseConnection, error)
	Update(ctx context.Context, conn *domain.DatabaseConnection) error
	Delete(ctx context.Context, id string) error

	// 查询操作
	List(ctx context.Context) ([]*domain.DatabaseConnection, error)
	GetByType(ctx context.Context, dbType string) ([]*domain.DatabaseConnection, error)

	// 连接测试
	TestConnection(ctx context.Context, id string) error
}

// databaseConnectionRepository GORM实现
type databaseConnectionRepository struct {
	db database.Database
}

// NewDatabaseConnectionRepository 创建新的数据库连接仓储实例
func NewDatabaseConnectionRepository(db database.Database) DatabaseConnectionRepository {
	return &databaseConnectionRepository{db: db}
}

// Create 创建数据库连接配置
func (r *databaseConnectionRepository) Create(ctx context.Context, conn *domain.DatabaseConnection) error {
	return r.db.GORMDB().WithContext(ctx).Create(conn).Error
}

// GetByID 根据ID获取数据库连接配置
func (r *databaseConnectionRepository) GetByID(ctx context.Context, id string) (*domain.DatabaseConnection, error) {
	var conn domain.DatabaseConnection
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&conn).Error
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

// Update 更新数据库连接配置
func (r *databaseConnectionRepository) Update(ctx context.Context, conn *domain.DatabaseConnection) error {
	return r.db.GORMDB().WithContext(ctx).Save(conn).Error
}

// Delete 删除数据库连接配置
func (r *databaseConnectionRepository) Delete(ctx context.Context, id string) error {
	return r.db.GORMDB().WithContext(ctx).Delete(&domain.DatabaseConnection{}, "id = ?", id).Error
}

// List 列出所有数据库连接配置
func (r *databaseConnectionRepository) List(ctx context.Context) ([]*domain.DatabaseConnection, error) {
	var connections []*domain.DatabaseConnection
	err := r.db.GORMDB().WithContext(ctx).Find(&connections).Error
	return connections, err
}

// GetByType 根据数据库类型获取连接配置
func (r *databaseConnectionRepository) GetByType(ctx context.Context, dbType string) ([]*domain.DatabaseConnection, error) {
	var connections []*domain.DatabaseConnection
	err := r.db.GORMDB().WithContext(ctx).Where("db_type = ?", dbType).Find(&connections).Error
	return connections, err
}

// TestConnection 测试数据库连接
func (r *databaseConnectionRepository) TestConnection(ctx context.Context, id string) error {
	// 获取连接配置
	conn, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 这里可以实现实际的数据库连接测试逻辑
	// 暂时返回nil表示测试成功
	_ = conn
	return nil
}
