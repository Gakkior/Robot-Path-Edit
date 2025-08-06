// Package services Mock服务实现
// 用于演示模式，提供基本功能
package services

import (
	"context"

	"robot-path-editor/internal/domain"
)

// MockPathService Mock路径服务实现
type MockPathService struct{}

// CreatePath 创建路径（Mock实现）
func (s *MockPathService) CreatePath(ctx context.Context, req CreatePathRequest) (*domain.Path, error) {
	return domain.NewPath(req.Name, req.StartNodeID, req.EndNodeID), nil
}

// GetPath 获取路径（Mock实现）
func (s *MockPathService) GetPath(ctx context.Context, id domain.PathID) (*domain.Path, error) {
	return &domain.Path{ID: id, Name: "模拟路径"}, nil
}

// UpdatePath 更新路径（Mock实现）
func (s *MockPathService) UpdatePath(ctx context.Context, req UpdatePathRequest) (*domain.Path, error) {
	return &domain.Path{ID: req.ID, Name: "更新的模拟路径"}, nil
}

// DeletePath 删除路径（Mock实现）
func (s *MockPathService) DeletePath(ctx context.Context, id domain.PathID) error {
	return nil
}

// CreatePaths 批量创建路径（Mock实现）
func (s *MockPathService) CreatePaths(ctx context.Context, req CreatePathsRequest) ([]*domain.Path, error) {
	paths := make([]*domain.Path, len(req.Paths))
	for i, pathReq := range req.Paths {
		paths[i] = domain.NewPath(pathReq.Name, pathReq.StartNodeID, pathReq.EndNodeID)
	}
	return paths, nil
}

// DeletePaths 批量删除路径（Mock实现）
func (s *MockPathService) DeletePaths(ctx context.Context, ids []domain.PathID) error {
	return nil
}

// ListPaths 列出路径（Mock实现）
func (s *MockPathService) ListPaths(ctx context.Context, req ListPathsRequest) (*ListPathsResponse, error) {
	return &ListPathsResponse{
		Paths:      []*domain.Path{},
		Total:      0,
		Page:       1,
		PageSize:   20,
		TotalPages: 0,
	}, nil
}

// GetPathsByNode 获取节点相关路径（Mock实现）
func (s *MockPathService) GetPathsByNode(ctx context.Context, nodeID domain.NodeID) ([]*domain.Path, error) {
	return []*domain.Path{}, nil
}

// GetPathsBetweenNodes 获取节点间路径（Mock实现）
func (s *MockPathService) GetPathsBetweenNodes(ctx context.Context, startNodeID, endNodeID domain.NodeID) ([]*domain.Path, error) {
	return []*domain.Path{}, nil
}

// ValidatePath 验证路径（Mock实现）
func (s *MockPathService) ValidatePath(ctx context.Context, path *domain.Path) error {
	return nil
}

// CalculatePathWeight 计算路径权重（Mock实现）
func (s *MockPathService) CalculatePathWeight(ctx context.Context, startNodeID, endNodeID domain.NodeID) (float64, error) {
	return 1.0, nil
}

// MockDatabaseService Mock数据库服务实现
type MockDatabaseService struct{}

// CreateConnection 创建数据库连接（Mock实现）
func (s *MockDatabaseService) CreateConnection(ctx context.Context, req CreateConnectionRequest) (*domain.DatabaseConnection, error) {
	return &domain.DatabaseConnection{
		ID:   "mock-connection",
		Name: req.Name,
		Type: req.Type,
		Host: req.Host,
		Port: req.Port,
	}, nil
}

// GetConnection 获取数据库连接（Mock实现）
func (s *MockDatabaseService) GetConnection(ctx context.Context, id string) (*domain.DatabaseConnection, error) {
	return &domain.DatabaseConnection{ID: id, Name: "Mock连接"}, nil
}

// UpdateConnection 更新数据库连接（Mock实现）
func (s *MockDatabaseService) UpdateConnection(ctx context.Context, req UpdateConnectionRequest) (*domain.DatabaseConnection, error) {
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	return &domain.DatabaseConnection{ID: string(req.ID), Name: name}, nil
}

// DeleteConnection 删除数据库连接（Mock实现）
func (s *MockDatabaseService) DeleteConnection(ctx context.Context, id string) error {
	return nil
}

// ListConnections 列出数据库连接（Mock实现）
func (s *MockDatabaseService) ListConnections(ctx context.Context) ([]*domain.DatabaseConnection, error) {
	return []*domain.DatabaseConnection{}, nil
}

// TestConnection 测试数据库连接（Mock实现）
func (s *MockDatabaseService) TestConnection(ctx context.Context, id string) error {
	return nil
}

// CreateTableMapping 创建表映射配置（Mock实现）
func (s *MockDatabaseService) CreateTableMapping(ctx context.Context, req CreateTableMappingRequest) (*domain.TableMapping, error) {
	return &domain.TableMapping{
		ID:           "mock-mapping",
		ConnectionID: req.ConnectionID,
		TableName:    req.TableName,
	}, nil
}

// GetTableMapping 获取表映射配置（Mock实现）
func (s *MockDatabaseService) GetTableMapping(ctx context.Context, id string) (*domain.TableMapping, error) {
	return &domain.TableMapping{ID: id, TableName: "mock_table"}, nil
}

// UpdateTableMapping 更新表映射配置（Mock实现）
func (s *MockDatabaseService) UpdateTableMapping(ctx context.Context, req UpdateTableMappingRequest) (*domain.TableMapping, error) {
	tableName := ""
	if req.TableName != nil {
		tableName = *req.TableName
	}
	return &domain.TableMapping{ID: string(req.ID), TableName: tableName}, nil
}

// DeleteTableMapping 删除表映射配置（Mock实现）
func (s *MockDatabaseService) DeleteTableMapping(ctx context.Context, id string) error {
	return nil
}

// ListTableMappings 列出表映射配置（Mock实现）
func (s *MockDatabaseService) ListTableMappings(ctx context.Context) ([]*domain.TableMapping, error) {
	return []*domain.TableMapping{}, nil
}
