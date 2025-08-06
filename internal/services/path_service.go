// Package services 路径业务服务实现
//
// 设计参考：
// - 图算法的路径管理
// - 导航系统的路径规划
// - 网络拓扑的连接管理
//
// 特点：
// 1. 路径数据管理：CRUD操作
// 2. 图关系维护：节点连接管理
// 3. 路径验证：业务规则检查
// 4. 性能优化：批量操作支持
package services

import (
	"context"
	"fmt"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// PathService 路径业务服务接口
type PathService interface {
	// 基础CRUD操作
	CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error)
	GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error)
	UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error)
	DeletePath(ctx context.Context, id domain.PathID) error

	// 批量操作
	CreatePaths(ctx context.Context, req CreatePathsRequest) ([]*domain.Path, error)
	DeletePaths(ctx context.Context, ids []domain.PathID) error

	// 查询操作
	ListPaths(ctx context.Context, req ListPathsRequest) (*ListPathsResponse, error)
	GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
	GetPathsBetweenNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error)

	// 业务操作
	ValidatePath(ctx context.Context, path *domain.Path) error
	CalculatePathWeight(ctx context.Context, startNodeID, endNodeID domain.NodeID) (float64, error)
}

// CreatePathRequest 创建路径请求
type CreatePathRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Type        domain.PathType        `json:"type"`
	StartNodeID domain.NodeID          `json:"start_node_id" binding:"required"`
	EndNodeID   domain.NodeID          `json:"end_node_id" binding:"required"`
	Weight      float64                `json:"weight"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Style       domain.PathStyle       `json:"style"`
}

// UpdatePathRequest 更新路径请求
type UpdatePathRequest struct {
	ID         domain.PathID          `json:"id" binding:"required"`
	Name       *string                `json:"name,omitempty"`
	Type       *domain.PathType       `json:"type,omitempty"`
	Status     *domain.PathStatus     `json:"status,omitempty"`
	Weight     *float64               `json:"weight,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Style      *domain.PathStyle      `json:"style,omitempty"`
}

// CreatePathsRequest 批量创建路径请求
type CreatePathsRequest struct {
	Paths []CreatePathRequest `json:"paths" binding:"required"`
}

// ListPathsRequest 路径列表请求
type ListPathsRequest struct {
	StartNodeID domain.NodeID     `json:"start_node_id,omitempty"`
	EndNodeID   domain.NodeID     `json:"end_node_id,omitempty"`
	Type        domain.PathType   `json:"type,omitempty"`
	Status      domain.PathStatus `json:"status,omitempty"`
	Page        int               `json:"page,omitempty"`
	PageSize    int               `json:"page_size,omitempty"`
}

// ListPathsResponse 路径列表响应
type ListPathsResponse struct {
	Paths      []*domain.Path `json:"paths"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// pathService 路径服务实现
type pathService struct {
	pathRepo repositories.PathRepository
	nodeRepo repositories.NodeRepository
}

// NewPathService 创建新的路径服务实例
func NewPathService(pathRepo repositories.PathRepository, nodeRepo repositories.NodeRepository) PathService {
	return &pathService{
		pathRepo: pathRepo,
		nodeRepo: nodeRepo,
	}
}

// CreatePath 创建路径
func (s *pathService) CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error) {
	// 验证起始和结束节点存在
	if _, err := s.nodeRepo.GetByID(ctx, req.StartNodeID); err != nil {
		return nil, fmt.Errorf("起始节点不存在: %w", err)
	}

	if _, err := s.nodeRepo.GetByID(ctx, req.EndNodeID); err != nil {
		return nil, fmt.Errorf("结束节点不存在: %w", err)
	}

	// 创建路径实体
	path := domain.NewPath(req.Name, req.StartNodeID, req.EndNodeID)
	path.Type = req.Type
	path.Weight = req.Weight
	path.Properties = req.Properties
	path.Style = req.Style

	// 验证路径
	if err := s.ValidatePath(ctx, path); err != nil {
		return nil, fmt.Errorf("路径验证失败: %w", err)
	}

	// 持久化
	if err := s.pathRepo.Create(ctx, path); err != nil {
		return nil, fmt.Errorf("创建路径失败: %w", err)
	}

	return path, nil
}

// GetPath 获取路径
func (s *pathService) GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	path, err := s.pathRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %w", err)
	}
	return path, nil
}

// UpdatePath 更新路径
func (s *pathService) UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error) {
	// 获取现有路径
	path, err := s.pathRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %w", err)
	}

	// 应用更新
	if req.Name != nil {
		path.Name = *req.Name
	}
	if req.Type != nil {
		path.Type = *req.Type
	}
	if req.Status != nil {
		path.Status = *req.Status
	}
	if req.Weight != nil {
		path.Weight = *req.Weight
	}
	if req.Properties != nil {
		path.Properties = req.Properties
	}
	if req.Style != nil {
		path.Style = *req.Style
	}

	// 验证更新后的路径
	if err := s.ValidatePath(ctx, path); err != nil {
		return nil, fmt.Errorf("路径验证失败: %w", err)
	}

	// 更新时间戳
	path.UpdatedAt()

	// 持久化
	if err := s.pathRepo.Update(ctx, path); err != nil {
		return nil, fmt.Errorf("更新路径失败: %w", err)
	}

	return path, nil
}

// DeletePath 删除路径
func (s *pathService) DeletePath(ctx context.Context, id domain.PathID) error {
	// 检查路径是否存在
	_, err := s.pathRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("路径不存在: %w", err)
	}

	// 执行删除
	if err := s.pathRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除路径失败: %w", err)
	}

	return nil
}

// CreatePaths 批量创建路径
func (s *pathService) CreatePaths(ctx context.Context, req CreatePathsRequest) ([]*domain.Path, error) {
	paths := make([]*domain.Path, 0, len(req.Paths))

	for _, pathReq := range req.Paths {
		path, err := s.CreatePath(ctx, pathReq)
		if err != nil {
			return nil, fmt.Errorf("批量创建路径失败: %w", err)
		}
		paths = append(paths, path)
	}

	return paths, nil
}

// DeletePaths 批量删除路径
func (s *pathService) DeletePaths(ctx context.Context, ids []domain.PathID) error {
	return s.pathRepo.DeleteBatch(ctx, ids)
}

// ListPaths 获取路径列表
func (s *pathService) ListPaths(ctx context.Context, req ListPathsRequest) (*ListPathsResponse, error) {
	// 构建查询选项
	filter := repositories.PathFilter{
		StartNodeID: req.StartNodeID,
		EndNodeID:   req.EndNodeID,
		Type:        req.Type,
		Status:      req.Status,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	// 设置默认分页
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	// 获取路径列表
	paths, err := s.pathRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取路径列表失败: %w", err)
	}

	// 获取总数
	total, err := s.pathRepo.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("统计路径数量失败: %w", err)
	}

	// 计算总页数
	totalPages := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))

	return &ListPathsResponse{
		Paths:      paths,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetPathsByNode 获取节点相关的路径
func (s *pathService) GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	paths, err := s.pathRepo.GetByNode(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("获取节点路径失败: %w", err)
	}
	return paths, nil
}

// GetPathsBetweenNodes 获取两个节点之间的路径
func (s *pathService) GetPathsBetweenNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error) {
	paths, err := s.pathRepo.GetByNodes(ctx, startNodeID, endNodeID)
	if err != nil {
		return nil, fmt.Errorf("获取节点间路径失败: %w", err)
	}
	return paths, nil
}

// ValidatePath 验证路径
func (s *pathService) ValidatePath(ctx context.Context, path *domain.Path) error {
	// 使用域对象的验证方法
	if err := path.IsValid(); err != nil {
		return err
	}

	// 额外的业务验证
	if path.StartNodeID == path.EndNodeID {
		return fmt.Errorf("起始节点和结束节点不能相同")
	}

	// 检查权重范围
	if path.Weight < 0 {
		return fmt.Errorf("路径权重不能为负数")
	}

	if path.Weight > 10000 {
		return fmt.Errorf("路径权重不能超过10000")
	}

	return nil
}

// CalculatePathWeight 计算路径权重（基于节点距离）
func (s *pathService) CalculatePathWeight(ctx context.Context, startNodeID, endNodeID domain.NodeID) (float64, error) {
	// 获取起始节点
	startNode, err := s.nodeRepo.GetByID(ctx, startNodeID)
	if err != nil {
		return 0, fmt.Errorf("获取起始节点失败: %w", err)
	}

	// 获取结束节点
	endNode, err := s.nodeRepo.GetByID(ctx, endNodeID)
	if err != nil {
		return 0, fmt.Errorf("获取结束节点失败: %w", err)
	}

	// 计算欧几里得距离作为权重
	distance := startNode.Position.DistanceTo(endNode.Position)
	return distance, nil
}
