// Package services 实现业务逻辑层
//
// 设计参考：
// - DDD的应用服务模式
// - Kubernetes的控制器模式
// - 微服务的业务逻辑封装
//
// 特点：
// 1. 业务规则封装：包含所有业务逻辑
// 2. 事务管理：确保数据一致性
// 3. 事件发布：支持事件驱动架构
// 4. 验证和授权：统一的业务验证
package services

import (
	"context"
	"fmt"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// NodeService 节点服务接口
type NodeService interface {
	// 基础操作
	CreateNode(ctx context.Context, req CreateNodeRequest) (*domain.Node, error)
	GetNode(ctx context.Context, id domain.NodeID) (*domain.Node, error)
	UpdateNode(ctx context.Context, req UpdateNodeRequest) (*domain.Node, error)
	DeleteNode(ctx context.Context, id domain.NodeID) error

	// 批量操作
	CreateNodes(ctx context.Context, req CreateNodesRequest) ([]*domain.Node, error)
	GetNodes(ctx context.Context, req GetNodesRequest) (*GetNodesResponse, error)
	ListNodes(ctx context.Context) ([]*domain.Node, error)

	// 位置操作
	UpdateNodePosition(ctx context.Context, id domain.NodeID, position domain.Position) error
	MoveNodes(ctx context.Context, moves []NodeMove) error

	// 查询操作
	SearchNodes(ctx context.Context, req SearchNodesRequest) ([]*domain.Node, error)
	GetNodesInArea(ctx context.Context, req GetNodesInAreaRequest) ([]*domain.Node, error)
	GetNearbyNodes(ctx context.Context, req GetNearbyNodesRequest) ([]*domain.Node, error)

	// 分析操作
	GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error)
	GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error)
	ValidateNode(ctx context.Context, node *domain.Node) error
}

// 请求和响应结构体定义

// CreateNodeRequest 创建节点请求
type CreateNodeRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Type        domain.NodeType          `json:"type"`
	Position    domain.Position          `json:"position"`
	RobotCoords *domain.RobotCoordinates `json:"robot_coords,omitempty"`
	Properties  map[string]interface{}   `json:"properties,omitempty"`
	Style       domain.NodeStyle         `json:"style"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
}

// UpdateNodeRequest 更新节点请求
type UpdateNodeRequest struct {
	ID          domain.NodeID            `json:"id" binding:"required"`
	Name        *string                  `json:"name,omitempty"`
	Type        *domain.NodeType         `json:"type,omitempty"`
	Status      *domain.NodeStatus       `json:"status,omitempty"`
	Position    *domain.Position         `json:"position,omitempty"`
	RobotCoords *domain.RobotCoordinates `json:"robot_coords,omitempty"`
	Properties  map[string]interface{}   `json:"properties,omitempty"`
	Style       *domain.NodeStyle        `json:"style,omitempty"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
}

// CreateNodesRequest 批量创建节点请求
type CreateNodesRequest struct {
	Nodes []CreateNodeRequest `json:"nodes" binding:"required,dive"`
}

// GetNodesRequest 获取节点列表请求
type GetNodesRequest struct {
	Filter   repositories.NodeFilter `json:"filter"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	OrderBy  string                  `json:"order_by"`
	Order    string                  `json:"order"`
}

// GetNodesResponse 获取节点列表响应
type GetNodesResponse struct {
	Nodes      []*domain.Node `json:"nodes"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// SearchNodesRequest 搜索节点请求
type SearchNodesRequest struct {
	Query  string            `json:"query"`
	Type   domain.NodeType   `json:"type,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Limit  int               `json:"limit"`
}

// GetNodesInAreaRequest 获取区域内节点请�?
type GetNodesInAreaRequest struct {
	MinX float64 `json:"min_x" binding:"required"`
	MinY float64 `json:"min_y" binding:"required"`
	MaxX float64 `json:"max_x" binding:"required"`
	MaxY float64 `json:"max_y" binding:"required"`
}

// GetNearbyNodesRequest 获取附近节点请求
type GetNearbyNodesRequest struct {
	Position domain.Position `json:"position" binding:"required"`
	Radius   float64         `json:"radius" binding:"required,gt=0"`
	Limit    int             `json:"limit"`
}

// NodeMove 节点移动请求
type NodeMove struct {
	NodeID      domain.NodeID   `json:"node_id" binding:"required"`
	NewPosition domain.Position `json:"new_position" binding:"required"`
}

// nodeService 节点服务实现
type nodeService struct {
	nodeRepo repositories.NodeRepository
}

// NewNodeService 创建节点服务实例
func NewNodeService(nodeRepo repositories.NodeRepository) NodeService {
	return &nodeService{
		nodeRepo: nodeRepo,
	}
}

// CreateNode 创建节点
func (s *nodeService) CreateNode(ctx context.Context, req CreateNodeRequest) (*domain.Node, error) {
	// 1. 创建节点实体
	node := domain.NewNode(req.Name, req.Position)

	// 2. 设置可选属�?
	if req.Type != "" {
		node.Type = req.Type
	}

	if req.RobotCoords != nil {
		node.RobotCoords = req.RobotCoords
	}

	if req.Properties != nil {
		node.Properties = req.Properties
	}

	if req.Style.Shape != "" {
		node.Style = req.Style
	}

	if req.Labels != nil {
		node.Metadata.Labels = req.Labels
	}

	if req.Annotations != nil {
		node.Metadata.Annotations = req.Annotations
	}

	// 3. 业务规则验证
	if err := s.ValidateNode(ctx, node); err != nil {
		return nil, fmt.Errorf("节点验证失败: %w", err)
	}

	// 4. 持久�?
	if err := s.nodeRepo.Create(ctx, node); err != nil {
		return nil, fmt.Errorf("创建节点失败: %w", err)
	}

	return node, nil
}

// GetNode 获取节点
func (s *nodeService) GetNode(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取节点失败: %w", err)
	}

	return node, nil
}

// UpdateNode 更新节点
func (s *nodeService) UpdateNode(ctx context.Context, req UpdateNodeRequest) (*domain.Node, error) {
	// 1. 获取现有节点
	node, err := s.nodeRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("节点不存�? %w", err)
	}

	// 2. 应用更新
	if req.Name != nil {
		node.Name = *req.Name
	}

	if req.Type != nil {
		node.Type = *req.Type
	}

	if req.Status != nil {
		node.Status = *req.Status
	}

	if req.Position != nil {
		node.Position = *req.Position
	}

	if req.RobotCoords != nil {
		node.RobotCoords = req.RobotCoords
	}

	if req.Properties != nil {
		node.Properties = req.Properties
	}

	if req.Style != nil {
		node.Style = *req.Style
	}

	if req.Labels != nil {
		node.Metadata.Labels = req.Labels
	}

	if req.Annotations != nil {
		node.Metadata.Annotations = req.Annotations
	}

	// 3. 验证更新后的节点
	if err := s.ValidateNode(ctx, node); err != nil {
		return nil, fmt.Errorf("节点验证失败: %w", err)
	}

	// 4. 持久化更�?
	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return nil, fmt.Errorf("更新节点失败: %w", err)
	}

	return node, nil
}

// DeleteNode 删除节点
func (s *nodeService) DeleteNode(ctx context.Context, id domain.NodeID) error {
	// 1. 检查节点是否存�?
	_, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("节点不存�? %w", err)
	}

	// 2. 检查是否有路径连接 (业务规则: 不能删除有连接的节点)
	connectedNodes, err := s.nodeRepo.GetConnectedNodes(ctx, id)
	if err != nil {
		return fmt.Errorf("检查节点连接失败: %w", err)
	}

	if len(connectedNodes) > 0 {
		return fmt.Errorf("不能删除有路径连接的节点，请先删除相关路径")
	}

	// 3. 执行删除
	if err := s.nodeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除节点失败: %w", err)
	}

	return nil
}

// CreateNodes 批量创建节点
func (s *nodeService) CreateNodes(ctx context.Context, req CreateNodesRequest) ([]*domain.Node, error) {
	nodes := make([]*domain.Node, 0, len(req.Nodes))

	// 1. 创建所有节点实�?
	for _, nodeReq := range req.Nodes {
		node := domain.NewNode(nodeReq.Name, nodeReq.Position)

		// 设置属�?
		if nodeReq.Type != "" {
			node.Type = nodeReq.Type
		}
		if nodeReq.RobotCoords != nil {
			node.RobotCoords = nodeReq.RobotCoords
		}
		if nodeReq.Properties != nil {
			node.Properties = nodeReq.Properties
		}
		if nodeReq.Style.Shape != "" {
			node.Style = nodeReq.Style
		}
		if nodeReq.Labels != nil {
			node.Metadata.Labels = nodeReq.Labels
		}
		if nodeReq.Annotations != nil {
			node.Metadata.Annotations = nodeReq.Annotations
		}

		// 验证节点
		if err := s.ValidateNode(ctx, node); err != nil {
			return nil, fmt.Errorf("节点 %s 验证失败: %w", node.Name, err)
		}

		nodes = append(nodes, node)
	}

	// 2. 批量创建
	if err := s.nodeRepo.CreateBatch(ctx, nodes); err != nil {
		return nil, fmt.Errorf("批量创建节点失败: %w", err)
	}

	return nodes, nil
}

// GetNodes 获取节点列表
func (s *nodeService) GetNodes(ctx context.Context, req GetNodesRequest) (*GetNodesResponse, error) {
	// 1. 设置默认分页参数
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	// 2. 构建查询选项
	options := repositories.ListOptions{
		Filter:   req.Filter,
		Page:     req.Page,
		PageSize: req.PageSize,
		OrderBy:  req.OrderBy,
		Order:    req.Order,
	}

	// 3. 查询节点和总数
	nodes, err := s.nodeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("查询节点列表失败: %w", err)
	}

	total, err := s.nodeRepo.Count(ctx, req.Filter)
	if err != nil {
		return nil, fmt.Errorf("统计节点数量失败: %w", err)
	}

	// 4. 计算总页�?
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &GetNodesResponse{
		Nodes:      nodes,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListNodes 获取所有节点列�?
func (s *nodeService) ListNodes(ctx context.Context) ([]*domain.Node, error) {
	// 构建查�选项，不分页
	options := repositories.ListOptions{
		PageSize: 0, // 0 表示不分�?
	}

	nodes, err := s.nodeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %w", err)
	}

	return nodes, nil
}

// UpdateNodePosition 更新节点位置
func (s *nodeService) UpdateNodePosition(ctx context.Context, id domain.NodeID, position domain.Position) error {
	// 获取节点
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("节点不存�? %w", err)
	}

	// 更新位置
	node.Position = position

	// 保存更新
	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return fmt.Errorf("更新节点位置失败: %w", err)
	}

	return nil
}

// MoveNodes 批量移动节点
func (s *nodeService) MoveNodes(ctx context.Context, moves []NodeMove) error {
	// 获取所有需要移动的节点
	nodeIDs := make([]domain.NodeID, len(moves))
	for i, move := range moves {
		nodeIDs[i] = move.NodeID
	}

	nodes, err := s.nodeRepo.GetByIDs(ctx, nodeIDs)
	if err != nil {
		return fmt.Errorf("获取节点失败: %w", err)
	}

	// 创建节点ID到位置的映射
	nodePositions := make(map[domain.NodeID]domain.Position)
	for _, move := range moves {
		nodePositions[move.NodeID] = move.NewPosition
	}

	// 更新节点位置
	for _, node := range nodes {
		if newPos, exists := nodePositions[node.ID]; exists {
			node.Position = newPos
		}
	}

	// 批量更新
	if err := s.nodeRepo.UpdateBatch(ctx, nodes); err != nil {
		return fmt.Errorf("批量更新节点位置失败: %w", err)
	}

	return nil
}

// SearchNodes 搜索节点
func (s *nodeService) SearchNodes(ctx context.Context, req SearchNodesRequest) ([]*domain.Node, error) {
	filter := repositories.NodeFilter{
		Name:   req.Query,
		Type:   req.Type,
		Labels: req.Labels,
	}

	options := repositories.ListOptions{
		Filter:   filter,
		PageSize: req.Limit,
	}

	nodes, err := s.nodeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("搜索节点失败: %w", err)
	}

	return nodes, nil
}

// GetNodesInArea 获取区域内的节点
func (s *nodeService) GetNodesInArea(ctx context.Context, req GetNodesInAreaRequest) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetByArea(ctx, req.MinX, req.MinY, req.MaxX, req.MaxY)
	if err != nil {
		return nil, fmt.Errorf("获取区域内节点失�? %w", err)
	}

	return nodes, nil
}

// GetNearbyNodes 获取附近的节�?
func (s *nodeService) GetNearbyNodes(ctx context.Context, req GetNearbyNodesRequest) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetNearby(ctx, req.Position, req.Radius)
	if err != nil {
		return nil, fmt.Errorf("获取附近节点失败: %w", err)
	}

	// 如果设置了限制，截取结果
	if req.Limit > 0 && len(nodes) > req.Limit {
		nodes = nodes[:req.Limit]
	}

	return nodes, nil
}

// GetConnectedNodes 获取连接的节�?
func (s *nodeService) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetConnectedNodes(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("获取连接节点失败: %w", err)
	}

	return nodes, nil
}

// GetIsolatedNodes 获取孤立节点
func (s *nodeService) GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetIsolatedNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取孤立节点失败: %w", err)
	}

	return nodes, nil
}

// ValidateNode 验证节点
func (s *nodeService) ValidateNode(ctx context.Context, node *domain.Node) error {
	// 1. 基础验证
	if err := node.IsValid(); err != nil {
		return err
	}

	// 2. 业务规则验证

	// 名称不能重复（可选的业务规则�?
	// 这里可以添加更多业务验证逻辑

	// 3. 位置合理性验�?
	// 可以验证位置是否在允许的范围�?

	return nil
}
