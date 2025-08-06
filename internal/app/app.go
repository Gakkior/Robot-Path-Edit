// Package app 是应用程序的主要组装器
//
// 设计参考：
// - Uber FX的依赖注入模式
// - Kubernetes Controller Manager的组件协调
// - Grafana的应用程序架构
//
// 职责：
// 1. 组件初始化和依赖注入
// 2. 应用程序生命周期管理
// 3. 各个服务层的协调
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"robot-path-editor/internal/config"
	"robot-path-editor/internal/database"
	"robot-path-editor/internal/handlers"
	"robot-path-editor/internal/repositories"
	"robot-path-editor/internal/services"
	"robot-path-editor/pkg/logger"
	"robot-path-editor/pkg/middleware"
	"robot-path-editor/web"
)

// Application 应用程序主结构
// 采用依赖注入模式，管理所有组件的生命周期
type Application struct {
	config *config.Config
	log    *logrus.Entry

	// 服务器 - 核心业务逻辑
	server *http.Server
	router *gin.Engine

	// 数据层
	db               database.Database
	nodeRepo         repositories.NodeRepository
	pathRepo         repositories.PathRepository
	dbConnRepo       repositories.DatabaseConnectionRepository
	tableMappingRepo repositories.TableMappingRepository

	// 业务服务层
	nodeService     services.NodeService
	pathService     services.PathService
	layoutService   services.LayoutService
	pluginService   services.PluginService
	databaseService services.DatabaseService

	// HTTP处理器
	handlers *handlers.Handlers
}

// New 创建新的应用程序实例
// 采用构造器模式，确保所有依赖正确初始化
func New(cfg *config.Config) (*Application, error) {
	log := logrus.WithField("component", "app")

	app := &Application{
		config: cfg,
		log:    log,
	}

	// 1. 初始化日志系统
	logger.Init(cfg.Logger)

	// 2. 初始化数据库和仓储层
	var nodeRepo repositories.NodeRepository
	var pathRepo repositories.PathRepository
	var dbConnRepo repositories.DatabaseConnectionRepository
	var tableMappingRepo repositories.TableMappingRepository
	var templateRepo repositories.TemplateRepository
	var db database.Database

	// 尝试初始化数据库
	database, err := database.New(cfg.Database)
	if err != nil {
		// 如果数据库初始化失败，使用内存仓储作为后备
		log.WithError(err).Warn("数据库初始化失败，使用内存存储")

		// 使用内存仓储
		nodeRepo = repositories.NewMemoryNodeRepository()
		// 暂时使用nil，稍后实现其他内存仓储
		pathRepo = nil
		dbConnRepo = nil
		tableMappingRepo = nil
		templateRepo = nil
		db = nil
	} else {
		// 使用数据库仓储
		nodeRepo = repositories.NewNodeRepository(database)
		pathRepo = repositories.NewPathRepository(database)
		dbConnRepo = repositories.NewDatabaseConnectionRepository(database)
		tableMappingRepo = repositories.NewTableMappingRepository(database)
		templateRepo = repositories.NewTemplateRepository(database)
		db = database
	}

	// 3. 初始化业务服务层
	var nodeService services.NodeService
	var pathService services.PathService
	var layoutService services.LayoutService
	var pluginService services.PluginService
	var databaseService services.DatabaseService
	var dataSyncService services.DataSyncService
	var templateService services.TemplateService

	if pathRepo == nil {
		// 如果使用内存模式，创建简化的服务
		nodeService = services.NewNodeService(nodeRepo, nil)
		pathService = &services.MockPathService{}
		layoutService = services.NewLayoutService()
		pluginService = services.NewPluginService()
		databaseService = &services.MockDatabaseService{}
		dataSyncService = &services.MockDataSyncService{}
		templateService = &services.MockTemplateService{}
	} else {
		// 内存模式下的简化服务
		nodeService = services.NewNodeService(nodeRepo, pathRepo)
		pathService = services.NewPathService(pathRepo, nodeRepo)
		layoutService = services.NewLayoutService()
		pluginService = services.NewPluginService()
		databaseService = services.NewDatabaseService(dbConnRepo, tableMappingRepo)
		dataSyncService = services.NewDataSyncService(dbConnRepo, tableMappingRepo, nodeRepo, pathRepo)
		templateService = services.NewTemplateService(templateRepo, nodeRepo, pathRepo)
	}

	// 4. 初始化处理器层 - API接口
	handlers := handlers.New(
		nodeService,
		pathService,
		layoutService,
		databaseService,
		dataSyncService,
		templateService,
	)

	// 5. 创建HTTP服务器
	router := gin.New()
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 设置应用程序字段
	app.db = db
	app.nodeRepo = nodeRepo
	app.pathRepo = pathRepo
	app.dbConnRepo = dbConnRepo
	app.tableMappingRepo = tableMappingRepo
	app.nodeService = nodeService
	app.pathService = pathService
	app.layoutService = layoutService
	app.pluginService = pluginService
	app.databaseService = databaseService
	app.handlers = handlers
	app.router = router
	app.server = server

	// 6. 配置路由
	if err := app.setupRoutes(); err != nil {
		return nil, fmt.Errorf("配置路由失败: %w", err)
	}

	log.Info("应用程序初始化完成")
	return app, nil
}

// Start 启动应用程序
// 参考Kubernetes Controller的启动模式
func (a *Application) Start(ctx context.Context) error {
	a.log.Info("启动应用程序...")

	// 启动后台服务
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.WithError(err).Error("HTTP服务器启动失败")
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	return nil
}

// Stop 停止应用程序
// 实现优雅关闭，参考Kubernetes的优雅终止
func (a *Application) Stop(ctx context.Context) error {
	a.log.Info("停止应用程序...")

	// 1. 停止HTTP服务器
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.WithError(err).Error("HTTP服务器关闭失败")
		return err
	}

	// 2. 关闭数据库连接
	if err := a.db.Close(); err != nil {
		a.log.WithError(err).Error("数据库关闭失败")
		return err
	}

	a.log.Info("应用程序已停止")
	return nil
}

// setupRoutes 配置HTTP路由
// 参考Kubernetes API Server的路由设计
func (a *Application) setupRoutes() error {
	// 根据环境设置Gin模式
	if a.config.Logger.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 基础中间件 - 参考Kubernetes API Server的中间件设计
	a.router.Use(middleware.Logger())
	a.router.Use(middleware.Recovery())
	a.router.Use(middleware.CORS())

	// 健康检查端点 - 参考Kubernetes的健康检查
	a.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// 指标端点 - 参考Prometheus的指标暴露
	a.router.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, gin.H{"metrics": "placeholder"})
	})

	// API路由组 - RESTful API设计
	api := a.router.Group("/api/v1")
	{
		// 节点管理
		nodes := api.Group("/nodes")
		{
			nodes.GET("", a.handlers.ListNodes)
			nodes.POST("", a.handlers.CreateNode)
			nodes.GET("/:id", a.handlers.GetNode)
			nodes.PUT("/:id", a.handlers.UpdateNode)
			nodes.DELETE("/:id", a.handlers.DeleteNode)
			nodes.POST("/batch", a.handlers.BatchCreateNodes)
			nodes.PUT("/batch", a.handlers.BatchUpdateNodes)
			nodes.DELETE("/batch", a.handlers.BatchDeleteNodes)
			nodes.GET("/search", a.handlers.SearchNodes)
			nodes.GET("/:id/connected", a.handlers.GetConnectedNodes)
		}

		// 路径管理
		paths := api.Group("/paths")
		{
			paths.GET("", a.handlers.ListPaths)
			paths.POST("", a.handlers.CreatePath)
			paths.GET("/:id", a.handlers.GetPath)
			paths.PUT("/:id", a.handlers.UpdatePath)
			paths.DELETE("/:id", a.handlers.DeletePath)
			paths.GET("/node/:nodeId", a.handlers.GetPathsByNode)
		}

		// 布局算法
		layout := api.Group("/layout")
		{
			layout.POST("/force-directed", a.handlers.ApplyForceDirectedLayout)
			layout.POST("/hierarchical", a.handlers.ApplyHierarchicalLayout)
			layout.POST("/circular", a.handlers.ApplyCircularLayout)
			layout.POST("/grid", a.handlers.ApplyGridLayout)
		}

		// 路径生成算法
		generation := api.Group("/generation")
		{
			generation.POST("/shortest-paths", a.handlers.GenerateShortestPaths)
			generation.POST("/full-connectivity", a.handlers.GenerateFullConnectivity)
			generation.POST("/tree-structure", a.handlers.GenerateTreeStructure)
			generation.POST("/nearest-neighbor", a.handlers.GenerateNearestNeighborPaths)
			generation.POST("/grid-paths", a.handlers.GenerateGridPaths)
		}

		// 数据库连接管理
		db := api.Group("/database")
		{
			db.GET("/connections", a.handlers.ListDatabaseConnections)
			db.POST("/connections", a.handlers.CreateDatabaseConnection)
			db.PUT("/connections/:id", a.handlers.UpdateDatabaseConnection)
			db.DELETE("/connections/:id", a.handlers.DeleteDatabaseConnection)
			db.POST("/connections/:id/test", a.handlers.TestDatabaseConnection)
		}

		// 表映射管理
		mapping := api.Group("/mapping")
		{
			mapping.GET("/tables", a.handlers.ListTableMappings)
			mapping.POST("/tables", a.handlers.CreateTableMapping)
			mapping.PUT("/tables/:id", a.handlers.UpdateTableMapping)
			mapping.DELETE("/tables/:id", a.handlers.DeleteTableMapping)
		}

		// 分析相关处理器
		analysis := api.Group("/analysis")
		{
			analysis.POST("/shortest-path", a.handlers.FindShortestPath)
			analysis.GET("/connectivity", a.handlers.AnalyzeConnectivity)
			analysis.GET("/cycles", a.handlers.DetectCycles)
		}

		// 数据同步相关处理器
		sync := api.Group("/sync")
		{
			sync.POST("/mappings/:mappingId/nodes", a.handlers.SyncNodesFromExternal)
			sync.POST("/mappings/:mappingId/paths", a.handlers.SyncPathsFromExternal)
			sync.POST("/mappings/:mappingId/all", a.handlers.SyncAllDataFromExternal)
			sync.GET("/validate-table", a.handlers.ValidateExternalTable)
		}

		// 模板相关处理器
		templates := api.Group("/templates")
		{
			templates.GET("", a.handlers.ListTemplates)
			templates.POST("", a.handlers.CreateTemplate)
			templates.GET("/public", a.handlers.GetPublicTemplates)
			templates.GET("/search", a.handlers.SearchTemplates)
			templates.GET("/category/:category", a.handlers.GetTemplatesByCategory)
			templates.GET("/stats", a.handlers.GetTemplateStats)
			templates.GET("/:id", a.handlers.GetTemplate)
			templates.PUT("/:id", a.handlers.UpdateTemplate)
			templates.DELETE("/:id", a.handlers.DeleteTemplate)
			templates.POST("/:id/apply", a.handlers.ApplyTemplate)
			templates.POST("/:id/clone", a.handlers.CloneTemplate)
			templates.GET("/:id/export", a.handlers.ExportTemplate)
			templates.POST("/import", a.handlers.ImportTemplate)
			templates.POST("/save-as", a.handlers.SaveAsTemplate)
		}
	}

	// WebSocket端点
	a.router.GET("/ws/canvas", a.handlers.CanvasWebSocket)

	// 静态文件服务 - 前端资源
	a.router.StaticFS("/static", http.FS(web.StaticFiles))

	// 新前端静态文件 (如果构建了的话)
	a.router.Static("/app/new", "./web/static/new-frontend")

	// 首页 - 介绍页面
	a.router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", web.IndexHTML)
	})

	// 应用界面 - 主要的编辑器界面
	a.router.GET("/app", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", web.AppHTML)
	})

	// 新前端应用界面
	a.router.GET("/app/new", func(c *gin.Context) {
		c.File("./web/static/new-frontend/index.html")
	})

	return nil
}
