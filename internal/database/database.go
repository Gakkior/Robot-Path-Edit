// Package database 提供数据库访问的抽象�?
//
// 设计参考：
// - GORM的数据库抽象设计
// - Kubernetes的存储抽象层
// - Docker的存储驱动模�?
//
// 特点�?
// 1. 数据库无关：支持多种数据库类�?
// 2. 连接池管理：自动管理连接生命周期
// 3. 事务支持：支持事务操�?
// 4. 迁移管理：自动化数据库结构迁�?
package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"robot-path-editor/internal/config"
	"robot-path-editor/internal/domain"
)

// Database 数据库接口抽�?
// 提供统一的数据库访问接口，支持不同的数据库实�?
type Database interface {
	// 基础操作
	DB() interface{}  // 修改为interface{}以支持不同的数据库实�?
	GORMDB() *gorm.DB // 添加GORM特定的�闖��?
	Close() error

	// 健康检�?
	Ping(ctx context.Context) error

	// 事务操作
	Transaction(ctx context.Context, fn func(tx interface{}) error) error

	// 迁移操作
	AutoMigrate() error
}

// database 数据库实�?
type database struct {
	db     *gorm.DB
	config config.DatabaseConfig
}

// New 创建数据库实�?
// 根据配置类型自动选择对应的数据库驱动
func New(cfg config.DatabaseConfig) (Database, error) {
	// 如果是SQLite，使用内存数据库作为后备方案
	if cfg.Type == "sqlite" {
		fmt.Printf("信息: 使用内存数据库模式\n")
		return NewMemoryDatabaseFromConfig(cfg)
	}

	var dialector gorm.Dialector

	// 根据数据库类型选择驱动 - 适配器模�?
	switch cfg.Type {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("不支持的数据库类�? %s", cfg.Type)
	}

	// GORM配置 - 参考最佳实�?
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 使用自定义日�?
		NowFunc: func() time.Time {
			return time.Now().UTC() // 统一使用UTC时间
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 支持SQLite
	}

	// 建立数据库连�?
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失�? %w", err)
	}

	// 获取底层SQL DB实例进行连接池配�?
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取SQL DB实例失败: %w", err)
	}

	// 配置连接�?- 参考数据库连接池最佳实�?
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	instance := &database{
		db:     db,
		config: cfg,
	}

	// 自动迁移数据库结�?
	if cfg.AutoMigrate {
		if err := instance.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("数据库迁移失�? %w", err)
		}
	}

	return instance, nil
}

// DB 返回GORM数据库实�?
func (d *database) DB() interface{} {
	return d.db
}

// GORMDB 返回GORM数据库实�?
func (d *database) GORMDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连�?
func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 检查数据库连接状�?
func (d *database) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Transaction 执行事务操作
// 参考GORM的事务最佳实�?
func (d *database) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// AutoMigrate 自动迁移数据库结�?
// 参考Kubernetes的声明式资源管理
func (d *database) AutoMigrate() error {
	// 定义需要迁移的模型
	models := []interface{}{
		&domain.Node{},
		&domain.Path{},
		&domain.DatabaseConnection{},
		&domain.TableMapping{},
	}

	// 执行迁移
	for _, model := range models {
		if err := d.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型 %T 失败: %w", model, err)
		}
	}

	return nil
}
