// Package services Mock服务实现
// 用于演示模式，提供基本功能
package services

import (
	"context"
	"fmt"

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
	return &domain.DatabaseConnection{ID: req.ID, Name: name}, nil
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
	return &domain.TableMapping{ID: req.ID, TableName: tableName}, nil
}

// DeleteTableMapping 删除表映射配置（Mock实现）
func (s *MockDatabaseService) DeleteTableMapping(ctx context.Context, id string) error {
	return nil
}

// ListTableMappings 列出表映射配置（Mock实现）
func (s *MockDatabaseService) ListTableMappings(ctx context.Context) ([]*domain.TableMapping, error) {
	return []*domain.TableMapping{}, nil
}

// MockDataSyncService Mock数据同步服务实现
type MockDataSyncService struct{}

// SyncNodesFromExternal 从外部数据库同步节点数据（Mock实现）
func (s *MockDataSyncService) SyncNodesFromExternal(ctx context.Context, mappingID string) (*SyncResult, error) {
	return &SyncResult{
		NodesCreated: 0,
		NodesUpdated: 0,
		Errors:       []string{"内存模式下不支持外部数据同步"},
	}, nil
}

// SyncPathsFromExternal 从外部数据库同步路径数据（Mock实现）
func (s *MockDataSyncService) SyncPathsFromExternal(ctx context.Context, mappingID string) (*SyncResult, error) {
	return &SyncResult{
		PathsCreated: 0,
		PathsUpdated: 0,
		Errors:       []string{"内存模式下不支持外部数据同步"},
	}, nil
}

// SyncAllDataFromExternal 全量同步数据（Mock实现）
func (s *MockDataSyncService) SyncAllDataFromExternal(ctx context.Context, mappingID string) (*SyncResult, error) {
	return &SyncResult{
		NodesCreated: 0,
		NodesUpdated: 0,
		PathsCreated: 0,
		PathsUpdated: 0,
		Errors:       []string{"内存模式下不支持外部数据同步"},
	}, nil
}

// ValidateExternalTable 验证外部数据库表结构（Mock实现）
func (s *MockDataSyncService) ValidateExternalTable(ctx context.Context, connectionID, tableName string) (*TableValidationResult, error) {
	return &TableValidationResult{
		Valid:   false,
		Columns: []string{},
		Message: "内存模式下不支持外部表验证",
	}, nil
}

// MockTemplateService Mock模板服务实现
type MockTemplateService struct{}

// CreateTemplate 创建模板（Mock实现）
func (s *MockTemplateService) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*domain.Template, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// GetTemplate 获取模板（Mock实现）
func (s *MockTemplateService) GetTemplate(ctx context.Context, id string) (*domain.Template, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// UpdateTemplate 更新模板（Mock实现）
func (s *MockTemplateService) UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*domain.Template, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// DeleteTemplate 删除模板（Mock实现）
func (s *MockTemplateService) DeleteTemplate(ctx context.Context, id string) error {
	return fmt.Errorf("内存模式下不支持模板功能")
}

// ListTemplates 列出模板（Mock实现）
func (s *MockTemplateService) ListTemplates(ctx context.Context, req ListTemplatesRequest) (*ListTemplatesResponse, error) {
	return &ListTemplatesResponse{
		Templates:  []*domain.Template{},
		Total:      0,
		Page:       1,
		PageSize:   20,
		TotalPages: 0,
	}, nil
}

// SearchTemplates 搜索模板（Mock实现）
func (s *MockTemplateService) SearchTemplates(ctx context.Context, query string) ([]*domain.Template, error) {
	return []*domain.Template{}, nil
}

// GetPublicTemplates 获取公开模板（Mock实现）
func (s *MockTemplateService) GetPublicTemplates(ctx context.Context) ([]*domain.Template, error) {
	return []*domain.Template{}, nil
}

// GetTemplatesByCategory 根据分类获取模板（Mock实现）
func (s *MockTemplateService) GetTemplatesByCategory(ctx context.Context, category string) ([]*domain.Template, error) {
	return []*domain.Template{}, nil
}

// ApplyTemplate 应用模板（Mock实现）
func (s *MockTemplateService) ApplyTemplate(ctx context.Context, templateID string, canvasConfig domain.CanvasConfig) (*ApplyTemplateResponse, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// SaveAsTemplate 保存为模板（Mock实现）
func (s *MockTemplateService) SaveAsTemplate(ctx context.Context, req SaveAsTemplateRequest) (*domain.Template, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// CloneTemplate 克隆模板（Mock实现）
func (s *MockTemplateService) CloneTemplate(ctx context.Context, templateID string, newName string) (*domain.Template, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// ExportTemplate 导出模板（Mock实现）
func (s *MockTemplateService) ExportTemplate(ctx context.Context, templateID string) (*ExportTemplateResponse, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}

// ImportTemplate 导入模板（Mock实现）
func (s *MockTemplateService) ImportTemplate(ctx context.Context, req ImportTemplateRequest) (*domain.Template, error) {
	return nil, fmt.Errorf("内存模式下不支持模板功能")
}
