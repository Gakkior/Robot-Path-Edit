// Package services 数据同步服务
package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// DataSyncService 数据同步服务接口
type DataSyncService interface {
	// 从外部数据库同步节点数据
	SyncNodesFromExternal(ctx context.Context, mappingID string) (*SyncResult, error)
	// 从外部数据库同步路径数据
	SyncPathsFromExternal(ctx context.Context, mappingID string) (*SyncResult, error)
	// 全量同步数据
	SyncAllDataFromExternal(ctx context.Context, mappingID string) (*SyncResult, error)
	// 验证外部数据库表结构
	ValidateExternalTable(ctx context.Context, connectionID, tableName string) (*TableValidationResult, error)
}

// SyncResult 同步结果
type SyncResult struct {
	NodesCreated int      `json:"nodes_created"`
	NodesUpdated int      `json:"nodes_updated"`
	PathsCreated int      `json:"paths_created"`
	PathsUpdated int      `json:"paths_updated"`
	Errors       []string `json:"errors,omitempty"`
}

// TableValidationResult 表验证结果
type TableValidationResult struct {
	Valid   bool     `json:"valid"`
	Columns []string `json:"columns"`
	Message string   `json:"message,omitempty"`
}

// dataSyncService 数据同步服务实现
type dataSyncService struct {
	dbConnRepo       repositories.DatabaseConnectionRepository
	tableMappingRepo repositories.TableMappingRepository
	nodeRepo         repositories.NodeRepository
	pathRepo         repositories.PathRepository
}

// NewDataSyncService 创建新的数据同步服务实例
func NewDataSyncService(
	dbConnRepo repositories.DatabaseConnectionRepository,
	tableMappingRepo repositories.TableMappingRepository,
	nodeRepo repositories.NodeRepository,
	pathRepo repositories.PathRepository,
) DataSyncService {
	return &dataSyncService{
		dbConnRepo:       dbConnRepo,
		tableMappingRepo: tableMappingRepo,
		nodeRepo:         nodeRepo,
		pathRepo:         pathRepo,
	}
}

// SyncNodesFromExternal 从外部数据库同步节点数据
func (s *dataSyncService) SyncNodesFromExternal(ctx context.Context, mappingID string) (*SyncResult, error) {
	result := &SyncResult{}

	// 获取表映射配置
	mapping, err := s.tableMappingRepo.GetByID(ctx, mappingID)
	if err != nil {
		return nil, fmt.Errorf("获取表映射失败: %w", err)
	}

	if mapping.NodeMapping == nil {
		return nil, fmt.Errorf("表映射中未配置节点映射")
	}

	// 获取数据库连接配置
	conn, err := s.dbConnRepo.GetByID(ctx, mapping.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 连接外部数据库
	externalDB, err := s.connectToExternalDB(conn)
	if err != nil {
		return nil, fmt.Errorf("连接外部数据库失败: %w", err)
	}
	defer externalDB.Close()

	// 构建查询SQL
	query := s.buildNodeSelectQuery(mapping)

	// 执行查询
	rows, err := externalDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询外部数据库失败: %w", err)
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		node, err := s.scanNodeFromRow(rows, mapping.NodeMapping)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}

		// 检查节点是否已存在
		existingNode, err := s.nodeRepo.GetByID(ctx, node.ID)
		if err != nil {
			// 节点不存在，创建新节点
			err = s.nodeRepo.Create(ctx, node)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("创建节点失败: %v", err))
			} else {
				result.NodesCreated++
			}
		} else {
			// 节点存在，更新节点
			node.Metadata = existingNode.Metadata
			node.Metadata.UpdatedAt = time.Now()
			node.Metadata.Version++
			err = s.nodeRepo.Update(ctx, node)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("更新节点失败: %v", err))
			} else {
				result.NodesUpdated++
			}
		}
	}

	return result, nil
}

// SyncPathsFromExternal 从外部数据库同步路径数据
func (s *dataSyncService) SyncPathsFromExternal(ctx context.Context, mappingID string) (*SyncResult, error) {
	result := &SyncResult{}

	// 获取表映射配置
	mapping, err := s.tableMappingRepo.GetByID(ctx, mappingID)
	if err != nil {
		return nil, fmt.Errorf("获取表映射失败: %w", err)
	}

	if mapping.PathMapping == nil {
		return nil, fmt.Errorf("表映射中未配置路径映射")
	}

	// 获取数据库连接配置
	conn, err := s.dbConnRepo.GetByID(ctx, mapping.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 连接外部数据库
	externalDB, err := s.connectToExternalDB(conn)
	if err != nil {
		return nil, fmt.Errorf("连接外部数据库失败: %w", err)
	}
	defer externalDB.Close()

	// 构建查询SQL
	query := s.buildPathSelectQuery(mapping)

	// 执行查询
	rows, err := externalDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询外部数据库失败: %w", err)
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		path, err := s.scanPathFromRow(rows, mapping.PathMapping)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			continue
		}

		// 检查路径是否已存在
		existingPath, err := s.pathRepo.GetByID(ctx, path.ID)
		if err != nil {
			// 路径不存在，创建新路径
			err = s.pathRepo.Create(ctx, path)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("创建路径失败: %v", err))
			} else {
				result.PathsCreated++
			}
		} else {
			// 路径存在，更新路径
			path.Metadata = existingPath.Metadata
			path.Metadata.UpdatedAt = time.Now()
			path.Metadata.Version++
			err = s.pathRepo.Update(ctx, path)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("更新路径失败: %v", err))
			} else {
				result.PathsUpdated++
			}
		}
	}

	return result, nil
}

// SyncAllDataFromExternal 全量同步数据
func (s *dataSyncService) SyncAllDataFromExternal(ctx context.Context, mappingID string) (*SyncResult, error) {
	totalResult := &SyncResult{}

	// 同步节点数据
	nodeResult, err := s.SyncNodesFromExternal(ctx, mappingID)
	if err != nil {
		return nil, fmt.Errorf("同步节点数据失败: %w", err)
	}
	totalResult.NodesCreated = nodeResult.NodesCreated
	totalResult.NodesUpdated = nodeResult.NodesUpdated
	totalResult.Errors = append(totalResult.Errors, nodeResult.Errors...)

	// 同步路径数据
	pathResult, err := s.SyncPathsFromExternal(ctx, mappingID)
	if err != nil {
		return nil, fmt.Errorf("同步路径数据失败: %w", err)
	}
	totalResult.PathsCreated = pathResult.PathsCreated
	totalResult.PathsUpdated = pathResult.PathsUpdated
	totalResult.Errors = append(totalResult.Errors, pathResult.Errors...)

	return totalResult, nil
}

// ValidateExternalTable 验证外部数据库表结构
func (s *dataSyncService) ValidateExternalTable(ctx context.Context, connectionID, tableName string) (*TableValidationResult, error) {
	// 获取数据库连接配置
	conn, err := s.dbConnRepo.GetByID(ctx, connectionID)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 连接外部数据库
	externalDB, err := s.connectToExternalDB(conn)
	if err != nil {
		return &TableValidationResult{
			Valid:   false,
			Message: fmt.Sprintf("连接数据库失败: %v", err),
		}, nil
	}
	defer externalDB.Close()

	// 获取表结构信息
	var query string
	switch conn.Type {
	case "mysql":
		query = fmt.Sprintf("DESCRIBE %s", tableName)
	case "sqlite", "sqlite3":
		query = fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	default:
		return &TableValidationResult{
			Valid:   false,
			Message: fmt.Sprintf("不支持的数据库类型: %s", conn.Type),
		}, nil
	}

	rows, err := externalDB.Query(query)
	if err != nil {
		return &TableValidationResult{
			Valid:   false,
			Message: fmt.Sprintf("查询表结构失败: %v", err),
		}, nil
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var columnName string
		var others []interface{}

		// 根据数据库类型处理不同的列结构
		switch conn.Type {
		case "mysql":
			// MySQL DESCRIBE 结果: Field, Type, Null, Key, Default, Extra
			var field, typ, null, key, defaultVal, extra sql.NullString
			err = rows.Scan(&field, &typ, &null, &key, &defaultVal, &extra)
			if err == nil && field.Valid {
				columnName = field.String
			}
		case "sqlite", "sqlite3":
			// SQLite PRAGMA table_info 结果: cid, name, type, notnull, dflt_value, pk
			var cid int
			var name, typ sql.NullString
			var notnull, pk int
			var dfltValue sql.NullString
			err = rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk)
			if err == nil && name.Valid {
				columnName = name.String
			}
		default:
			// 通用处理，尝试扫描第一列作为列名
			cols, _ := rows.Columns()
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			err = rows.Scan(valuePtrs...)
			if err == nil && len(values) > 0 {
				if str, ok := values[0].(string); ok {
					columnName = str
				}
			}
		}

		if err != nil {
			continue
		}

		if columnName != "" {
			columns = append(columns, columnName)
		}

		// 清理others变量以避免未使用的警告
		_ = others
	}

	return &TableValidationResult{
		Valid:   len(columns) > 0,
		Columns: columns,
		Message: fmt.Sprintf("找到 %d 个列", len(columns)),
	}, nil
}

// connectToExternalDB 连接到外部数据库
func (s *dataSyncService) connectToExternalDB(conn *domain.DatabaseConnection) (*sql.DB, error) {
	var dsn string
	switch conn.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conn.Username, conn.Password, conn.Host, conn.Port, conn.Database)
	case "sqlite", "sqlite3":
		dsn = conn.Database
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", conn.Type)
	}

	db, err := sql.Open(conn.Type, dsn)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// buildNodeSelectQuery 构建节点查询SQL
func (s *dataSyncService) buildNodeSelectQuery(mapping *domain.TableMapping) string {
	nodeMapping := mapping.NodeMapping
	query := fmt.Sprintf("SELECT %s", nodeMapping.IDField)

	if nodeMapping.NameField != "" {
		query += fmt.Sprintf(", %s", nodeMapping.NameField)
	}
	if nodeMapping.TypeField != "" {
		query += fmt.Sprintf(", %s", nodeMapping.TypeField)
	}
	if nodeMapping.XField != "" {
		query += fmt.Sprintf(", %s", nodeMapping.XField)
	}
	if nodeMapping.YField != "" {
		query += fmt.Sprintf(", %s", nodeMapping.YField)
	}
	if nodeMapping.ZField != "" {
		query += fmt.Sprintf(", %s", nodeMapping.ZField)
	}

	query += fmt.Sprintf(" FROM %s", mapping.TableName)
	return query
}

// buildPathSelectQuery 构建路径查询SQL
func (s *dataSyncService) buildPathSelectQuery(mapping *domain.TableMapping) string {
	pathMapping := mapping.PathMapping
	query := fmt.Sprintf("SELECT %s, %s, %s",
		pathMapping.IDField,
		pathMapping.StartNodeField,
		pathMapping.EndNodeField)

	if pathMapping.NameField != "" {
		query += fmt.Sprintf(", %s", pathMapping.NameField)
	}
	if pathMapping.WeightField != "" {
		query += fmt.Sprintf(", %s", pathMapping.WeightField)
	}

	query += fmt.Sprintf(" FROM %s", mapping.TableName)
	return query
}

// scanNodeFromRow 从查询结果扫描节点数据
func (s *dataSyncService) scanNodeFromRow(rows *sql.Rows, nodeMapping *domain.NodeTableMapping) (*domain.Node, error) {
	var id string
	var name, nodeType sql.NullString
	var x, y, z sql.NullFloat64

	// 准备扫描目标
	scanArgs := []interface{}{&id}

	if nodeMapping.NameField != "" {
		scanArgs = append(scanArgs, &name)
	}
	if nodeMapping.TypeField != "" {
		scanArgs = append(scanArgs, &nodeType)
	}
	if nodeMapping.XField != "" {
		scanArgs = append(scanArgs, &x)
	}
	if nodeMapping.YField != "" {
		scanArgs = append(scanArgs, &y)
	}
	if nodeMapping.ZField != "" {
		scanArgs = append(scanArgs, &z)
	}

	// 扫描数据
	if err := rows.Scan(scanArgs...); err != nil {
		return nil, fmt.Errorf("扫描节点数据失败: %w", err)
	}

	// 创建节点对象
	node := &domain.Node{
		ID: domain.NodeID(id),
		Position: domain.Position{
			X: x.Float64,
			Y: y.Float64,
			Z: z.Float64,
		},
	}

	if name.Valid {
		node.Name = name.String
	} else {
		node.Name = id // 默认使用ID作为名称
	}

	if nodeType.Valid {
		node.Type = domain.NodeType(nodeType.String)
	} else {
		node.Type = domain.NodeTypePoint
	}

	node.Status = domain.NodeStatusActive

	// 设置默认样式
	node.Style = domain.NodeStyle{
		Color:       "#007bff",
		Size:        10.0,
		Shape:       "circle",
		BorderColor: "#000000",
		BorderWidth: 1.0,
		Opacity:     1.0,
	}

	return node, nil
}

// scanPathFromRow 从查询结果扫描路径数据
func (s *dataSyncService) scanPathFromRow(rows *sql.Rows, pathMapping *domain.PathTableMapping) (*domain.Path, error) {
	var id, startNodeID, endNodeID string
	var name sql.NullString
	var weight sql.NullFloat64

	// 准备扫描目标
	scanArgs := []interface{}{&id, &startNodeID, &endNodeID}

	if pathMapping.NameField != "" {
		scanArgs = append(scanArgs, &name)
	}
	if pathMapping.WeightField != "" {
		scanArgs = append(scanArgs, &weight)
	}

	// 扫描数据
	if err := rows.Scan(scanArgs...); err != nil {
		return nil, fmt.Errorf("扫描路径数据失败: %w", err)
	}

	// 创建路径对象
	path := &domain.Path{
		ID:          domain.PathID(id),
		StartNodeID: domain.NodeID(startNodeID),
		EndNodeID:   domain.NodeID(endNodeID),
		Weight:      weight.Float64,
		Type:        domain.PathTypeNormal,
		Status:      domain.PathStatusActive,
		Direction:   "bidirectional",
		CurveType:   domain.CurveTypeLinear,
	}

	if name.Valid {
		path.Name = name.String
	} else {
		path.Name = fmt.Sprintf("路径_%s", id)
	}

	// 设置默认样式
	path.Style = domain.PathStyle{
		Color:   "#6c757d",
		Width:   2.0,
		Style:   "solid",
		Opacity: 1.0,
	}

	return path, nil
}
