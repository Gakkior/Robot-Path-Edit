// Package main 是应用程序的入口点
//
// 设计参考：
// - Cobra CLI框架的最佳实践
// - Grafana Server的启动流程
// - Kubernetes kubectl的命令结构
//
// 特点：
// 1. 命令行参数支持
// 2. 配置文件加载
// 3. 优雅关闭
// 4. 信号处理
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"robot-path-editor/internal/app"
	"robot-path-editor/internal/config"
	"robot-path-editor/pkg/logger"
)

var (
	// 版本信息
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"

	// 配置文件路径
	configPath string
)

func main() {
	// 创建根命令 - 参考kubectl的命令结构
	rootCmd := &cobra.Command{
		Use:   "robot-path-editor",
		Short: "机器人路径编辑器服务端",
		Long:  `一个现代化的机器人路径编辑器，支持可视化编辑和数据库管理。`,
		RunE:  runServer,
	}

	// 添加命令行参数 - 参考Grafana的配置选项
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./configs", "配置文件路径")
	rootCmd.PersistentFlags().String("host", "0.0.0.0", "服务监听地址")
	rootCmd.PersistentFlags().Int("port", 8080, "服务监听端口")
	rootCmd.PersistentFlags().String("db-type", "sqlite", "数据库类型(sqlite, mysql)")
	rootCmd.PersistentFlags().String("db-dsn", "data.db", "数据库连接字符串")
	rootCmd.PersistentFlags().String("log-level", "info", "日志级别")

	// 绑定参数到viper - 参考Kubernetes的配置管理
	viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("database.type", rootCmd.PersistentFlags().Lookup("db-type"))
	viper.BindPFlag("database.dsn", rootCmd.PersistentFlags().Lookup("db-dsn"))
	viper.BindPFlag("logger.level", rootCmd.PersistentFlags().Lookup("log-level"))

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "启动失败: %v\n", err)
		os.Exit(1)
	}
}

// runServer 启动服务器
// 参考Grafana Server的启动流程
func runServer(cmd *cobra.Command, args []string) error {
	// 1. 加载配置 - 参考Viper的配置加载流程
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 2. 初始化日志 - 参考Kubernetes的结构化日志
	logger.Init(cfg.Logger)

	log := logrus.WithFields(logrus.Fields{
		"component": "server",
		"version":   Version,
	})

	log.Info("机器人路径编辑器启动中...")

	// 3. 创建应用实例 - 参考DDD的应用服务模式
	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("创建应用失败: %w", err)
	}

	// 4. 设置信号处理 - 参考Docker的优雅停机
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.WithField("signal", sig).Info("接收到停止信号，开始优雅关闭...")
		cancel()
	}()

	// 5. 启动服务 - 参考HTTP服务器的启动模式
	log.WithFields(logrus.Fields{
		"host": cfg.Server.Host,
		"port": cfg.Server.Port,
	}).Info("服务器启动完成")

	if err := application.Start(ctx); err != nil {
		log.WithError(err).Error("服务器运行错误")
		return err
	}

	log.Info("服务器已优雅关闭")
	return nil
}
