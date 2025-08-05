// Package services 璺緞鏈嶅姟瀹炵幇
package services

import (
	"context"
	"fmt"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// PathService 璺緞鏈嶅姟鎺ュ彛
type PathService interface {
	CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error)
	GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error)
	UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error)
	DeletePath(ctx context.Context, id domain.PathID) error
	GetPaths(ctx context.Context, req GetPathsRequest) (*GetPathsResponse, error)
	ListPaths(ctx context.Context) ([]*domain.Path, error)
	GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error)
}

// CreatePathRequest 鍒涘缓璺緞璇锋眰
type CreatePathRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Type        domain.PathType        `json:"type"`
	StartNodeID domain.NodeID          `json:"start_node_id" binding:"required"`
	EndNodeID   domain.NodeID          `json:"end_node_id" binding:"required"`
	Direction   domain.PathDirection   `json:"direction"`
	Weight      float64                `json:"weight"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// UpdatePathRequest 鏇存柊璺緞璇锋眰
type UpdatePathRequest struct {
	ID         domain.PathID          `json:"id" binding:"required"`
	Name       *string                `json:"name,omitempty"`
	Type       *domain.PathType       `json:"type,omitempty"`
	Status     *domain.PathStatus     `json:"status,omitempty"`
	Direction  *domain.PathDirection  `json:"direction,omitempty"`
	Weight     *float64               `json:"weight,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// GetPathsRequest 鑾峰彇璺緞鍒楄〃璇锋眰
type GetPathsRequest struct {
	Filter   repositories.PathFilter `json:"filter"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	OrderBy  string                  `json:"order_by"`
	Order    string                  `json:"order"`
}

// GetPathsResponse 鑾峰彇璺緞鍒楄〃鍝嶅簲
type GetPathsResponse struct {
	Paths      []*domain.Path `json:"paths"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// pathService 璺緞鏈嶅姟瀹炵幇
type pathService struct {
	pathRepo repositories.PathRepository
	nodeRepo repositories.NodeRepository
}

// NewPathService 鍒涘缓璺緞鏈嶅姟瀹炰緥
func NewPathService(pathRepo repositories.PathRepository, nodeRepo repositories.NodeRepository) PathService {
	return &pathService{
		pathRepo: pathRepo,
		nodeRepo: nodeRepo,
	}
}

// CreatePath 鍒涘缓璺緞
func (s *pathService) CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error) {
	// 楠岃瘉璧峰鍜岀粨鏉熻妭鐐瑰瓨鍦?
	if _, err := s.nodeRepo.GetByID(ctx, req.StartNodeID); err != nil {
		return nil, fmt.Errorf("璧峰鑺傜偣涓嶅瓨鍦? %w", err)
	}

	if _, err := s.nodeRepo.GetByID(ctx, req.EndNodeID); err != nil {
		return nil, fmt.Errorf("缁撴潫鑺傜偣涓嶅瓨鍦? %w", err)
	}

	// 鍒涘缓璺緞瀹炰綋
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

	// 鎸佷箙鍖?
	if err := s.pathRepo.Create(ctx, path); err != nil {
		return nil, fmt.Errorf("鍒涘缓璺緞澶辫触: %w", err)
	}

	return path, nil
}

// GetPath 鑾峰彇璺緞
func (s *pathService) GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	return s.pathRepo.GetByID(ctx, id)
}

// UpdatePath 鏇存柊璺緞
func (s *pathService) UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error) {
	path, err := s.pathRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("璺緞涓嶅瓨鍦? %w", err)
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
		return nil, fmt.Errorf("鏇存柊璺緞澶辫触: %w", err)
	}

	return path, nil
}

// DeletePath 鍒犻櫎璺緞
func (s *pathService) DeletePath(ctx context.Context, id domain.PathID) error {
	return s.pathRepo.Delete(ctx, id)
}

// GetPaths 鑾峰彇璺緞鍒楄〃
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
		return nil, fmt.Errorf("鏌ヨ璺緞鍒楄〃澶辫触: %w", err)
	}

	total, err := s.pathRepo.Count(ctx, req.Filter)
	if err != nil {
		return nil, fmt.Errorf("缁熻璺緞鏁伴噺澶辫触: %w", err)
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

// ListPaths 鑾峰彇鎵€鏈夎矾寰勫垪琛?
func (s *pathService) ListPaths(ctx context.Context) ([]*domain.Path, error) {
	// 鏋勫缓鏌ヨ閫夐」锛屼笉鍒嗛〉
	options := repositories.PathListOptions{
		PageSize: 0, // 0 琛ㄧず涓嶅垎椤?
	}

	paths, err := s.pathRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇璺緞鍒楄〃澶辫触: %w", err)
	}

	return paths, nil
}

// GetPathsByNode 鑾峰彇鑺傜偣鐩稿叧鐨勮矾寰?
func (s *pathService) GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	return s.pathRepo.GetByNode(ctx, nodeID)
}
