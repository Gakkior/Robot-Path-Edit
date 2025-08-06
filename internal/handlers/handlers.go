// Package handlers HTTP请求处理器
//
// 设计参考：
// - RESTful API设计原则
// - Kubernetes API Server的处理器模式
// - Gin框架的最佳实践
//
// 特点：
// 1. 统一的响应格式
// 2. 错误处理中间件
// 3. 请求参数验证
// 4. 业务逻辑委托
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/services"
)

// Handlers HTTP处理器集合
type Handlers struct {
	nodeService     services.NodeService
	pathService     services.PathService
	layoutService   services.LayoutService
	databaseService services.DatabaseService
	dataSyncService services.DataSyncService
	templateService services.TemplateService
}

// New 创建新的处理器实例
func New(
	nodeService services.NodeService,
	pathService services.PathService,
	layoutService services.LayoutService,
	databaseService services.DatabaseService,
	dataSyncService services.DataSyncService,
	templateService services.TemplateService,
) *Handlers {
	return &Handlers{
		nodeService:     nodeService,
		pathService:     pathService,
		layoutService:   layoutService,
		databaseService: databaseService,
		dataSyncService: dataSyncService,
		templateService: templateService,
	}
}

// 节点相关处理器
func (h *Handlers) ListNodes(c *gin.Context) {
	nodes, err := h.nodeService.ListNodes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

func (h *Handlers) CreateNode(c *gin.Context) {
	var req services.CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := h.nodeService.CreateNode(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"node": node})
}

func (h *Handlers) GetNode(c *gin.Context) {
	id := c.Param("id")
	node, err := h.nodeService.GetNode(c.Request.Context(), domain.NodeID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}

func (h *Handlers) UpdateNode(c *gin.Context) {
	var req services.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = domain.NodeID(c.Param("id"))
	node, err := h.nodeService.UpdateNode(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node})
}

func (h *Handlers) DeleteNode(c *gin.Context) {
	id := c.Param("id")
	err := h.nodeService.DeleteNode(c.Request.Context(), domain.NodeID(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "节点删除成功"})
}

func (h *Handlers) BatchCreateNodes(c *gin.Context) {
	var req services.BatchCreateNodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodes, err := h.nodeService.BatchCreateNodes(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"nodes": nodes})
}

func (h *Handlers) BatchUpdateNodes(c *gin.Context) {
	var req services.BatchUpdateNodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodes, err := h.nodeService.BatchUpdateNodes(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

func (h *Handlers) BatchDeleteNodes(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodeIDs := make([]domain.NodeID, len(req.IDs))
	for i, id := range req.IDs {
		nodeIDs[i] = domain.NodeID(id)
	}
	err := h.nodeService.BatchDeleteNodes(c.Request.Context(), nodeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量删除节点成功"})
}

func (h *Handlers) SearchNodes(c *gin.Context) {
	var req services.SearchNodesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.nodeService.SearchNodes(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetConnectedNodes(c *gin.Context) {
	nodeID := c.Param("id")
	nodes, err := h.nodeService.GetConnectedNodes(c.Request.Context(), domain.NodeID(nodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

// 路径相关处理器
func (h *Handlers) ListPaths(c *gin.Context) {
	var req services.ListPathsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.pathService.ListPaths(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) CreatePath(c *gin.Context) {
	var req services.CreatePathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path, err := h.pathService.CreatePath(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"path": path})
}

func (h *Handlers) GetPath(c *gin.Context) {
	id := c.Param("id")
	path, err := h.pathService.GetPath(c.Request.Context(), domain.PathID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"path": path})
}

func (h *Handlers) UpdatePath(c *gin.Context) {
	var req services.UpdatePathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = domain.PathID(c.Param("id"))
	path, err := h.pathService.UpdatePath(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": path})
}

func (h *Handlers) DeletePath(c *gin.Context) {
	id := c.Param("id")
	err := h.pathService.DeletePath(c.Request.Context(), domain.PathID(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "路径删除成功"})
}

func (h *Handlers) GetPathsByNode(c *gin.Context) {
	nodeID := c.Param("nodeId")
	paths, err := h.pathService.GetPathsByNode(c.Request.Context(), domain.NodeID(nodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"paths": paths})
}

// 布局相关处理器
func (h *Handlers) ApplyForceDirectedLayout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "应用力导向布局"})
}

func (h *Handlers) ApplyHierarchicalLayout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "应用层次化布局"})
}

func (h *Handlers) ApplyCircularLayout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "应用圆形布局"})
}

func (h *Handlers) ApplyGridLayout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "应用网格布局"})
}

// 路径生成相关处理器
func (h *Handlers) GenerateShortestPaths(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "生成最短路径"})
}

func (h *Handlers) GenerateFullConnectivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "生成完全连通图"})
}

func (h *Handlers) GenerateTreeStructure(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "生成树状结构"})
}

func (h *Handlers) GenerateNearestNeighborPaths(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "生成最近邻路径"})
}

func (h *Handlers) GenerateGridPaths(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "生成网格路径"})
}

// 数据库相关处理器
func (h *Handlers) ListDatabaseConnections(c *gin.Context) {
	connections, err := h.databaseService.ListConnections(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"connections": connections})
}

func (h *Handlers) CreateDatabaseConnection(c *gin.Context) {
	var req services.CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := h.databaseService.CreateConnection(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"connection": conn})
}

func (h *Handlers) UpdateDatabaseConnection(c *gin.Context) {
	var req services.UpdateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = c.Param("id")
	conn, err := h.databaseService.UpdateConnection(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connection": conn})
}

func (h *Handlers) DeleteDatabaseConnection(c *gin.Context) {
	id := c.Param("id")
	err := h.databaseService.DeleteConnection(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "数据库连接删除成功"})
}

func (h *Handlers) TestDatabaseConnection(c *gin.Context) {
	id := c.Param("id")
	err := h.databaseService.TestConnection(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "数据库连接测试成功"})
}

func (h *Handlers) ListTableMappings(c *gin.Context) {
	mappings, err := h.databaseService.ListTableMappings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mappings": mappings})
}

func (h *Handlers) CreateTableMapping(c *gin.Context) {
	var req services.CreateTableMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mapping, err := h.databaseService.CreateTableMapping(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mapping": mapping})
}

func (h *Handlers) UpdateTableMapping(c *gin.Context) {
	var req services.UpdateTableMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = c.Param("id")
	mapping, err := h.databaseService.UpdateTableMapping(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mapping": mapping})
}

func (h *Handlers) DeleteTableMapping(c *gin.Context) {
	id := c.Param("id")
	err := h.databaseService.DeleteTableMapping(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "表映射删除成功"})
}

// 分析相关处理器
func (h *Handlers) FindShortestPath(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "查找最短路径"})
}

func (h *Handlers) AnalyzeConnectivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "分析连通性"})
}

func (h *Handlers) DetectCycles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "检测环路"})
}

// WebSocket处理器
func (h *Handlers) CanvasWebSocket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "画布WebSocket"})
}

// 数据同步相关处理器
func (h *Handlers) SyncNodesFromExternal(c *gin.Context) {
	mappingID := c.Param("mappingId")
	result, err := h.dataSyncService.SyncNodesFromExternal(c.Request.Context(), mappingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *Handlers) SyncPathsFromExternal(c *gin.Context) {
	mappingID := c.Param("mappingId")
	result, err := h.dataSyncService.SyncPathsFromExternal(c.Request.Context(), mappingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *Handlers) SyncAllDataFromExternal(c *gin.Context) {
	mappingID := c.Param("mappingId")
	result, err := h.dataSyncService.SyncAllDataFromExternal(c.Request.Context(), mappingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *Handlers) ValidateExternalTable(c *gin.Context) {
	connectionID := c.Query("connection_id")
	tableName := c.Query("table_name")

	if connectionID == "" || tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供connection_id和table_name参数"})
		return
	}

	result, err := h.dataSyncService.ValidateExternalTable(c.Request.Context(), connectionID, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}
