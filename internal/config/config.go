// Package config 提供应用程序配置管理
//
// 设计参考：
// - Grafana的配置管理系统
// - Kubernetes的配置结构
// - Prometheus的配置验证机制
//
// 特点：
// 1. 分层配置：支持文件、环境变量、命令行参数
// 2. 配置验证：启动时验证配置的正确性
// 3. 热重载：支持配置文件变更时自动重载
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用程序主配置结��?
// 参考Grafana的配置结构设��?
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Canvas   CanvasConfig   `mapstructure:"canvas"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// ServerConfig HTTP服务器配��?
type ServerConfig struct {
	Addr         string        `mapstructure:"addr"`          // 监听地址
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 读取超时
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // 写入超时
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`  // 空闲超时
	TLS          TLSConfig     `mapstructure:"tls"`           // TLS配置
}

// TLSConfig TLS/SSL配置
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// DatabaseConfig 数据库配��?
// 支持多种数据库类型的通用配置
type DatabaseConfig struct {
	Type            string        `mapstructure:"type"`              // 数据库类��? sqlite, mysql, postgres
	DSN             string        `mapstructure:"dsn"`               // 数据源名��?
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大打开连接��?
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周��?
	AutoMigrate     bool          `mapstructure:"auto_migrate"`      // 是否自动迁移
}

// LoggerConfig 日志配置
// 参考Kubernetes的结构化日志设计
type LoggerConfig struct {
	Level      string `mapstructure:"level"`       // 日志级别
	Format     string `mapstructure:"format"`      // 日志格式: json, text
	Output     string `mapstructure:"output"`      // 输出目标: stdout, file
	File       string `mapstructure:"file"`        // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个文件最大大��?MB)
	MaxBackups int    `mapstructure:"max_backups"` // 保留的备份文件数
	MaxAge     int    `mapstructure:"max_age"`     // 保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩备份文件
}

// CanvasConfig 画布相关配置
type CanvasConfig struct {
	DefaultWidth  int     `mapstructure:"default_width"`  // 默认画布宽度
	DefaultHeight int     `mapstructure:"default_height"` // 默认画布高度
	GridSize      int     `mapstructure:"grid_size"`      // 网格大小
	ZoomMin       float64 `mapstructure:"zoom_min"`       // 最小缩放倍数
	ZoomMax       float64 `mapstructure:"zoom_max"`       // 最大缩放倍数
	NodeRadius    int     `mapstructure:"node_radius"`    // 默认节点半径
}

// MetricsConfig 监控指标配置
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"` // 是否启用指标收集
	Path    string `mapstructure:"path"`    // 指标接口路径
}

// Load 加载配置
// 参考Viper的最佳实践和Grafana的配置加载流��?
func Load() (*Config, error) {
	// 设置配置文件搜索路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/robot-path-editor")

	// 设置环境变量前缀 - 参考Kubernetes的环境变量命��?
	viper.SetEnvPrefix("RPE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 设置默认��?- 参考Grafana的默认配��?
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		// 配置文件不存在时继续，使用默认配��?
	}

	// 解析配置到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认配置��?
// 参考各个优秀项目的默认配��?
func setDefaults() {
	// 服务器配置默认��?
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.tls.enabled", false)

	// 数据库配置默认��?
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.dsn", "./data/app.db")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "1h")
	viper.SetDefault("database.auto_migrate", true)

	// 日志配置默认��?
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.output", "stdout")
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_backups", 3)
	viper.SetDefault("logger.max_age", 7)
	viper.SetDefault("logger.compress", true)

	// 画布配置默认��?
	viper.SetDefault("canvas.default_width", 1920)
	viper.SetDefault("canvas.default_height", 1080)
	viper.SetDefault("canvas.grid_size", 20)
	viper.SetDefault("canvas.zoom_min", 0.1)
	viper.SetDefault("canvas.zoom_max", 5.0)
	viper.SetDefault("canvas.node_radius", 20)

	// 监控配置默认��?
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
}

// Validate 验证配置的有效��?
// 参考Kubernetes的配置验证机��?
func (c *Config) Validate() error {
	// 验证服务器配��?
	if c.Server.Addr == "" {
		return fmt.Errorf("server.addr 不能为空")
	}

	// 验证数据库配��?
	if c.Database.Type == "" {
		return fmt.Errorf("database.type 不能为空")
	}

	validDBTypes := map[string]bool{
		"sqlite":   true,
		"mysql":    true,
		"postgres": true,
	}
	if !validDBTypes[c.Database.Type] {
		return fmt.Errorf("不支持的数据库类��? %s", c.Database.Type)
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("database.dsn 不能为空")
	}

	// 验证日志配置
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
		"panic": true,
	}
	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("无效的日志级��? %s", c.Logger.Level)
	}

	// 验证画布配置
	if c.Canvas.ZoomMin <= 0 || c.Canvas.ZoomMax <= 0 || c.Canvas.ZoomMin >= c.Canvas.ZoomMax {
		return fmt.Errorf("无效的缩放配��? min=%f, max=%f", c.Canvas.ZoomMin, c.Canvas.ZoomMax)
	}

	return nil
}
