// Package database 鎻愪緵鏁版嵁搴撹闂殑鎶借薄灞?
//
// 璁捐鍙傝€冿細
// - GORM鐨勬暟鎹簱鎶借薄璁捐
// - Kubernetes鐨勫瓨鍌ㄦ娊璞″眰
// - Docker鐨勫瓨鍌ㄩ┍鍔ㄦā寮?
//
// 鐗圭偣锛?
// 1. 鏁版嵁搴撴棤鍏筹細鏀寔澶氱鏁版嵁搴撶被鍨?
// 2. 杩炴帴姹犵鐞嗭細鑷姩绠＄悊杩炴帴鐢熷懡鍛ㄦ湡
// 3. 浜嬪姟鏀寔锛氭敮鎸佷簨鍔℃搷浣?
// 4. 杩佺Щ绠＄悊锛氳嚜鍔ㄥ寲鏁版嵁搴撶粨鏋勮縼绉?
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

// Database 鏁版嵁搴撴帴鍙ｆ娊璞?
// 鎻愪緵缁熶竴鐨勬暟鎹簱璁块棶鎺ュ彛锛屾敮鎸佷笉鍚岀殑鏁版嵁搴撳疄鐜?
type Database interface {
	// 鍩虹鎿嶄綔
	DB() interface{}  // 淇敼涓篿nterface{}浠ユ敮鎸佷笉鍚岀殑鏁版嵁搴撳疄鐜?
	GORMDB() *gorm.DB // 娣诲姞GORM鐗瑰畾鐨勮闂柟娉?
	Close() error

	// 鍋ュ悍妫€鏌?
	Ping(ctx context.Context) error

	// 浜嬪姟鎿嶄綔
	Transaction(ctx context.Context, fn func(tx interface{}) error) error

	// 杩佺Щ鎿嶄綔
	AutoMigrate() error
}

// database 鏁版嵁搴撳疄鐜?
type database struct {
	db     *gorm.DB
	config config.DatabaseConfig
}

// New 鍒涘缓鏁版嵁搴撳疄渚?
// 鏍规嵁閰嶇疆绫诲瀷鑷姩閫夋嫨瀵瑰簲鐨勬暟鎹簱椹卞姩
func New(cfg config.DatabaseConfig) (Database, error) {
	// 濡傛灉鏄疭QLite锛屼娇鐢ㄥ唴瀛樻暟鎹簱浣滀负鍚庡鏂规
	if cfg.Type == "sqlite" {
		fmt.Printf("淇℃伅: 浣跨敤鍐呭瓨鏁版嵁搴撴ā寮廫n")
		return NewMemoryDatabaseFromConfig(cfg)
	}

	var dialector gorm.Dialector

	// 鏍规嵁鏁版嵁搴撶被鍨嬮€夋嫨椹卞姩 - 閫傞厤鍣ㄦā寮?
	switch cfg.Type {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("涓嶆敮鎸佺殑鏁版嵁搴撶被鍨? %s", cfg.Type)
	}

	// GORM閰嶇疆 - 鍙傝€冩渶浣冲疄璺?
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 浣跨敤鑷畾涔夋棩蹇?
		NowFunc: func() time.Time {
			return time.Now().UTC() // 缁熶竴浣跨敤UTC鏃堕棿
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 鏀寔SQLite
	}

	// 寤虹珛鏁版嵁搴撹繛鎺?
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("杩炴帴鏁版嵁搴撳け璐? %w", err)
	}

	// 鑾峰彇搴曞眰SQL DB瀹炰緥杩涜杩炴帴姹犻厤缃?
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇SQL DB瀹炰緥澶辫触: %w", err)
	}

	// 閰嶇疆杩炴帴姹?- 鍙傝€冩暟鎹簱杩炴帴姹犳渶浣冲疄璺?
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	instance := &database{
		db:     db,
		config: cfg,
	}

	// 鑷姩杩佺Щ鏁版嵁搴撶粨鏋?
	if cfg.AutoMigrate {
		if err := instance.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("鏁版嵁搴撹縼绉诲け璐? %w", err)
		}
	}

	return instance, nil
}

// DB 杩斿洖GORM鏁版嵁搴撳疄渚?
func (d *database) DB() interface{} {
	return d.db
}

// GORMDB 杩斿洖GORM鏁版嵁搴撳疄渚?
func (d *database) GORMDB() *gorm.DB {
	return d.db
}

// Close 鍏抽棴鏁版嵁搴撹繛鎺?
func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 妫€鏌ユ暟鎹簱杩炴帴鐘舵€?
func (d *database) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Transaction 鎵ц浜嬪姟鎿嶄綔
// 鍙傝€僄ORM鐨勪簨鍔℃渶浣冲疄璺?
func (d *database) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// AutoMigrate 鑷姩杩佺Щ鏁版嵁搴撶粨鏋?
// 鍙傝€僈ubernetes鐨勫０鏄庡紡璧勬簮绠＄悊
func (d *database) AutoMigrate() error {
	// 瀹氫箟闇€瑕佽縼绉荤殑妯″瀷
	models := []interface{}{
		&domain.Node{},
		&domain.Path{},
		&domain.DatabaseConnection{},
		&domain.TableMapping{},
	}

	// 鎵ц杩佺Щ
	for _, model := range models {
		if err := d.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("杩佺Щ妯″瀷 %T 澶辫触: %w", model, err)
		}
	}

	return nil
}
