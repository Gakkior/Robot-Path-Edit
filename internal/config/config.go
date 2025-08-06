// Package config 应用程序配置管理
//
// 设计参考：
// - Viper配置管理最佳实践
// - Grafana的配置结构设计
// - Kubernetes的配置加载机制
//
// 特点：
// 1. 多数据源支持：文件、环境变量、命令行参数
// 2. 配置验证：确保配置的有效性
// 3. 热重载：支持配置动态更新（可选）
// 4. 结构化配置：使用struct tag进行配置映射
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用程序主配置结构
// 参考Grafana的配置结构设计
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Canvas   CanvasConfig   `mapstructure:"canvas"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Host            string        `mapstructure:"host"`             // 服务监听地址
	Port            int           `mapstructure:"port"`             // 服务监听端口
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // 读取超时
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // 写入超时
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`     // 空闲超时
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // 关闭超时
	EnableCORS      bool          `mapstructure:"enable_cors"`      // 启用CORS
	EnableHTTPS     bool          `mapstructure:"enable_https"`     // 启用HTTPS
	TLSCertPath     string        `mapstructure:"tls_cert_path"`    // TLS证书路径
	TLSKeyPath      string        `mapstructure:"tls_key_path"`     // TLS私钥路径
	APIPrefix       string        `mapstructure:"api_prefix"`       // API路径前缀
	WebRoot         string        `mapstructure:"web_root"`         // 静态文件根目录
	MaxRequestSize  int64         `mapstructure:"max_request_size"` // 最大请求大小
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string        `mapstructure:"type"`              // 数据库类型: sqlite, mysql, postgres
	DSN             string        `mapstructure:"dsn"`               // 数据源名称
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大打开连接数
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
	Debug           bool          `mapstructure:"debug"`             // 调试模式
	AutoMigrate     bool          `mapstructure:"auto_migrate"`      // 自动迁移
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`       // 日志级别
	Format     string `mapstructure:"format"`      // 日志格式: json, text
	Output     string `mapstructure:"output"`      // 输出目标: stdout, stderr, file
	FilePath   string `mapstructure:"file_path"`   // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个文件最大大小(MB)
	MaxBackups int    `mapstructure:"max_backups"` // 保留备份文件数
	MaxAge     int    `mapstructure:"max_age"`     // 保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

// CanvasConfig 画布配置
type CanvasConfig struct {
	Width      int     `mapstructure:"width"`        // 画布宽度
	Height     int     `mapstructure:"height"`       // 画布高度
	ZoomMin    float64 `mapstructure:"zoom_min"`     // 最小缩放
	ZoomMax    float64 `mapstructure:"zoom_max"`     // 最大缩放
	GridSize   int     `mapstructure:"grid_size"`    // 网格大小
	SnapToGrid bool    `mapstructure:"snap_to_grid"` // 对齐网格
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"` // 启用监控
	Port    int    `mapstructure:"port"`    // 监控端口
	Path    string `mapstructure:"path"`    // 监控路径
}

// Load 加载配置
// 参考Viper的最佳实践和Grafana的配置加载流程
func Load(configPath string) (*Config, error) {
	// 配置viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置环境变量前缀 - 参考Kubernetes的环境变量命名
	viper.SetEnvPrefix("ROBOT_PATH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 设置默认值 - 参考Grafana的默认配置
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在时继续，使用默认配置
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认配置值
// 参考各个优秀项目的默认配置
func setDefaults() {
	// 服务器配置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("server.enable_cors", true)
	viper.SetDefault("server.api_prefix", "/api/v1")
	viper.SetDefault("server.max_request_size", 32<<20) // 32MB

	// 数据库配置默认值
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.dsn", "data.db")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.auto_migrate", true)

	// 日志配置默认值
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.output", "stdout")
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_backups", 3)
	viper.SetDefault("logger.max_age", 28)
	viper.SetDefault("logger.compress", true)

	// 画布配置默认值
	viper.SetDefault("canvas.width", 1200)
	viper.SetDefault("canvas.height", 800)
	viper.SetDefault("canvas.zoom_min", 0.1)
	viper.SetDefault("canvas.zoom_max", 5.0)
	viper.SetDefault("canvas.grid_size", 20)
	viper.SetDefault("canvas.snap_to_grid", true)

	// 监控配置默认值
	viper.SetDefault("metrics.enabled", false)
	viper.SetDefault("metrics.port", 9090)
	viper.SetDefault("metrics.path", "/metrics")
}

// Validate 验证配置的有效性
// 参考Kubernetes的配置验证机制
func (c *Config) Validate() error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("无效的服务器端口: %d", c.Server.Port)
	}

	// 验证数据库配置
	validDBTypes := []string{"sqlite", "mysql", "postgres"}
	isValidDBType := false
	for _, validType := range validDBTypes {
		if c.Database.Type == validType {
			isValidDBType = true
			break
		}
	}
	if !isValidDBType {
		return fmt.Errorf("不支持的数据库类型: %s", c.Database.Type)
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("数据源名称不能为空")
	}

	if c.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("最大打开连接数必须大于0")
	}

	if c.Database.MaxIdleConns <= 0 {
		return fmt.Errorf("最大空闲连接数必须大于0")
	}

	// 验证日志配置
	validLogLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	isValidLogLevel := false
	for _, validLevel := range validLogLevels {
		if c.Logger.Level == validLevel {
			isValidLogLevel = true
			break
		}
	}
	if !isValidLogLevel {
		return fmt.Errorf("无效的日志级别: %s", c.Logger.Level)
	}

	// 验证画布配置
	if c.Canvas.ZoomMin <= 0 || c.Canvas.ZoomMax <= c.Canvas.ZoomMin {
		return fmt.Errorf("无效的缩放配置: min=%f, max=%f", c.Canvas.ZoomMin, c.Canvas.ZoomMax)
	}

	return nil
}
