// Package database 内存数据库实现
// 用于演示和开发环境，不依赖外部数据库
package database

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"robot-path-editor/internal/domain"
)

// memoryDatabase 内存数据库实现
type memoryDatabase struct {
	nodes   map[domain.NodeID]*domain.Node
	paths   map[domain.PathID]*domain.Path
	nodesMu sync.RWMutex
	pathsMu sync.RWMutex
}

// NewMemoryDatabase 创建内存数据库实例
func NewMemoryDatabase() Database {
	return &memoryDatabase{
		nodes: make(map[domain.NodeID]*domain.Node),
		paths: make(map[domain.PathID]*domain.Path),
	}
}

// Connect 连接数据库（内存数据库无需连接）
func (db *memoryDatabase) Connect() error {
	return nil
}

// Close 关闭数据库连接（内存数据库无需关闭）
func (db *memoryDatabase) Close() error {
	return nil
}

// AutoMigrate 自动迁移（在内存数据库中是空操作）
func (db *memoryDatabase) AutoMigrate(dst ...interface{}) error {
	// 内存数据库不需要迁移
	return nil
}

// Transaction 事务执行
func (db *memoryDatabase) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	// 简化实现：内存数据库直接执行
	return fn(db)
}

// DB 返回nil（内存数据库不需要GORM）
func (db *memoryDatabase) DB() interface{} {
	return nil
}

// GORMDB 返回nil（内存数据库不使用GORM）
func (db *memoryDatabase) GORMDB() *gorm.DB {
	return nil
}

// Ping 检查数据库连接状态（内存数据库总是可用）
func (db *memoryDatabase) Ping() error {
	return nil
}

// === 节点操作 ===

// CreateNode 创建节点
func (db *memoryDatabase) CreateNode(node *domain.Node) error {
	db.nodesMu.Lock()
	defer db.nodesMu.Unlock()

	if _, exists := db.nodes[node.ID]; exists {
		return fmt.Errorf("节点已存在: %s", node.ID)
	}

	// 创建副本以避免外部修改
	nodeCopy := *node
	db.nodes[node.ID] = &nodeCopy
	return nil
}

// GetNode 获取节点
func (db *memoryDatabase) GetNode(id domain.NodeID) (*domain.Node, error) {
	db.nodesMu.RLock()
	defer db.nodesMu.RUnlock()

	node, exists := db.nodes[id]
	if !exists {
		return nil, fmt.Errorf("节点不存在: %s", id)
	}

	// 返回副本以避免并发修改
	nodeCopy := *node
	return &nodeCopy, nil
}

// UpdateNode 更新节点
func (db *memoryDatabase) UpdateNode(node *domain.Node) error {
	db.nodesMu.Lock()
	defer db.nodesMu.Unlock()

	if _, exists := db.nodes[node.ID]; !exists {
		return fmt.Errorf("节点不存在: %s", node.ID)
	}

	// 创建副本
	nodeCopy := *node
	db.nodes[node.ID] = &nodeCopy
	return nil
}

// DeleteNode 删除节点
func (db *memoryDatabase) DeleteNode(id domain.NodeID) error {
	db.nodesMu.Lock()
	defer db.nodesMu.Unlock()

	if _, exists := db.nodes[id]; !exists {
		return fmt.Errorf("节点不存在: %s", id)
	}

	delete(db.nodes, id)
	return nil
}

// ListNodes 列出所有节点
func (db *memoryDatabase) ListNodes() ([]*domain.Node, error) {
	db.nodesMu.RLock()
	defer db.nodesMu.RUnlock()

	nodes := make([]*domain.Node, 0, len(db.nodes))
	for _, node := range db.nodes {
		// 返回副本
		nodeCopy := *node
		nodes = append(nodes, &nodeCopy)
	}

	return nodes, nil
}

// === 路径操作 ===

// CreatePath 创建路径
func (db *memoryDatabase) CreatePath(path *domain.Path) error {
	db.pathsMu.Lock()
	defer db.pathsMu.Unlock()

	if _, exists := db.paths[path.ID]; exists {
		return fmt.Errorf("路径已存在: %s", path.ID)
	}

	// 创建副本
	pathCopy := *path
	db.paths[path.ID] = &pathCopy
	return nil
}

// GetPath 获取路径
func (db *memoryDatabase) GetPath(id domain.PathID) (*domain.Path, error) {
	db.pathsMu.RLock()
	defer db.pathsMu.RUnlock()

	path, exists := db.paths[id]
	if !exists {
		return nil, fmt.Errorf("路径不存在: %s", id)
	}

	// 返回副本
	pathCopy := *path
	return &pathCopy, nil
}

// UpdatePath 更新路径
func (db *memoryDatabase) UpdatePath(path *domain.Path) error {
	db.pathsMu.Lock()
	defer db.pathsMu.Unlock()

	if _, exists := db.paths[path.ID]; !exists {
		return fmt.Errorf("路径不存在: %s", path.ID)
	}

	// 创建副本
	pathCopy := *path
	db.paths[path.ID] = &pathCopy
	return nil
}

// DeletePath 删除路径
func (db *memoryDatabase) DeletePath(id domain.PathID) error {
	db.pathsMu.Lock()
	defer db.pathsMu.Unlock()

	if _, exists := db.paths[id]; !exists {
		return fmt.Errorf("路径不存在: %s", id)
	}

	delete(db.paths, id)
	return nil
}

// ListPaths 列出所有路径
func (db *memoryDatabase) ListPaths() ([]*domain.Path, error) {
	db.pathsMu.RLock()
	defer db.pathsMu.RUnlock()

	paths := make([]*domain.Path, 0, len(db.paths))
	for _, path := range db.paths {
		// 返回副本
		pathCopy := *path
		paths = append(paths, &pathCopy)
	}

	return paths, nil
}
