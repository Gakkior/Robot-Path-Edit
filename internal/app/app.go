// Package app 鏄簲鐢ㄧ▼搴忕殑涓昏缁勮灞?
//
// 璁捐鍙傝€冿細
// - Uber FX鐨勪緷璧栨敞鍏ユā寮?
// - Kubernetes Controller Manager鐨勭粍浠跺崗璋?
// - Grafana鐨勫簲鐢ㄧ▼搴忔灦鏋?
//
// 鑱岃矗锛?
// 1. 缁勪欢鍒濆鍖栧拰渚濊禆娉ㄥ叆
// 2. 搴旂敤绋嬪簭鐢熷懡鍛ㄦ湡绠＄悊
// 3. 鍚勪釜鏈嶅姟灞傜殑鍗忚皟
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

// Application 搴旂敤绋嬪簭涓荤粨鏋?
// 閲囩敤渚濊禆娉ㄥ叆妯″紡锛岀鐞嗘墍鏈夌粍浠剁殑鐢熷懡鍛ㄦ湡
type Application struct {
	config *config.Config
	server *http.Server
	db     database.Database

	// 鏈嶅姟灞?- 鏍稿績涓氬姟閫昏緫
	nodeService   services.NodeService
	pathService   services.PathService
	layoutService services.LayoutService
	dbService     services.DatabaseService

	// 澶勭悊鍣ㄥ眰 - HTTP API澶勭悊
	handlers *handlers.Handlers

	log *logrus.Entry
}

// New 鍒涘缓鏂扮殑搴旂敤绋嬪簭瀹炰緥
// 閲囩敤鏋勯€犲櫒妯″紡锛岀‘淇濇墍鏈変緷璧栨纭垵濮嬪寲
func New(cfg *config.Config) (*Application, error) {
	log := logrus.WithField("component", "app")

	// 1. 鍒濆鍖栨暟鎹簱鍜屼粨鍌ㄥ眰
	var nodeRepo repositories.NodeRepository
	var pathRepo repositories.PathRepository
	var dbConnRepo repositories.DatabaseConnectionRepository
	var tableMappingRepo repositories.TableMappingRepository
	var db database.Database

	// 灏濊瘯鍒濆鍖栨暟鎹簱
	database, err := database.New(cfg.Database)
	if err != nil {
		// 濡傛灉鏁版嵁搴撳垵濮嬪寲澶辫触锛屼娇鐢ㄥ唴瀛樹粨鍌ㄤ綔涓哄悗澶?
		log.WithError(err).Warn("鏁版嵁搴撳垵濮嬪寲澶辫触锛屼娇鐢ㄥ唴瀛樺瓨鍌?)

		// 浣跨敤鍐呭瓨浠撳偍
		nodeRepo = repositories.NewMemoryNodeRepository()
		// 鏆傛椂浣跨敤nil锛岀◢鍚庡疄鐜板叾浠栧唴瀛樹粨鍌?
		pathRepo = nil
		dbConnRepo = nil
		tableMappingRepo = nil
		db = nil
	} else {
		// 浣跨敤鏁版嵁搴撲粨鍌?
		db = database
		nodeRepo = repositories.NewNodeRepository(db)
		pathRepo = repositories.NewPathRepository(db)
		dbConnRepo = repositories.NewDatabaseConnectionRepository(db)
		tableMappingRepo = repositories.NewTableMappingRepository(db)
	}

	// 2. 鍒濆鍖栨湇鍔″眰 - 涓氬姟閫昏緫
	nodeService := services.NewNodeService(nodeRepo)

	// 濡傛灉浣跨敤鍐呭瓨妯″紡锛屽垱寤虹畝鍖栫殑鏈嶅姟
	var pathService services.PathService
	var layoutService services.LayoutService
	var dbService services.DatabaseService

	if pathRepo != nil {
		pathService = services.NewPathService(pathRepo, nodeRepo)
		layoutService = services.NewLayoutService(nodeService, pathService)
		dbService = services.NewDatabaseService(dbConnRepo, tableMappingRepo)
	} else {
		// 鍐呭瓨妯″紡涓嬬殑绠€鍖栨湇鍔?
		pathService = &services.MockPathService{}
		layoutService = &services.MockLayoutService{}
		dbService = &services.MockDatabaseService{}
	}

	// 4. 鍒濆鍖栧鐞嗗櫒灞?- API鎺ュ彛
	handlers := handlers.New(
		nodeService,
		pathService,
		layoutService,
		dbService,
	)

	// 5. 鍒涘缓HTTP鏈嶅姟鍣?
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

	// 6. 閰嶇疆璺敱
	if err := app.setupRoutes(); err != nil {
		return nil, fmt.Errorf("閰嶇疆璺敱澶辫触: %w", err)
	}

	log.Info("搴旂敤绋嬪簭鍒濆鍖栧畬鎴?)
	return app, nil
}

// Start 鍚姩搴旂敤绋嬪簭
// 鍙傝€僈ubernetes Controller鐨勫惎鍔ㄦā寮?
func (a *Application) Start(ctx context.Context) error {
	a.log.Info("鍚姩搴旂敤绋嬪簭...")

	// 鍚姩鍚庡彴鏈嶅姟
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.WithError(err).Error("HTTP鏈嶅姟鍣ㄥ惎鍔ㄥけ璐?)
		}
	}()

	// 绛夊緟涓婁笅鏂囧彇娑?
	<-ctx.Done()
	return nil
}

// Stop 鍋滄搴旂敤绋嬪簭
// 瀹炵幇浼橀泤鍏抽棴锛屽弬鑰僈ubernetes鐨勪紭闆呯粓姝?
func (a *Application) Stop(ctx context.Context) error {
	a.log.Info("鍋滄搴旂敤绋嬪簭...")

	// 1. 鍋滄HTTP鏈嶅姟鍣?
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.WithError(err).Error("HTTP鏈嶅姟鍣ㄥ叧闂け璐?)
		return err
	}

	// 2. 鍏抽棴鏁版嵁搴撹繛鎺?
	if err := a.db.Close(); err != nil {
		a.log.WithError(err).Error("鏁版嵁搴撳叧闂け璐?)
		return err
	}

	a.log.Info("搴旂敤绋嬪簭宸插仠姝?)
	return nil
}

// setupRoutes 閰嶇疆HTTP璺敱
// 鍙傝€僈ubernetes API Server鐨勮矾鐢辫璁?
func (a *Application) setupRoutes() error {
	// 鏍规嵁鐜璁剧疆Gin妯″紡
	if a.config.Logger.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 鍩虹涓棿浠?- 鍙傝€僈ubernetes API Server鐨勪腑闂翠欢閾?
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// 鍋ュ悍妫€鏌ョ鐐?- 鍙傝€僈ubernetes鐨勫仴搴锋鏌?
	router.GET("/health", a.handlers.HealthCheck)
	router.GET("/ready", a.handlers.ReadinessCheck)

	// 鎸囨爣绔偣 - 鍙傝€働rometheus鐨勬寚鏍囨毚闇?
	if a.config.Metrics.Enabled {
		router.GET(a.config.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// API璺敱缁?- RESTful API璁捐
	api := router.Group("/api/v1")
	{
		// 鑺傜偣绠＄悊API
		nodes := api.Group("/nodes")
		{
			nodes.GET("", a.handlers.ListNodes)
			nodes.POST("", a.handlers.CreateNode)
			nodes.GET("/:id", a.handlers.GetNode)
			nodes.PUT("/:id", a.handlers.UpdateNode)
			nodes.DELETE("/:id", a.handlers.DeleteNode)
			nodes.PUT("/:id/position", a.handlers.UpdateNodePosition)
		}

		// 璺緞绠＄悊API
		paths := api.Group("/paths")
		{
			paths.GET("", a.handlers.ListPaths)
			paths.POST("", a.handlers.CreatePath)
			paths.GET("/:id", a.handlers.GetPath)
			paths.PUT("/:id", a.handlers.UpdatePath)
			paths.DELETE("/:id", a.handlers.DeletePath)
			paths.POST("/generate", a.handlers.GeneratePaths)
		}

		// 甯冨眬绠＄悊API
		layouts := api.Group("/layouts")
		{
			layouts.POST("/arrange", a.handlers.ArrangeNodes)
			layouts.GET("/algorithms", a.handlers.ListLayoutAlgorithms)
		}

		// 鏁版嵁搴撶鐞咥PI
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

		// 鍥惧垎鏋怉PI
		analysis := api.Group("/analysis")
		{
			analysis.POST("/shortest-path", a.handlers.FindShortestPath)
			analysis.GET("/connectivity", a.handlers.AnalyzeConnectivity)
			analysis.GET("/cycles", a.handlers.DetectCycles)
		}
	}

	// WebSocket API - 瀹炴椂閫氫俊
	ws := router.Group("/ws")
	{
		ws.GET("/canvas", a.handlers.CanvasWebSocket)
	}

	// 闈欐€佹枃浠舵湇鍔?- 鍓嶇璧勬簮
	router.StaticFS("/static", http.FS(web.StaticFiles))
	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", web.IndexHTML)
	})

	a.server.Handler = router
	return nil
}
