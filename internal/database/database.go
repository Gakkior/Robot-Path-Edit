// Package database 提供数据库访问的抽象层
//
// 设计参考：
// - GORM的数据库抽象设计
// - Kubernetes的存储抽象层
// - Docker的存储驱动模式
//
// 特点：
// 1. 数据库无关：支持多种数据库类型
// 2. 连接池管理：自动管理连接生命周期
// 3. 事务支持：支持事务操作
// 4. 迁移管理：自动化数据库结构迁移
package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"robot-path-editor/internal/config"
	"robot-path-editor/internal/domain"
)

// Database 数据库接口抽象
// 提供统一的数据库访问接口，支持不同的数据库实现
type Database interface {
	Connect() error
	Close() error

	DB() interface{}  // 修改为interface{}以支持不同的数据库实现
	GORMDB() *gorm.DB // 添加GORM特定的访问器

	// 健康检查
	Ping() error

	// 数据库操作
	AutoMigrate(dst ...interface{}) error
	Transaction(ctx context.Context, fn func(tx interface{}) error) error
}

// database 数据库实现
type database struct {
	db     *gorm.DB
	config config.DatabaseConfig
}

// New 创建数据库实例
func New(cfg config.DatabaseConfig) (Database, error) {
	db := &database{
		config: cfg,
	}

	if err := db.Connect(); err != nil {
		return nil, err
	}

	return db, nil
}

// Connect 连接数据库
func (d *database) Connect() error {
	var dialector gorm.Dialector

	// 根据数据库类型选择驱动 - 适配器模式
	switch d.config.Type {
	case "sqlite":
		dialector = sqlite.Open(d.config.DSN)
	case "mysql":
		dialector = mysql.Open(d.config.DSN)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", d.config.Type)
	}

	// GORM配置 - 参考最佳实践
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 使用自定义日志
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 建立数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层SQL DB实例进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 配置连接池 - 参考数据库连接池最佳实践
	sqlDB.SetMaxIdleConns(d.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(d.config.ConnMaxLifetime) * time.Second)

	d.db = db

	// 执行健康检查
	if err := d.Ping(); err != nil {
		return fmt.Errorf("数据库连接检查失败: %w", err)
	}

	// 自动迁移数据库结构
	if err := d.AutoMigrate(
		&domain.Node{},
		&domain.Path{},
		&domain.DatabaseConnection{},
		&domain.TableMapping{},
		&domain.Template{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// DB 返回GORM数据库实例
func (d *database) DB() interface{} {
	return d.db
}

// GORMDB 返回GORM数据库实例
func (d *database) GORMDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *database) Close() error {
	if d.db == nil {
		return nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping 检查数据库连接状态
func (d *database) Ping() error {
	if d.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// AutoMigrate 自动迁移数据库结构
func (d *database) AutoMigrate(dst ...interface{}) error {
	if d.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	return d.db.AutoMigrate(dst...)
}

// Transaction 执行事务
func (d *database) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	if d.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
