// Package services 数据库服务实现
package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// DatabaseService 数据库业务服务接口
type DatabaseService interface {
	// 数据库连接管理
	CreateConnection(ctx context.Context, req CreateConnectionRequest) (*domain.DatabaseConnection, error)
	GetConnection(ctx context.Context, id string) (*domain.DatabaseConnection, error)
	UpdateConnection(ctx context.Context, req UpdateConnectionRequest) (*domain.DatabaseConnection, error)
	DeleteConnection(ctx context.Context, id string) error
	ListConnections(ctx context.Context) ([]*domain.DatabaseConnection, error)
	TestConnection(ctx context.Context, id string) error

	// 表映射管理
	CreateTableMapping(ctx context.Context, req CreateTableMappingRequest) (*domain.TableMapping, error)
	GetTableMapping(ctx context.Context, id string) (*domain.TableMapping, error)
	UpdateTableMapping(ctx context.Context, req UpdateTableMappingRequest) (*domain.TableMapping, error)
	DeleteTableMapping(ctx context.Context, id string) error
	ListTableMappings(ctx context.Context) ([]*domain.TableMapping, error)
}

// CreateConnectionRequest 创建数据库连接请求
type CreateConnectionRequest struct {
	Name       string            `json:"name" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Host       string            `json:"host" binding:"required"`
	Port       int               `json:"port" binding:"required"`
	Database   string            `json:"database" binding:"required"`
	Username   string            `json:"username" binding:"required"`
	Password   string            `json:"password" binding:"required"`
	Properties map[string]string `json:"properties,omitempty"`
}

// UpdateConnectionRequest 更新数据库连接请求
type UpdateConnectionRequest struct {
	ID         string            `json:"id" binding:"required"`
	Name       *string           `json:"name,omitempty"`
	Type       *string           `json:"type,omitempty"`
	Host       *string           `json:"host,omitempty"`
	Port       *int              `json:"port,omitempty"`
	Database   *string           `json:"database,omitempty"`
	Username   *string           `json:"username,omitempty"`
	Password   *string           `json:"password,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// CreateTableMappingRequest 创建表映射请求
type CreateTableMappingRequest struct {
	ConnectionID string                   `json:"connection_id" binding:"required"`
	TableName    string                   `json:"table_name" binding:"required"`
	NodeMapping  *domain.NodeTableMapping `json:"node_mapping,omitempty"`
	PathMapping  *domain.PathTableMapping `json:"path_mapping,omitempty"`
}

// UpdateTableMappingRequest 更新表映射请求
type UpdateTableMappingRequest struct {
	ID          string                   `json:"id" binding:"required"`
	TableName   *string                  `json:"table_name,omitempty"`
	NodeMapping *domain.NodeTableMapping `json:"node_mapping,omitempty"`
	PathMapping *domain.PathTableMapping `json:"path_mapping,omitempty"`
}

// databaseService 数据库服务实现
type databaseService struct {
	dbConnRepo       repositories.DatabaseConnectionRepository
	tableMappingRepo repositories.TableMappingRepository
}

// NewDatabaseService 创建新的数据库服务实例
func NewDatabaseService(
	dbConnRepo repositories.DatabaseConnectionRepository,
	tableMappingRepo repositories.TableMappingRepository,
) DatabaseService {
	return &databaseService{
		dbConnRepo:       dbConnRepo,
		tableMappingRepo: tableMappingRepo,
	}
}

// CreateConnection 创建数据库连接
func (s *databaseService) CreateConnection(ctx context.Context, req CreateConnectionRequest) (*domain.DatabaseConnection, error) {
	// 创建新的数据库连接对象
	conn := &domain.DatabaseConnection{
		ID:         generateID(),
		Name:       req.Name,
		Type:       req.Type,
		Host:       req.Host,
		Port:       req.Port,
		Database:   req.Database,
		Username:   req.Username,
		Password:   req.Password,
		Properties: req.Properties,
	}

	// 保存到数据库
	err := s.dbConnRepo.Create(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败: %w", err)
	}

	return conn, nil
}

// GetConnection 获取数据库连接
func (s *databaseService) GetConnection(ctx context.Context, id string) (*domain.DatabaseConnection, error) {
	return s.dbConnRepo.GetByID(ctx, id)
}

// UpdateConnection 更新数据库连接
func (s *databaseService) UpdateConnection(ctx context.Context, req UpdateConnectionRequest) (*domain.DatabaseConnection, error) {
	// 获取现有连接
	conn, err := s.dbConnRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 更新非空字段
	if req.Name != nil {
		conn.Name = *req.Name
	}
	if req.Type != nil {
		conn.Type = *req.Type
	}
	if req.Host != nil {
		conn.Host = *req.Host
	}
	if req.Port != nil {
		conn.Port = *req.Port
	}
	if req.Database != nil {
		conn.Database = *req.Database
	}
	if req.Username != nil {
		conn.Username = *req.Username
	}
	if req.Password != nil {
		conn.Password = *req.Password
	}
	if req.Properties != nil {
		conn.Properties = req.Properties
	}

	// 保存更新
	err = s.dbConnRepo.Update(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("更新数据库连接失败: %w", err)
	}

	return conn, nil
}

// DeleteConnection 删除数据库连接
func (s *databaseService) DeleteConnection(ctx context.Context, id string) error {
	return s.dbConnRepo.Delete(ctx, id)
}

// ListConnections 列出数据库连接
func (s *databaseService) ListConnections(ctx context.Context) ([]*domain.DatabaseConnection, error) {
	return s.dbConnRepo.List(ctx)
}

// TestConnection 测试数据库连接
func (s *databaseService) TestConnection(ctx context.Context, id string) error {
	return s.dbConnRepo.TestConnection(ctx, id)
}

// CreateTableMapping 创建表映射
func (s *databaseService) CreateTableMapping(ctx context.Context, req CreateTableMappingRequest) (*domain.TableMapping, error) {
	// 验证连接是否存在
	_, err := s.dbConnRepo.GetByID(ctx, req.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("数据库连接不存在: %w", err)
	}

	// 创建表映射对象
	mapping := &domain.TableMapping{
		ID:           generateID(),
		ConnectionID: req.ConnectionID,
		TableName:    req.TableName,
		NodeMapping:  req.NodeMapping,
		PathMapping:  req.PathMapping,
	}

	// 保存到数据库
	err = s.tableMappingRepo.Create(ctx, mapping)
	if err != nil {
		return nil, fmt.Errorf("创建表映射失败: %w", err)
	}

	return mapping, nil
}

// GetTableMapping 获取表映射
func (s *databaseService) GetTableMapping(ctx context.Context, id string) (*domain.TableMapping, error) {
	return s.tableMappingRepo.GetByID(ctx, id)
}

// UpdateTableMapping 更新表映射
func (s *databaseService) UpdateTableMapping(ctx context.Context, req UpdateTableMappingRequest) (*domain.TableMapping, error) {
	// 获取现有映射
	mapping, err := s.tableMappingRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("获取表映射失败: %w", err)
	}

	// 更新非空字段
	if req.TableName != nil {
		mapping.TableName = *req.TableName
	}
	if req.NodeMapping != nil {
		mapping.NodeMapping = req.NodeMapping
	}
	if req.PathMapping != nil {
		mapping.PathMapping = req.PathMapping
	}

	// 保存更新
	err = s.tableMappingRepo.Update(ctx, mapping)
	if err != nil {
		return nil, fmt.Errorf("更新表映射失败: %w", err)
	}

	return mapping, nil
}

// DeleteTableMapping 删除表映射
func (s *databaseService) DeleteTableMapping(ctx context.Context, id string) error {
	return s.tableMappingRepo.Delete(ctx, id)
}

// generateID 生成唯一ID
func generateID() string {
	return uuid.New().String()
}

// ListTableMappings 列出表映射
func (s *databaseService) ListTableMappings(ctx context.Context) ([]*domain.TableMapping, error) {
	return s.tableMappingRepo.List(ctx)
}
