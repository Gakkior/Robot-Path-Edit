// Package main 是应用程序的入口�?
//
// 架构设计参考：
// - Kubernetes API Server的启动模�?
// - Grafana Server的配置管�?
// - Docker的命令行接口设计
//
// 设计理念�?
// 1. 单一职责：main函数只负责启动，不包含业务逻辑
// 2. 配置驱动：通过配置文件和环境变量控制行�?
// 3. 优雅启停：支持信号处理和资源清理
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"robot-path-editor/internal/app"
	"robot-path-editor/internal/config"
	"robot-path-editor/pkg/logger"
)

var (
	// 版本信息，构建时注入
	version   = "dev"
	buildTime = "unknown"
	gitHash   = "unknown"
)

func main() {
	// 创建根命�?- 参考kubectl的命令结�?
	rootCmd := &cobra.Command{
		Use:     "robot-path-editor",
		Short:   "机器人路径编辑器 - 通用的点位和路径管理工具",
		Long:    `一个现代化的三端兼容机器人路径编辑器，支持可视化编辑和数据库管理。`,
		Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitHash),
		RunE:    runServer,
	}

	// 添加命令行参�?- 参考Grafana的配置选项
	rootCmd.PersistentFlags().String("config", "", "配置文件路径")
	rootCmd.PersistentFlags().String("log-level", "info", "日志级别 (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("addr", ":8080", "服务监听地址")
	rootCmd.PersistentFlags().String("db-type", "sqlite", "数据库类�?(sqlite, mysql)")
	rootCmd.PersistentFlags().String("db-dsn", "./data/app.db", "数据库连接字符串")

	// 绑定参数到viper - 参考Kubernetes的配置管�?
	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("应用启动失败")
	}
}

// runServer 启动服务�?
// 参考Grafana Server的启动流�?
func runServer(cmd *cobra.Command, args []string) error {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 2. 初始化日�?- 参考Kubernetes的结构化日志
	logger.Init(cfg.Logger)

	log := logrus.WithFields(logrus.Fields{
		"version":    version,
		"build_time": buildTime,
		"git_hash":   gitHash,
	})

	log.Info("机器人路径编辑器启动�?..")

	// 3. 创建应用实例 - 依赖注入模式，参考Uber FX
	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("创建应用实例失败: %w", err)
	}

	// 4. 设置信号处理 - 参考Docker的优雅停�?
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 5. 启动应用
	go func() {
		if err := application.Start(ctx); err != nil {
			log.WithError(err).Error("应用启动失败")
			cancel()
		}
	}()

	log.WithField("addr", cfg.Server.Addr).Info("服务器启动成�?)

	// 6. 等待停止信号
	select {
	case sig := <-sigChan:
		log.WithField("signal", sig).Info("收到停止信号，开始优雅关�?..")
	case <-ctx.Done():
		log.Info("应用上下文被取消")
	}

	// 7. 优雅关闭 - 参考Kubernetes的优雅终�?
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := application.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("应用关闭时发生错�?)
		return err
	}

	log.Info("应用已成功关�?)
	return nil
}
