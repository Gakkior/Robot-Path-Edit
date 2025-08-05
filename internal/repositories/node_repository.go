// Package repositories 瀹炵幇鏁版嵁璁块棶灞?
//
// 璁捐鍙傝€冿細
// - DDD鐨勪粨鍌ㄦā寮?
// - Kubernetes鐨勫瓨鍌ㄦ娊璞?
// - GitHub鐨勪粨鍌ㄥ疄鐜版ā寮?
//
// 鐗圭偣锛?
// 1. 鎺ュ彛鎶借薄锛氬畾涔夋竻鏅扮殑鏁版嵁璁块棶鎺ュ彛
// 2. 瀹炵幇鍒嗙锛氭敮鎸佷笉鍚岀殑瀛樺偍鍚庣
// 3. 鏌ヨ浼樺寲锛氭敮鎸佸鏉傛煡璇㈠拰鍒嗛〉
// 4. 缂撳瓨鍙嬪ソ锛氳璁′究浜庣紦瀛橀泦鎴?
package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// NodeRepository 鑺傜偣浠撳偍鎺ュ彛
// 瀹氫箟鑺傜偣鏁版嵁璁块棶鐨勬墍鏈夋搷浣?
type NodeRepository interface {
	// 鍩虹CRUD鎿嶄綔
	Create(ctx context.Context, node *domain.Node) error
	GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error)
	Update(ctx context.Context, node *domain.Node) error
	Delete(ctx context.Context, id domain.NodeID) error

	// 鎵归噺鎿嶄綔
	CreateBatch(ctx context.Context, nodes []*domain.Node) error
	GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error)
	UpdateBatch(ctx context.Context, nodes []*domain.Node) error
	DeleteBatch(ctx context.Context, ids []domain.NodeID) error

	// 鏌ヨ鎿嶄綔
	List(ctx context.Context, options ListOptions) ([]*domain.Node, error)
	Count(ctx context.Context, filter NodeFilter) (int64, error)

	// 绌洪棿鏌ヨ - 鐢ㄤ簬鐢诲竷鎿嶄綔
	GetByArea(ctx context.Context, minX, minY, maxX, maxY float64) ([]*domain.Node, error)
	GetNearby(ctx context.Context, position domain.Position, radius float64) ([]*domain.Node, error)

	// 鍏崇郴鏌ヨ
	GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error)
	GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error)

	// 鍏冩暟鎹煡璇?
	GetByLabels(ctx context.Context, labels map[string]string) ([]*domain.Node, error)
	GetByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error)
	GetByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error)
}

// NodeFilter 鑺傜偣鏌ヨ杩囨护鍣?
type NodeFilter struct {
	IDs    []domain.NodeID   `json:"ids,omitempty"`
	Name   string            `json:"name,omitempty"`
	Type   domain.NodeType   `json:"type,omitempty"`
	Status domain.NodeStatus `json:"status,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`

	// 浣嶇疆杩囨护
	MinX *float64 `json:"min_x,omitempty"`
	MaxX *float64 `json:"max_x,omitempty"`
	MinY *float64 `json:"min_y,omitempty"`
	MaxY *float64 `json:"max_y,omitempty"`
	MinZ *float64 `json:"min_z,omitempty"`
	MaxZ *float64 `json:"max_z,omitempty"`
}

// ListOptions 鍒楄〃鏌ヨ閫夐」
type ListOptions struct {
	Filter   NodeFilter `json:"filter"`
	Page     int        `json:"page"`      // 椤电爜锛屼粠1寮€濮?
	PageSize int        `json:"page_size"` // 椤靛ぇ灏?
	OrderBy  string     `json:"order_by"`  // 鎺掑簭瀛楁
	Order    string     `json:"order"`     // 鎺掑簭鏂瑰悜: asc, desc
}

// nodeRepository 鑺傜偣浠撳偍瀹炵幇
type nodeRepository struct {
	db database.Database
}

// NewNodeRepository 鍒涘缓鑺傜偣浠撳偍瀹炰緥
func NewNodeRepository(db database.Database) NodeRepository {
	return &nodeRepository{
		db: db,
	}
}

// Create 鍒涘缓鑺傜偣
func (r *nodeRepository) Create(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
	}

	// 妫€鏌ユ槸鍚︿负GORM鏁版嵁搴?
	if gormDB := r.db.GORMDB(); gormDB != nil {
		return gormDB.WithContext(ctx).Create(node).Error
	}

	// 濡傛灉鏄唴瀛樻暟鎹簱锛岄渶瑕佺被鍨嬫柇瑷€
	if memDB, ok := r.db.(interface {
		CreateNode(*domain.Node) error
	}); ok {
		return memDB.CreateNode(node)
	}

	return fmt.Errorf("涓嶆敮鎸佺殑鏁版嵁搴撶被鍨?)
}

// GetByID 鏍规嵁ID鑾峰彇鑺傜偣
func (r *nodeRepository) GetByID(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	var node domain.Node
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&node).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", id)
		}
		return nil, err
	}
	return &node, nil
}

// Update 鏇存柊鑺傜偣
func (r *nodeRepository) Update(ctx context.Context, node *domain.Node) error {
	if err := node.IsValid(); err != nil {
		return fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
	}

	result := r.db.DB().WithContext(ctx).Save(node)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", node.ID)
	}

	return nil
}

// Delete 鍒犻櫎鑺傜偣
func (r *nodeRepository) Delete(ctx context.Context, id domain.NodeID) error {
	result := r.db.DB().WithContext(ctx).Delete(&domain.Node{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %s", id)
	}

	return nil
}

// CreateBatch 鎵归噺鍒涘缓鑺傜偣
func (r *nodeRepository) CreateBatch(ctx context.Context, nodes []*domain.Node) error {
	// 楠岃瘉鎵€鏈夎妭鐐?
	for _, node := range nodes {
		if err := node.IsValid(); err != nil {
			return fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
		}
	}

	// 鎵归噺鎻掑叆 - 浣跨敤浜嬪姟纭繚涓€鑷存€?
	return r.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).CreateInBatches(nodes, 100).Error
	})
}

// GetByIDs 鏍规嵁ID鍒楄〃鑾峰彇鑺傜偣
func (r *nodeRepository) GetByIDs(ctx context.Context, ids []domain.NodeID) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 杞崲涓哄瓧绗︿覆鍒囩墖
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	err := r.db.DB().WithContext(ctx).Where("id IN ?", stringIDs).Find(&nodes).Error
	return nodes, err
}

// UpdateBatch 鎵归噺鏇存柊鑺傜偣
func (r *nodeRepository) UpdateBatch(ctx context.Context, nodes []*domain.Node) error {
	return r.db.Transaction(ctx, func(tx *gorm.DB) error {
		for _, node := range nodes {
			if err := node.IsValid(); err != nil {
				return fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
			}

			if err := tx.WithContext(ctx).Save(node).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteBatch 鎵归噺鍒犻櫎鑺傜偣
func (r *nodeRepository) DeleteBatch(ctx context.Context, ids []domain.NodeID) error {
	stringIDs := make([]string, len(ids))
	for i, id := range ids {
		stringIDs[i] = string(id)
	}

	return r.db.GORMDB().WithContext(ctx).Delete(&domain.Node{}, "id IN ?", stringIDs).Error
}

// List 鍒楄〃鏌ヨ鑺傜偣
func (r *nodeRepository) List(ctx context.Context, options ListOptions) ([]*domain.Node, error) {
	var nodes []*domain.Node

	query := r.db.DB().WithContext(ctx)

	// 搴旂敤杩囨护鍣?
	query = r.applyFilter(query, options.Filter)

	// 搴旂敤鎺掑簭
	if options.OrderBy != "" {
		order := "asc"
		if options.Order == "desc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", options.OrderBy, order))
	} else {
		query = query.Order("created_at desc") // 榛樿鎸夊垱寤烘椂闂撮檷搴?
	}

	// 搴旂敤鍒嗛〉
	if options.PageSize > 0 {
		offset := 0
		if options.Page > 1 {
			offset = (options.Page - 1) * options.PageSize
		}
		query = query.Offset(offset).Limit(options.PageSize)
	}

	err := query.Find(&nodes).Error
	return nodes, err
}

// Count 缁熻鑺傜偣鏁伴噺
func (r *nodeRepository) Count(ctx context.Context, filter NodeFilter) (int64, error) {
	var count int64

	query := r.db.DB().WithContext(ctx).Model(&domain.Node{})
	query = r.applyFilter(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// GetByArea 鑾峰彇鎸囧畾鍖哄煙鍐呯殑鑺傜偣
func (r *nodeRepository) GetByArea(ctx context.Context, minX, minY, maxX, maxY float64) ([]*domain.Node, error) {
	var nodes []*domain.Node

	err := r.db.DB().WithContext(ctx).
		Where("pos_x BETWEEN ? AND ?", minX, maxX).
		Where("pos_y BETWEEN ? AND ?", minY, maxY).
		Find(&nodes).Error

	return nodes, err
}

// GetNearby 鑾峰彇鎸囧畾浣嶇疆闄勮繎鐨勮妭鐐?
func (r *nodeRepository) GetNearby(ctx context.Context, position domain.Position, radius float64) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 浣跨敤绠€鍗曠殑鐭╁舰鑼冨洿鏌ヨ锛堝彲浼樺寲涓虹湡姝ｇ殑鍦嗗舰鑼冨洿锛?
	minX := position.X - radius
	maxX := position.X + radius
	minY := position.Y - radius
	maxY := position.Y + radius

	err := r.db.DB().WithContext(ctx).
		Where("pos_x BETWEEN ? AND ?", minX, maxX).
		Where("pos_y BETWEEN ? AND ?", minY, maxY).
		Find(&nodes).Error

	// TODO: 鍦ㄥ簲鐢ㄥ眰杩囨护鍑虹湡姝ｅ湪鍦嗗舰鑼冨洿鍐呯殑鑺傜偣
	return nodes, err
}

// GetConnectedNodes 鑾峰彇涓庢寚瀹氳妭鐐硅繛鎺ョ殑鎵€鏈夎妭鐐?
func (r *nodeRepository) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 閫氳繃璺緞琛ㄥ叧鑱旀煡璇?
	err := r.db.DB().WithContext(ctx).
		Joins("JOIN paths ON (nodes.id = paths.start_node_id OR nodes.id = paths.end_node_id)").
		Where("(paths.start_node_id = ? OR paths.end_node_id = ?) AND nodes.id != ?", nodeID, nodeID, nodeID).
		Where("paths.status = ?", domain.PathStatusActive).
		Distinct().
		Find(&nodes).Error

	return nodes, err
}

// GetIsolatedNodes 鑾峰彇瀛ょ珛鑺傜偣锛堟病鏈夎繛鎺ョ殑鑺傜偣锛?
func (r *nodeRepository) GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error) {
	var nodes []*domain.Node

	// 宸﹁繛鎺ヨ矾寰勮〃锛屾煡鎵炬病鏈夎矾寰勭殑鑺傜偣
	err := r.db.DB().WithContext(ctx).
		Where("NOT EXISTS (SELECT 1 FROM paths WHERE nodes.id = paths.start_node_id OR nodes.id = paths.end_node_id)").
		Find(&nodes).Error

	return nodes, err
}

// GetByLabels 鏍规嵁鏍囩鏌ヨ鑺傜偣
func (r *nodeRepository) GetByLabels(ctx context.Context, labels map[string]string) ([]*domain.Node, error) {
	var nodes []*domain.Node

	query := r.db.DB().WithContext(ctx)

	// 浣跨敤JSON鏌ヨ锛堥渶瑕佹暟鎹簱鏀寔锛?
	for key, value := range labels {
		query = query.Where("JSON_EXTRACT(labels, ?) = ?", "$."+key, value)
	}

	err := query.Find(&nodes).Error
	return nodes, err
}

// GetByType 鏍规嵁绫诲瀷鏌ヨ鑺傜偣
func (r *nodeRepository) GetByType(ctx context.Context, nodeType domain.NodeType) ([]*domain.Node, error) {
	var nodes []*domain.Node
	err := r.db.DB().WithContext(ctx).Where("type = ?", nodeType).Find(&nodes).Error
	return nodes, err
}

// GetByStatus 鏍规嵁鐘舵€佹煡璇㈣妭鐐?
func (r *nodeRepository) GetByStatus(ctx context.Context, status domain.NodeStatus) ([]*domain.Node, error) {
	var nodes []*domain.Node
	err := r.db.DB().WithContext(ctx).Where("status = ?", status).Find(&nodes).Error
	return nodes, err
}

// applyFilter 搴旂敤鏌ヨ杩囨护鍣?
func (r *nodeRepository) applyFilter(query *gorm.DB, filter NodeFilter) *gorm.DB {
	// ID杩囨护
	if len(filter.IDs) > 0 {
		stringIDs := make([]string, len(filter.IDs))
		for i, id := range filter.IDs {
			stringIDs[i] = string(id)
		}
		query = query.Where("id IN ?", stringIDs)
	}

	// 鍚嶇О杩囨护锛堟ā绯婃煡璇級
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	// 绫诲瀷杩囨护
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	// 鐘舵€佽繃婊?
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 浣嶇疆杩囨护
	if filter.MinX != nil {
		query = query.Where("pos_x >= ?", *filter.MinX)
	}
	if filter.MaxX != nil {
		query = query.Where("pos_x <= ?", *filter.MaxX)
	}
	if filter.MinY != nil {
		query = query.Where("pos_y >= ?", *filter.MinY)
	}
	if filter.MaxY != nil {
		query = query.Where("pos_y <= ?", *filter.MaxY)
	}
	if filter.MinZ != nil {
		query = query.Where("pos_z >= ?", *filter.MinZ)
	}
	if filter.MaxZ != nil {
		query = query.Where("pos_z <= ?", *filter.MaxZ)
	}

	// 鏍囩杩囨护
	for key, value := range filter.Labels {
		query = query.Where("JSON_EXTRACT(labels, ?) = ?", "$."+key, value)
	}

	return query
}
