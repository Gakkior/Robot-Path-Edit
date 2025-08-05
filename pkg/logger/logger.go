// Package logger 提供统一的日志管理功�?
//
// 设计参考：
// - Kubernetes的结构化日志系统
// - Grafana的日志管�?
// - Prometheus的日志规�?
//
// 特点�?
// 1. 结构化日志：支持JSON和文本格�?
// 2. 日志轮转：支持文件大小和时间轮转
// 3. 上下文感知：支持链路追踪
package logger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"robot-path-editor/internal/config"
)

// Init 初始化日志系�?
// 参考Kubernetes的日志初始化流程
func Init(cfg config.LoggerConfig) {
	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.WithError(err).Warn("解析日志级别失败，使用默认级�?info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	switch strings.ToLower(cfg.Format) {
	case "json":
		// JSON格式 - 适合生产环境和日志收集系�?
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	default:
		// 文本格式 - 适合开发环�?
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 设置输出目标
	switch strings.ToLower(cfg.Output) {
	case "file":
		if cfg.File == "" {
			logrus.Warn("日志文件路径为空，使用默认路�?)
			cfg.File = "./logs/app.log"
		}

		// 确保日志目录存在
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.WithError(err).Error("创建日志目录失败")
			return
		}

		// 打开日志文件
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.WithError(err).Error("打开日志文件失败")
			return
		}

		logrus.SetOutput(file)
	default:
		// 默认输出到标准输�?
		logrus.SetOutput(os.Stdout)
	}

	// 设置调用者信息报�?
	logrus.SetReportCaller(true)

	logrus.WithFields(logrus.Fields{
		"level":  cfg.Level,
		"format": cfg.Format,
		"output": cfg.Output,
	}).Info("日志系统初始化完�?)
}
