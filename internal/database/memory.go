// Package database 鍐呭瓨鏁版嵁搴撳疄鐜?
// 鐢ㄤ簬婕旂ず鍜屽紑鍙戠幆澧冿紝涓嶄緷璧栧閮ㄦ暟鎹簱
package database

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"robot-path-editor/internal/config"
	"robot-path-editor/internal/domain"
)

// memoryDatabase 鍐呭瓨鏁版嵁搴撳疄鐜?
type memoryDatabase struct {
	nodes       map[string]*domain.Node
	paths       map[string]*domain.Path
	connections map[string]*domain.DatabaseConnection
	mappings    map[string]*domain.TableMapping
	mu          sync.RWMutex
}

// NewMemoryDatabase 鍒涘缓鍐呭瓨鏁版嵁搴撳疄渚?
func NewMemoryDatabase() Database {
	return &memoryDatabase{
		nodes:       make(map[string]*domain.Node),
		paths:       make(map[string]*domain.Path),
		connections: make(map[string]*domain.DatabaseConnection),
		mappings:    make(map[string]*domain.TableMapping),
	}
}

// NewMemoryDatabaseFromConfig 浠庨厤缃垱寤哄唴瀛樻暟鎹簱锛堝吋瀹规帴鍙ｏ級
func NewMemoryDatabaseFromConfig(cfg config.DatabaseConfig) (Database, error) {
	db := NewMemoryDatabase()

	// 鑷姩杩佺Щ锛堝湪鍐呭瓨鏁版嵁搴撲腑鏄┖鎿嶄綔锛?
	if cfg.AutoMigrate {
		if err := db.AutoMigrate(); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// DB 杩斿洖nil锛堝唴瀛樻暟鎹簱涓嶉渶瑕丟ORM锛?
func (m *memoryDatabase) DB() interface{} {
	return nil
}

// GORMDB 杩斿洖nil锛堝唴瀛樻暟鎹簱涓嶄娇鐢℅ORM锛?
func (m *memoryDatabase) GORMDB() *gorm.DB {
	return nil
}

// Close 鍏抽棴鏁版嵁搴擄紙鍐呭瓨鏁版嵁搴撴棤闇€鍏抽棴锛?
func (m *memoryDatabase) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 娓呯悊鍐呭瓨
	m.nodes = make(map[string]*domain.Node)
	m.paths = make(map[string]*domain.Path)
	m.connections = make(map[string]*domain.DatabaseConnection)
	m.mappings = make(map[string]*domain.TableMapping)

	return nil
}

// Ping 妫€鏌ユ暟鎹簱杩炴帴鐘舵€侊紙鍐呭瓨鏁版嵁搴撴€绘槸鍙敤锛?
func (m *memoryDatabase) Ping(ctx context.Context) error {
	return nil
}

// Transaction 鎵ц浜嬪姟鎿嶄綔锛堝唴瀛樻暟鎹簱绠€鍖栧疄鐜帮級
func (m *memoryDatabase) Transaction(ctx context.Context, fn func(tx interface{}) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 鍦ㄥ唴瀛樻暟鎹簱涓紝鐩存帴鎵ц鍑芥暟
	return fn(m)
}

// AutoMigrate 鑷姩杩佺Щ鏁版嵁搴撶粨鏋勶紙鍐呭瓨鏁版嵁搴撴棤闇€杩佺Щ锛?
func (m *memoryDatabase) AutoMigrate() error {
	// 鍐呭瓨鏁版嵁搴撲笉闇€瑕佽縼绉伙紝鐩存帴杩斿洖鎴愬姛
	return nil
}

// 鑺傜偣鎿嶄綔

// CreateNode 鍒涘缓鑺傜偣
func (m *memoryDatabase) CreateNode(node *domain.Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[string(node.ID)]; exists {
		return fmt.Errorf("鑺傜偣宸插瓨鍦? %s", node.ID)
	}

	m.nodes[string(node.ID)] = node
	return nil
}

// GetNode 鑾峰彇鑺傜偣
func (m *memoryDatabase) GetNode(id string) (*domain.Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, exists := m.nodes[id]
	if !exists {
		return nil, fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", id)
	}

	return node, nil
}

// UpdateNode 鏇存柊鑺傜偣
func (m *memoryDatabase) UpdateNode(node *domain.Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[string(node.ID)]; !exists {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", node.ID)
	}

	m.nodes[string(node.ID)] = node
	return nil
}

// DeleteNode 鍒犻櫎鑺傜偣
func (m *memoryDatabase) DeleteNode(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[id]; !exists {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", id)
	}

	delete(m.nodes, id)
	return nil
}

// ListNodes 鍒楀嚭鎵€鏈夎妭鐐?
func (m *memoryDatabase) ListNodes() ([]*domain.Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*domain.Node, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// 璺緞鎿嶄綔

// CreatePath 鍒涘缓璺緞
func (m *memoryDatabase) CreatePath(path *domain.Path) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.paths[string(path.ID)]; exists {
		return fmt.Errorf("璺緞宸插瓨鍦? %s", path.ID)
	}

	m.paths[string(path.ID)] = path
	return nil
}

// GetPath 鑾峰彇璺緞
func (m *memoryDatabase) GetPath(id string) (*domain.Path, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path, exists := m.paths[id]
	if !exists {
		return nil, fmt.Errorf("璺緞涓嶅瓨鍦? %s", id)
	}

	return path, nil
}

// UpdatePath 鏇存柊璺緞
func (m *memoryDatabase) UpdatePath(path *domain.Path) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.paths[string(path.ID)]; !exists {
		return fmt.Errorf("璺緞涓嶅瓨鍦? %s", path.ID)
	}

	m.paths[string(path.ID)] = path
	return nil
}

// DeletePath 鍒犻櫎璺緞
func (m *memoryDatabase) DeletePath(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.paths[id]; !exists {
		return fmt.Errorf("璺緞涓嶅瓨鍦? %s", id)
	}

	delete(m.paths, id)
	return nil
}

// ListPaths 鍒楀嚭鎵€鏈夎矾寰?
func (m *memoryDatabase) ListPaths() ([]*domain.Path, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	paths := make([]*domain.Path, 0, len(m.paths))
	for _, path := range m.paths {
		paths = append(paths, path)
	}

	return paths, nil
}
