// Package services 鏁版嵁搴撴湇鍔″疄鐜?
package services

import (
	"context"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// DatabaseService 鏁版嵁搴撴湇鍔℃帴鍙?
type DatabaseService interface {
	CreateDatabaseConnection(ctx context.Context, conn *domain.DatabaseConnection) error
	GetDatabaseConnections(ctx context.Context) ([]*domain.DatabaseConnection, error)
	CreateTableMapping(ctx context.Context, mapping *domain.TableMapping) error
	GetTableMappings(ctx context.Context) ([]*domain.TableMapping, error)
}

type databaseService struct {
	dbConnRepo       repositories.DatabaseConnectionRepository
	tableMappingRepo repositories.TableMappingRepository
}

func NewDatabaseService(
	dbConnRepo repositories.DatabaseConnectionRepository,
	tableMappingRepo repositories.TableMappingRepository,
) DatabaseService {
	return &databaseService{
		dbConnRepo:       dbConnRepo,
		tableMappingRepo: tableMappingRepo,
	}
}

func (s *databaseService) CreateDatabaseConnection(ctx context.Context, conn *domain.DatabaseConnection) error {
	return s.dbConnRepo.Create(ctx, conn)
}

func (s *databaseService) GetDatabaseConnections(ctx context.Context) ([]*domain.DatabaseConnection, error) {
	return s.dbConnRepo.List(ctx)
}

func (s *databaseService) CreateTableMapping(ctx context.Context, mapping *domain.TableMapping) error {
	return s.tableMappingRepo.Create(ctx, mapping)
}

func (s *databaseService) GetTableMappings(ctx context.Context) ([]*domain.TableMapping, error) {
	return s.tableMappingRepo.List(ctx)
}
