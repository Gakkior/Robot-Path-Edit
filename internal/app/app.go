// Package app 是应用程序的主要组装�?
//
// 设计参考：
// - Uber FX的依赖注入模�?
// - Kubernetes Controller Manager的组件协�?
// - Grafana的应用程序架�?
//
// 职责�?
// 1. 组件初始化和依赖注入
// 2. 应用程序生命周期管理
// 3. 各个服务层的协调
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"robot-path-editor/internal/config"
	"robot-path-editor/internal/database"
	"robot-path-editor/internal/handlers"
	"robot-path-editor/internal/repositories"
	"robot-path-editor/internal/services"
	"robot-path-editor/pkg/middleware"
	"robot-path-editor/web"
)

// Application 应用程序主结�?
// 采用依赖注入模式，管理所有组件的生命周期
type Application struct {
	config *config.Config
	server *http.Server
	db     database.Database

	// 服务�?- 核心业务逻辑
	nodeService   services.NodeService
	pathService   services.PathService
	layoutService services.LayoutService
	dbService     services.DatabaseService

	// 处理器层 - HTTP API处理
	handlers *handlers.Handlers

	log *logrus.Entry
}

// New 创建新的应用程序实例
// 采用构造器模式，确保所有依赖正确初始化
func New(cfg *config.Config) (*Application, error) {
	log := logrus.WithField("component", "app")

	// 1. 初始化数据库和仓储层
	var nodeRepo repositories.NodeRepository
	var pathRepo repositories.PathRepository
	var dbConnRepo repositories.DatabaseConnectionRepository
	var tableMappingRepo repositories.TableMappingRepository
	var db database.Database

	// 尝试初始化数据库
	database, err := database.New(cfg.Database)
	if err != nil {
		// 如果数据库初始化失败，使用内存仓储作为后�?
		log.WithError(err).Warn("数据库初始化失败，使用内存存�?)

		// 使用内存仓储
		nodeRepo = repositories.NewMemoryNodeRepository()
		// 暂时使用nil，稍后实现其他内存仓�?
		pathRepo = nil
		dbConnRepo = nil
		tableMappingRepo = nil
		db = nil
	} else {
		// 使用数据库仓�?
		db = database
		nodeRepo = repositories.NewNodeRepository(db)
		pathRepo = repositories.NewPathRepository(db)
		dbConnRepo = repositories.NewDatabaseConnectionRepository(db)
		tableMappingRepo = repositories.NewTableMappingRepository(db)
	}

	// 2. 初始化服务层 - 业务逻辑
	nodeService := services.NewNodeService(nodeRepo)

	// 如果使用内存模式，创建简化的服务
	var pathService services.PathService
	var layoutService services.LayoutService
	var dbService services.DatabaseService

	if pathRepo != nil {
		pathService = services.NewPathService(pathRepo, nodeRepo)
		layoutService = services.NewLayoutService(nodeService, pathService)
		dbService = services.NewDatabaseService(dbConnRepo, tableMappingRepo)
	} else {
		// 内存模式下的简化服�?
		pathService = &services.MockPathService{}
		layoutService = &services.MockLayoutService{}
		dbService = &services.MockDatabaseService{}
	}

	// 4. 初始化处理器�?- API接口
	handlers := handlers.New(
		nodeService,
		pathService,
		layoutService,
		dbService,
	)

	// 5. 创建HTTP服务�?
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	app := &Application{
		config:        cfg,
		server:        server,
		db:            db,
		nodeService:   nodeService,
		pathService:   pathService,
		layoutService: layoutService,
		dbService:     dbService,
		handlers:      handlers,
		log:           log,
	}

	// 6. 配置路由
	if err := app.setupRoutes(); err != nil {
		return nil, fmt.Errorf("配置路由失败: %w", err)
	}

	log.Info("应用程序初始化完�?)
	return app, nil
}

// Start 启动应用程序
// 参考Kubernetes Controller的启动模�?
func (a *Application) Start(ctx context.Context) error {
	a.log.Info("启动应用程序...")

	// 启动后台服务
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.WithError(err).Error("HTTP服务器启动失�?)
		}
	}()

	// 等待上下文取�?
	<-ctx.Done()
	return nil
}

// Stop 停止应用程序
// 实现优雅关闭，参考Kubernetes的优雅终�?
func (a *Application) Stop(ctx context.Context) error {
	a.log.Info("停止应用程序...")

	// 1. 停止HTTP服务�?
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.WithError(err).Error("HTTP服务器关闭失�?)
		return err
	}

	// 2. 关闭数据库连�?
	if err := a.db.Close(); err != nil {
		a.log.WithError(err).Error("数据库关闭失�?)
		return err
	}

	a.log.Info("应用程序已停�?)
	return nil
}

// setupRoutes 配置HTTP路由
// 参考Kubernetes API Server的路由设�?
func (a *Application) setupRoutes() error {
	// 根据环境设置Gin模式
	if a.config.Logger.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 基础中间�?- 参考Kubernetes API Server的中间件�?
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// 健康检查端�?- 参考Kubernetes的健康检�?
	router.GET("/health", a.handlers.HealthCheck)
	router.GET("/ready", a.handlers.ReadinessCheck)

	// 指标端点 - 参考Prometheus的指标暴�?
	if a.config.Metrics.Enabled {
		router.GET(a.config.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// API路由�?- RESTful API设计
	api := router.Group("/api/v1")
	{
		// 节点管理API
		nodes := api.Group("/nodes")
		{
			nodes.GET("", a.handlers.ListNodes)
			nodes.POST("", a.handlers.CreateNode)
			nodes.GET("/:id", a.handlers.GetNode)
			nodes.PUT("/:id", a.handlers.UpdateNode)
			nodes.DELETE("/:id", a.handlers.DeleteNode)
			nodes.PUT("/:id/position", a.handlers.UpdateNodePosition)
		}

		// 路径管理API
		paths := api.Group("/paths")
		{
			paths.GET("", a.handlers.ListPaths)
			paths.POST("", a.handlers.CreatePath)
			paths.GET("/:id", a.handlers.GetPath)
			paths.PUT("/:id", a.handlers.UpdatePath)
			paths.DELETE("/:id", a.handlers.DeletePath)
			paths.POST("/generate", a.handlers.GeneratePaths)
		}

		// 布局管理API
		layouts := api.Group("/layouts")
		{
			layouts.POST("/arrange", a.handlers.ArrangeNodes)
			layouts.GET("/algorithms", a.handlers.ListLayoutAlgorithms)
		}

		// 数据库管理API
		databases := api.Group("/databases")
		{
			databases.GET("/connections", a.handlers.ListDatabaseConnections)
			databases.POST("/connections", a.handlers.CreateDatabaseConnection)
			databases.PUT("/connections/:id", a.handlers.UpdateDatabaseConnection)
			databases.DELETE("/connections/:id", a.handlers.DeleteDatabaseConnection)
			databases.POST("/connections/:id/test", a.handlers.TestDatabaseConnection)

			databases.GET("/connections/:id/tables", a.handlers.ListTables)
			databases.GET("/connections/:id/tables/:table/columns", a.handlers.ListColumns)

			databases.GET("/mappings", a.handlers.ListTableMappings)
			databases.POST("/mappings", a.handlers.CreateTableMapping)
			databases.PUT("/mappings/:id", a.handlers.UpdateTableMapping)
			databases.DELETE("/mappings/:id", a.handlers.DeleteTableMapping)
		}

		// 图分析API
		analysis := api.Group("/analysis")
		{
			analysis.POST("/shortest-path", a.handlers.FindShortestPath)
			analysis.GET("/connectivity", a.handlers.AnalyzeConnectivity)
			analysis.GET("/cycles", a.handlers.DetectCycles)
		}
	}

	// WebSocket API - 实时通信
	ws := router.Group("/ws")
	{
		ws.GET("/canvas", a.handlers.CanvasWebSocket)
	}

	// 静态文件服�?- 前端资源
	router.StaticFS("/static", http.FS(web.StaticFiles))
	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", web.IndexHTML)
	})

	a.server.Handler = router
	return nil
}
