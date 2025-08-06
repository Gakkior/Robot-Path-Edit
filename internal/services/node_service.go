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

// NodeService 节点业务服务接口
type NodeService interface {
	// 基础CRUD操作
	CreateNode(ctx context.Context, req CreateNodeRequest) (*domain.Node, error)
	GetNode(ctx context.Context, id domain.NodeID) (*domain.Node, error)
	UpdateNode(ctx context.Context, req UpdateNodeRequest) (*domain.Node, error)
	DeleteNode(ctx context.Context, id domain.NodeID) error

	// 批量操作
	BatchCreateNodes(ctx context.Context, req BatchCreateNodesRequest) ([]*domain.Node, error)
	BatchUpdateNodes(ctx context.Context, req BatchUpdateNodesRequest) ([]*domain.Node, error)
	BatchDeleteNodes(ctx context.Context, ids []domain.NodeID) error

	// 查询操作
	ListNodes(ctx context.Context) ([]*domain.Node, error)
	SearchNodes(ctx context.Context, req SearchNodesRequest) (*SearchNodesResponse, error)
	GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error)

	// 业务操作
	ValidateNodePosition(ctx context.Context, position domain.Position) error
	CalculateDistance(ctx context.Context, node1ID, node2ID domain.NodeID) (float64, error)
}

// CreateNodeRequest 创建节点请求
type CreateNodeRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Type        domain.NodeType          `json:"type"`
	Position    domain.Position          `json:"position"`
	RobotCoords *domain.RobotCoordinates `json:"robot_coords,omitempty"`
	Properties  map[string]interface{}   `json:"properties,omitempty"`
	Style       domain.NodeStyle         `json:"style"`
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
}

// BatchCreateNodesRequest 批量创建节点请求
type BatchCreateNodesRequest struct {
	Nodes []CreateNodeRequest `json:"nodes" binding:"required"`
}

// BatchUpdateNodesRequest 批量更新节点请求
type BatchUpdateNodesRequest struct {
	Nodes []UpdateNodeRequest `json:"nodes" binding:"required"`
}

// SearchNodesRequest 搜索节点请求
type SearchNodesRequest struct {
	Query    string            `json:"query"`
	Type     domain.NodeType   `json:"type,omitempty"`
	Status   domain.NodeStatus `json:"status,omitempty"`
	Page     int               `json:"page,omitempty"`
	PageSize int               `json:"page_size,omitempty"`
}

// SearchNodesResponse 搜索节点响应
type SearchNodesResponse struct {
	Nodes      []*domain.Node `json:"nodes"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// nodeService 节点服务实现
type nodeService struct {
	nodeRepo repositories.NodeRepository
	pathRepo repositories.PathRepository
}

// NewNodeService 创建新的节点服务实例
func NewNodeService(nodeRepo repositories.NodeRepository, pathRepo repositories.PathRepository) NodeService {
	return &nodeService{
		nodeRepo: nodeRepo,
		pathRepo: pathRepo,
	}
}

// CreateNode 创建节点
func (s *nodeService) CreateNode(ctx context.Context, req CreateNodeRequest) (*domain.Node, error) {
	// 1. 验证请求参数
	if req.Name == "" {
		return nil, fmt.Errorf("节点名称不能为空")
	}

	// 2. 验证位置信息
	if err := s.ValidateNodePosition(ctx, req.Position); err != nil {
		return nil, fmt.Errorf("位置验证失败: %w", err)
	}

	// 3. 创建节点实体
	node := domain.NewNode(req.Name, string(req.Type))
	node.Position = req.Position
	node.RobotCoords = req.RobotCoords
	node.Properties = req.Properties
	node.Style = req.Style

	// 4. 持久化
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
		return nil, fmt.Errorf("节点不存在: %w", err)
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
		if err := s.ValidateNodePosition(ctx, *req.Position); err != nil {
			return nil, fmt.Errorf("位置验证失败: %w", err)
		}
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

	// 3. 更新时间戳
	node.UpdatedAt()

	// 4. 持久化
	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return nil, fmt.Errorf("更新节点失败: %w", err)
	}

	return node, nil
}

// DeleteNode 删除节点
func (s *nodeService) DeleteNode(ctx context.Context, id domain.NodeID) error {
	// 1. 检查节点是否存在
	_, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("节点不存在: %w", err)
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

// BatchCreateNodes 批量创建节点
func (s *nodeService) BatchCreateNodes(ctx context.Context, req BatchCreateNodesRequest) ([]*domain.Node, error) {
	nodes := make([]*domain.Node, 0, len(req.Nodes))

	for _, nodeReq := range req.Nodes {
		node, err := s.CreateNode(ctx, nodeReq)
		if err != nil {
			return nil, fmt.Errorf("批量创建节点失败: %w", err)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// BatchUpdateNodes 批量更新节点
func (s *nodeService) BatchUpdateNodes(ctx context.Context, req BatchUpdateNodesRequest) ([]*domain.Node, error) {
	nodes := make([]*domain.Node, 0, len(req.Nodes))

	for _, nodeReq := range req.Nodes {
		node, err := s.UpdateNode(ctx, nodeReq)
		if err != nil {
			return nil, fmt.Errorf("批量更新节点失败: %w", err)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// BatchDeleteNodes 批量删除节点
func (s *nodeService) BatchDeleteNodes(ctx context.Context, ids []domain.NodeID) error {
	for _, id := range ids {
		if err := s.DeleteNode(ctx, id); err != nil {
			return fmt.Errorf("批量删除节点失败: %w", err)
		}
	}
	return nil
}

// ListNodes 获取节点列表
func (s *nodeService) ListNodes(ctx context.Context) ([]*domain.Node, error) {
	filter := repositories.NodeFilter{
		PageSize: 0, // 0 表示不分页
	}

	nodes, err := s.nodeRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %w", err)
	}

	return nodes, nil
}

// SearchNodes 搜索节点
func (s *nodeService) SearchNodes(ctx context.Context, req SearchNodesRequest) (*SearchNodesResponse, error) {
	// 构建过滤条件
	filter := repositories.NodeFilter{
		Name:     req.Query,
		Type:     req.Type,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 设置默认分页
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	// 获取节点列表
	nodes, err := s.nodeRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("搜索节点失败: %w", err)
	}

	// 获取总数
	total, err := s.nodeRepo.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("统计节点数量失败: %w", err)
	}

	// 计算总页数
	totalPages := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))

	return &SearchNodesResponse{
		Nodes:      nodes,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetConnectedNodes 获取连接的节点
func (s *nodeService) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetConnectedNodes(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("获取连接节点失败: %w", err)
	}
	return nodes, nil
}

// ValidateNodePosition 验证节点位置
func (s *nodeService) ValidateNodePosition(ctx context.Context, position domain.Position) error {
	// 基础位置验证
	if position.X < -10000 || position.X > 10000 {
		return fmt.Errorf("X坐标超出有效范围 (-10000, 10000)")
	}
	if position.Y < -10000 || position.Y > 10000 {
		return fmt.Errorf("Y坐标超出有效范围 (-10000, 10000)")
	}
	if position.Z < -10000 || position.Z > 10000 {
		return fmt.Errorf("Z坐标超出有效范围 (-10000, 10000)")
	}

	return nil
}

// CalculateDistance 计算两个节点之间的距离
func (s *nodeService) CalculateDistance(ctx context.Context, node1ID, node2ID domain.NodeID) (float64, error) {
	// 获取两个节点
	node1, err := s.nodeRepo.GetByID(ctx, node1ID)
	if err != nil {
		return 0, fmt.Errorf("获取节点1失败: %w", err)
	}

	node2, err := s.nodeRepo.GetByID(ctx, node2ID)
	if err != nil {
		return 0, fmt.Errorf("获取节点2失败: %w", err)
	}

	// 计算欧几里得距离
	distance := node1.Position.DistanceTo(node2.Position)
	return distance, nil
}
