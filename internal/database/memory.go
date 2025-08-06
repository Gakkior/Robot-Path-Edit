// Package database 内存数据库实现
// 用于演示和开发环境，不依赖外部数据库
package database

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"robot-path-editor/internal/config"
	"robot-path-editor/internal/domain"
)

// memoryDatabase 内存数据库实�?
type memoryDatabase struct {
	nodes       map[string]*domain.Node
	paths       map[string]*domain.Path
	connections map[string]*domain.DatabaseConnection
	mappings    map[string]*domain.TableMapping
	mu          sync.RWMutex
}

// NewMemoryDatabase 创建内存数据库实�?
func NewMemoryDatabase() Database {
	return &memoryDatabase{
		nodes:       make(map[string]*domain.Node),
		paths:       make(map[string]*domain.Path),
		connections: make(map[string]*domain.DatabaseConnection),
		mappings:    make(map[string]*domain.TableMapping),
	}
}

// NewMemoryDatabaseFromConfig 从配置创建内存数据库（兼容接口）
func NewMemoryDatabaseFromConfig(cfg config.DatabaseConfig) (Database, error) {
	db := NewMemoryDatabase()

	// 自动迁移（在内存数据库中是空操作�?
	if cfg.AutoMigrate {
		if err := db.AutoMigrate(); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// DB 返回nil（内存数据库不需要GORM�?
func (m *memoryDatabase) DB() interface{} {
	return nil
}

// GORMDB 返回nil（内存数捺�不使用GORM�?
func (m *memoryDatabase) GORMDB() *gorm.DB {
	return nil
}

// Close 关闭数据库（内存数据库无需关闭�?
func (m *memoryDatabase) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清理内存
	m.nodes = make(map[string]*domain.Node)
	m.paths = make(map[string]*domain.Path)
	m.connections = make(map[string]*domain.DatabaseConnection)
	m.mappings = make(map[string]*domain.TableMapping)

	return nil
}

// Ping 检查数据库连接状态（内存数据库总是可用�?
func (m *memoryDatabase) Ping(ctx context.Context) error {
	return nil
}

// Transaction 执行事务操作（内存数据库简化实现）
func (m *memoryDatabase) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 在内存数据库中，直接执行函数
	return fn(m)
}

// AutoMigrate 自动迁移数据库结构（内存数据库无需迁移�?
func (m *memoryDatabase) AutoMigrate() error {
	// 内存数据库不需要迁移，直接返回成功
	return nil
}

// 节点操作

// CreateNode 创建节点
func (m *memoryDatabase) CreateNode(node *domain.Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[string(node.ID)]; exists {
		return fmt.Errorf("节点已存�? %s", node.ID)
	}

	m.nodes[string(node.ID)] = node
	return nil
}

// GetNode 获取节点
func (m *memoryDatabase) GetNode(id string) (*domain.Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, exists := m.nodes[id]
	if !exists {
		return nil, fmt.Errorf("节点不存�? %s", id)
	}

	return node, nil
}

// UpdateNode 更新节点
func (m *memoryDatabase) UpdateNode(node *domain.Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[string(node.ID)]; !exists {
		return fmt.Errorf("节点不存�? %s", node.ID)
	}

	m.nodes[string(node.ID)] = node
	return nil
}

// DeleteNode 删除节点
func (m *memoryDatabase) DeleteNode(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[id]; !exists {
		return fmt.Errorf("节点不存�? %s", id)
	}

	delete(m.nodes, id)
	return nil
}

// ListNodes 列出所有节�?
func (m *memoryDatabase) ListNodes() ([]*domain.Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*domain.Node, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// 路径操作

// CreatePath 创建路径
func (m *memoryDatabase) CreatePath(path *domain.Path) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.paths[string(path.ID)]; exists {
		return fmt.Errorf("路径已存�? %s", path.ID)
	}

	m.paths[string(path.ID)] = path
	return nil
}

// GetPath 获取路径
func (m *memoryDatabase) GetPath(id string) (*domain.Path, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path, exists := m.paths[id]
	if !exists {
		return nil, fmt.Errorf("路径不存�? %s", id)
	}

	return path, nil
}

// UpdatePath 更新路径
func (m *memoryDatabase) UpdatePath(path *domain.Path) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.paths[string(path.ID)]; !exists {
		return fmt.Errorf("路径不存�? %s", path.ID)
	}

	m.paths[string(path.ID)] = path
	return nil
}

// DeletePath 删除路径
func (m *memoryDatabase) DeletePath(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.paths[id]; !exists {
		return fmt.Errorf("路径不存�? %s", id)
	}

	delete(m.paths, id)
	return nil
}

// ListPaths 列出所有路�?
func (m *memoryDatabase) ListPaths() ([]*domain.Path, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	paths := make([]*domain.Path, 0, len(m.paths))
	for _, path := range m.paths {
		paths = append(paths, path)
	}

	return paths, nil
}
