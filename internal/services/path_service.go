// Package services 路径服务实现
package services

import (
	"context"
	"fmt"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// PathService 路径服务接口
type PathService interface {
	CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error)
	GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error)
	UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error)
	DeletePath(ctx context.Context, id domain.PathID) error
	GetPaths(ctx context.Context, req GetPathsRequest) (*GetPathsResponse, error)
	ListPaths(ctx context.Context) ([]*domain.Path, error)
	GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
}

// CreatePathRequest 创建路径请求
type CreatePathRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Type        domain.PathType        `json:"type"`
	StartNodeID domain.NodeID          `json:"start_node_id" binding:"required"`
	EndNodeID   domain.NodeID          `json:"end_node_id" binding:"required"`
	Direction   domain.PathDirection   `json:"direction"`
	Weight      float64                `json:"weight"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// UpdatePathRequest 更新路径请求
type UpdatePathRequest struct {
	ID         domain.PathID          `json:"id" binding:"required"`
	Name       *string                `json:"name,omitempty"`
	Type       *domain.PathType       `json:"type,omitempty"`
	Status     *domain.PathStatus     `json:"status,omitempty"`
	Direction  *domain.PathDirection  `json:"direction,omitempty"`
	Weight     *float64               `json:"weight,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// GetPathsRequest 获取路径列表请求
type GetPathsRequest struct {
	Filter   repositories.PathFilter `json:"filter"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	OrderBy  string                  `json:"order_by"`
	Order    string                  `json:"order"`
}

// GetPathsResponse 获取路径列表响应
type GetPathsResponse struct {
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

// NewPathService 创建路径服务实例
func NewPathService(pathRepo repositories.PathRepository, nodeRepo repositories.NodeRepository) PathService {
	return &pathService{
		pathRepo: pathRepo,
		nodeRepo: nodeRepo,
	}
}

// CreatePath 创建路径
func (s *pathService) CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error) {
	// 验证起始和结束节点存�?
	if _, err := s.nodeRepo.GetByID(ctx, req.StartNodeID); err != nil {
		return nil, fmt.Errorf("起始节点不存�? %w", err)
	}

	if _, err := s.nodeRepo.GetByID(ctx, req.EndNodeID); err != nil {
		return nil, fmt.Errorf("结束节点不存�? %w", err)
	}

	// 创建路径实体
	path := domain.NewPath(req.Name, req.StartNodeID, req.EndNodeID)

	if req.Type != "" {
		path.Type = req.Type
	}
	if req.Direction != "" {
		path.Direction = req.Direction
	}
	if req.Weight > 0 {
		path.Weight = req.Weight
	}
	if req.Properties != nil {
		path.Properties = req.Properties
	}

	// 持久�?
	if err := s.pathRepo.Create(ctx, path); err != nil {
		return nil, fmt.Errorf("创建路径失败: %w", err)
	}

	return path, nil
}

// GetPath 获取路径
func (s *pathService) GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	return s.pathRepo.GetByID(ctx, id)
}

// UpdatePath 更新路径
func (s *pathService) UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error) {
	path, err := s.pathRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("路径不存�? %w", err)
	}

	if req.Name != nil {
		path.Name = *req.Name
	}
	if req.Type != nil {
		path.Type = *req.Type
	}
	if req.Status != nil {
		path.Status = *req.Status
	}
	if req.Direction != nil {
		path.Direction = *req.Direction
	}
	if req.Weight != nil {
		path.Weight = *req.Weight
	}
	if req.Properties != nil {
		path.Properties = req.Properties
	}

	if err := s.pathRepo.Update(ctx, path); err != nil {
		return nil, fmt.Errorf("更新路径失败: %w", err)
	}

	return path, nil
}

// DeletePath 删除路径
func (s *pathService) DeletePath(ctx context.Context, id domain.PathID) error {
	return s.pathRepo.Delete(ctx, id)
}

// GetPaths 获取路径列表
func (s *pathService) GetPaths(ctx context.Context, req GetPathsRequest) (*GetPathsResponse, error) {
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	options := repositories.PathListOptions{
		Filter:   req.Filter,
		Page:     req.Page,
		PageSize: req.PageSize,
		OrderBy:  req.OrderBy,
		Order:    req.Order,
	}

	paths, err := s.pathRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("查询路径列表失败: %w", err)
	}

	total, err := s.pathRepo.Count(ctx, req.Filter)
	if err != nil {
		return nil, fmt.Errorf("统计路径数量失败: %w", err)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &GetPathsResponse{
		Paths:      paths,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListPaths 获取所有路径列�?
func (s *pathService) ListPaths(ctx context.Context) ([]*domain.Path, error) {
	// 构建查�选项，不分页
	options := repositories.PathListOptions{
		PageSize: 0, // 0 表示不分�?
	}

	paths, err := s.pathRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("获取跾�列表失败: %w", err)
	}

	return paths, nil
}

// GetPathsByNode 获取节点相关的路�?
func (s *pathService) GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	return s.pathRepo.GetByNode(ctx, nodeID)
}
