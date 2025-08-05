// Package services Mock服务实现
// 用于演示模式，提供基本功�?
package services

import (
	"context"

	"robot-path-editor/internal/domain"
)

// MockPathService 模拟路径服务
type MockPathService struct{}

func (s *MockPathService) CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error) {
	path := domain.NewPath(req.Name, req.StartNodeID, req.EndNodeID)
	return path, nil
}

func (s *MockPathService) GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	return &domain.Path{ID: id, Name: "模拟路径"}, nil
}

func (s *MockPathService) UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error) {
	return &domain.Path{ID: req.ID, Name: "更新的模拟路�?}, nil
}

func (s *MockPathService) DeletePath(ctx context.Context, id domain.PathID) error {
	return nil
}

func (s *MockPathService) GetPaths(ctx context.Context, req GetPathsRequest) (*GetPathsResponse, error) {
	return &GetPathsResponse{
		Paths:      []*domain.Path{},
		Total:      0,
		Page:       1,
		PageSize:   20,
		TotalPages: 0,
	}, nil
}

func (s *MockPathService) ListPaths(ctx context.Context) ([]*domain.Path, error) {
	return []*domain.Path{}, nil
}

func (s *MockPathService) GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	return []*domain.Path{}, nil
}

// MockLayoutService 模拟布局服务
type MockLayoutService struct{}

func (s *MockLayoutService) ArrangeNodes(ctx context.Context, algorithm string) (map[string]domain.Position, error) {
	return make(map[string]domain.Position), nil
}

func (s *MockLayoutService) ApplyGridLayout(nodes []domain.Node, spacing float64) []domain.Node {
	return nodes
}

func (s *MockLayoutService) ApplyForceDirectedLayout(nodes []domain.Node, paths []domain.Path, iterations int) []domain.Node {
	return nodes
}

func (s *MockLayoutService) ApplyCircularLayout(nodes []domain.Node, radius, centerX, centerY float64) []domain.Node {
	return nodes
}

// MockDatabaseService 模拟数据库服�?
type MockDatabaseService struct{}

func (s *MockDatabaseService) CreateDatabaseConnection(ctx context.Context, conn *domain.DatabaseConnection) error {
	return nil
}

func (s *MockDatabaseService) GetDatabaseConnections(ctx context.Context) ([]*domain.DatabaseConnection, error) {
	return []*domain.DatabaseConnection{}, nil
}

func (s *MockDatabaseService) CreateTableMapping(ctx context.Context, mapping *domain.TableMapping) error {
	return nil
}

func (s *MockDatabaseService) GetTableMappings(ctx context.Context) ([]*domain.TableMapping, error) {
	return []*domain.TableMapping{}, nil
}
