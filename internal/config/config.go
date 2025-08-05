// Package config 鎻愪緵搴旂敤绋嬪簭閰嶇疆绠＄悊
//
// 璁捐鍙傝€冿細
// - Grafana鐨勯厤缃鐞嗙郴缁?
// - Kubernetes鐨勯厤缃粨鏋?
// - Prometheus鐨勯厤缃獙璇佹満鍒?
//
// 鐗圭偣锛?
// 1. 鍒嗗眰閰嶇疆锛氭敮鎸佹枃浠躲€佺幆澧冨彉閲忋€佸懡浠よ鍙傛暟
// 2. 閰嶇疆楠岃瘉锛氬惎鍔ㄦ椂楠岃瘉閰嶇疆鐨勬纭€?
// 3. 鐑噸杞斤細鏀寔閰嶇疆鏂囦欢鍙樻洿鏃惰嚜鍔ㄩ噸杞?
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 搴旂敤绋嬪簭涓婚厤缃粨鏋?
// 鍙傝€僄rafana鐨勯厤缃粨鏋勮璁?
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Canvas   CanvasConfig   `mapstructure:"canvas"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

// ServerConfig HTTP鏈嶅姟鍣ㄩ厤缃?
type ServerConfig struct {
	Addr         string        `mapstructure:"addr"`          // 鐩戝惉鍦板潃
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 璇诲彇瓒呮椂
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // 鍐欏叆瓒呮椂
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`  // 绌洪棽瓒呮椂
	TLS          TLSConfig     `mapstructure:"tls"`           // TLS閰嶇疆
}

// TLSConfig TLS/SSL閰嶇疆
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// DatabaseConfig 鏁版嵁搴撻厤缃?
// 鏀寔澶氱鏁版嵁搴撶被鍨嬬殑閫氱敤閰嶇疆
type DatabaseConfig struct {
	Type            string        `mapstructure:"type"`              // 鏁版嵁搴撶被鍨? sqlite, mysql, postgres
	DSN             string        `mapstructure:"dsn"`               // 鏁版嵁婧愬悕绉?
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 鏈€澶ф墦寮€杩炴帴鏁?
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 鏈€澶х┖闂茶繛鎺ユ暟
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 杩炴帴鏈€澶х敓鍛藉懆鏈?
	AutoMigrate     bool          `mapstructure:"auto_migrate"`      // 鏄惁鑷姩杩佺Щ
}

// LoggerConfig 鏃ュ織閰嶇疆
// 鍙傝€僈ubernetes鐨勭粨鏋勫寲鏃ュ織璁捐
type LoggerConfig struct {
	Level      string `mapstructure:"level"`       // 鏃ュ織绾у埆
	Format     string `mapstructure:"format"`      // 鏃ュ織鏍煎紡: json, text
	Output     string `mapstructure:"output"`      // 杈撳嚭鐩爣: stdout, file
	File       string `mapstructure:"file"`        // 鏃ュ織鏂囦欢璺緞
	MaxSize    int    `mapstructure:"max_size"`    // 鍗曚釜鏂囦欢鏈€澶уぇ灏?MB)
	MaxBackups int    `mapstructure:"max_backups"` // 淇濈暀鐨勫浠芥枃浠舵暟
	MaxAge     int    `mapstructure:"max_age"`     // 淇濈暀澶╂暟
	Compress   bool   `mapstructure:"compress"`    // 鏄惁鍘嬬缉澶囦唤鏂囦欢
}

// CanvasConfig 鐢诲竷鐩稿叧閰嶇疆
type CanvasConfig struct {
	DefaultWidth  int     `mapstructure:"default_width"`  // 榛樿鐢诲竷瀹藉害
	DefaultHeight int     `mapstructure:"default_height"` // 榛樿鐢诲竷楂樺害
	GridSize      int     `mapstructure:"grid_size"`      // 缃戞牸澶у皬
	ZoomMin       float64 `mapstructure:"zoom_min"`       // 鏈€灏忕缉鏀惧€嶆暟
	ZoomMax       float64 `mapstructure:"zoom_max"`       // 鏈€澶х缉鏀惧€嶆暟
	NodeRadius    int     `mapstructure:"node_radius"`    // 榛樿鑺傜偣鍗婂緞
}

// MetricsConfig 鐩戞帶鎸囨爣閰嶇疆
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"` // 鏄惁鍚敤鎸囨爣鏀堕泦
	Path    string `mapstructure:"path"`    // 鎸囨爣鎺ュ彛璺緞
}

// Load 鍔犺浇閰嶇疆
// 鍙傝€僔iper鐨勬渶浣冲疄璺靛拰Grafana鐨勯厤缃姞杞芥祦绋?
func Load() (*Config, error) {
	// 璁剧疆閰嶇疆鏂囦欢鎼滅储璺緞
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/robot-path-editor")

	// 璁剧疆鐜鍙橀噺鍓嶇紑 - 鍙傝€僈ubernetes鐨勭幆澧冨彉閲忓懡鍚?
	viper.SetEnvPrefix("RPE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 璁剧疆榛樿鍊?- 鍙傝€僄rafana鐨勯粯璁ら厤缃?
	setDefaults()

	// 璇诲彇閰嶇疆鏂囦欢
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("璇诲彇閰嶇疆鏂囦欢澶辫触: %w", err)
		}
		// 閰嶇疆鏂囦欢涓嶅瓨鍦ㄦ椂缁х画锛屼娇鐢ㄩ粯璁ら厤缃?
	}

	// 瑙ｆ瀽閰嶇疆鍒扮粨鏋勪綋
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("瑙ｆ瀽閰嶇疆澶辫触: %w", err)
	}

	// 楠岃瘉閰嶇疆
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("閰嶇疆楠岃瘉澶辫触: %w", err)
	}

	return &cfg, nil
}

// setDefaults 璁剧疆榛樿閰嶇疆鍊?
// 鍙傝€冨悇涓紭绉€椤圭洰鐨勯粯璁ら厤缃?
func setDefaults() {
	// 鏈嶅姟鍣ㄩ厤缃粯璁ゅ€?
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.tls.enabled", false)

	// 鏁版嵁搴撻厤缃粯璁ゅ€?
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.dsn", "./data/app.db")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "1h")
	viper.SetDefault("database.auto_migrate", true)

	// 鏃ュ織閰嶇疆榛樿鍊?
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.output", "stdout")
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_backups", 3)
	viper.SetDefault("logger.max_age", 7)
	viper.SetDefault("logger.compress", true)

	// 鐢诲竷閰嶇疆榛樿鍊?
	viper.SetDefault("canvas.default_width", 1920)
	viper.SetDefault("canvas.default_height", 1080)
	viper.SetDefault("canvas.grid_size", 20)
	viper.SetDefault("canvas.zoom_min", 0.1)
	viper.SetDefault("canvas.zoom_max", 5.0)
	viper.SetDefault("canvas.node_radius", 20)

	// 鐩戞帶閰嶇疆榛樿鍊?
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
}

// Validate 楠岃瘉閰嶇疆鐨勬湁鏁堟€?
// 鍙傝€僈ubernetes鐨勯厤缃獙璇佹満鍒?
func (c *Config) Validate() error {
	// 楠岃瘉鏈嶅姟鍣ㄩ厤缃?
	if c.Server.Addr == "" {
		return fmt.Errorf("server.addr 涓嶈兘涓虹┖")
	}

	// 楠岃瘉鏁版嵁搴撻厤缃?
	if c.Database.Type == "" {
		return fmt.Errorf("database.type 涓嶈兘涓虹┖")
	}

	validDBTypes := map[string]bool{
		"sqlite":   true,
		"mysql":    true,
		"postgres": true,
	}
	if !validDBTypes[c.Database.Type] {
		return fmt.Errorf("涓嶆敮鎸佺殑鏁版嵁搴撶被鍨? %s", c.Database.Type)
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("database.dsn 涓嶈兘涓虹┖")
	}

	// 楠岃瘉鏃ュ織閰嶇疆
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
		"panic": true,
	}
	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("鏃犳晥鐨勬棩蹇楃骇鍒? %s", c.Logger.Level)
	}

	// 楠岃瘉鐢诲竷閰嶇疆
	if c.Canvas.ZoomMin <= 0 || c.Canvas.ZoomMax <= 0 || c.Canvas.ZoomMin >= c.Canvas.ZoomMax {
		return fmt.Errorf("鏃犳晥鐨勭缉鏀鹃厤缃? min=%f, max=%f", c.Canvas.ZoomMin, c.Canvas.ZoomMax)
	}

	return nil
}
