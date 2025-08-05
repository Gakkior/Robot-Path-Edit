// Package services 瀹炵幇涓氬姟閫昏緫灞?
//
// 璁捐鍙傝€冿細
// - DDD鐨勫簲鐢ㄦ湇鍔℃ā寮?
// - Kubernetes鐨勬帶鍒跺櫒妯″紡
// - 寰湇鍔＄殑涓氬姟閫昏緫灏佽
//
// 鐗圭偣锛?
// 1. 涓氬姟瑙勫垯灏佽锛氬寘鍚墍鏈変笟鍔￠€昏緫
// 2. 浜嬪姟绠＄悊锛氱‘淇濇暟鎹竴鑷存€?
// 3. 浜嬩欢鍙戝竷锛氭敮鎸佷簨浠堕┍鍔ㄦ灦鏋?
// 4. 楠岃瘉鍜屾巿鏉冿細缁熶竴鐨勪笟鍔￠獙璇?
package services

import (
	"context"
	"fmt"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// NodeService 鑺傜偣鏈嶅姟鎺ュ彛
type NodeService interface {
	// 鍩虹鎿嶄綔
	CreateNode(ctx context.Context, req CreateNodeRequest) (*domain.Node, error)
	GetNode(ctx context.Context, id domain.NodeID) (*domain.Node, error)
	UpdateNode(ctx context.Context, req UpdateNodeRequest) (*domain.Node, error)
	DeleteNode(ctx context.Context, id domain.NodeID) error

	// 鎵归噺鎿嶄綔
	CreateNodes(ctx context.Context, req CreateNodesRequest) ([]*domain.Node, error)
	GetNodes(ctx context.Context, req GetNodesRequest) (*GetNodesResponse, error)
	ListNodes(ctx context.Context) ([]*domain.Node, error)

	// 浣嶇疆鎿嶄綔
	UpdateNodePosition(ctx context.Context, id domain.NodeID, position domain.Position) error
	MoveNodes(ctx context.Context, moves []NodeMove) error

	// 鏌ヨ鎿嶄綔
	SearchNodes(ctx context.Context, req SearchNodesRequest) ([]*domain.Node, error)
	GetNodesInArea(ctx context.Context, req GetNodesInAreaRequest) ([]*domain.Node, error)
	GetNearbyNodes(ctx context.Context, req GetNearbyNodesRequest) ([]*domain.Node, error)

	// 鍒嗘瀽鎿嶄綔
	GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error)
	GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error)
	ValidateNode(ctx context.Context, node *domain.Node) error
}

// 璇锋眰鍜屽搷搴旂粨鏋勪綋瀹氫箟

// CreateNodeRequest 鍒涘缓鑺傜偣璇锋眰
type CreateNodeRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Type        domain.NodeType          `json:"type"`
	Position    domain.Position          `json:"position"`
	RobotCoords *domain.RobotCoordinates `json:"robot_coords,omitempty"`
	Properties  map[string]interface{}   `json:"properties,omitempty"`
	Style       domain.NodeStyle         `json:"style"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
}

// UpdateNodeRequest 鏇存柊鑺傜偣璇锋眰
type UpdateNodeRequest struct {
	ID          domain.NodeID            `json:"id" binding:"required"`
	Name        *string                  `json:"name,omitempty"`
	Type        *domain.NodeType         `json:"type,omitempty"`
	Status      *domain.NodeStatus       `json:"status,omitempty"`
	Position    *domain.Position         `json:"position,omitempty"`
	RobotCoords *domain.RobotCoordinates `json:"robot_coords,omitempty"`
	Properties  map[string]interface{}   `json:"properties,omitempty"`
	Style       *domain.NodeStyle        `json:"style,omitempty"`
	Labels      map[string]string        `json:"labels,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
}

// CreateNodesRequest 鎵归噺鍒涘缓鑺傜偣璇锋眰
type CreateNodesRequest struct {
	Nodes []CreateNodeRequest `json:"nodes" binding:"required,dive"`
}

// GetNodesRequest 鑾峰彇鑺傜偣鍒楄〃璇锋眰
type GetNodesRequest struct {
	Filter   repositories.NodeFilter `json:"filter"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	OrderBy  string                  `json:"order_by"`
	Order    string                  `json:"order"`
}

// GetNodesResponse 鑾峰彇鑺傜偣鍒楄〃鍝嶅簲
type GetNodesResponse struct {
	Nodes      []*domain.Node `json:"nodes"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// SearchNodesRequest 鎼滅储鑺傜偣璇锋眰
type SearchNodesRequest struct {
	Query  string            `json:"query"`
	Type   domain.NodeType   `json:"type,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Limit  int               `json:"limit"`
}

// GetNodesInAreaRequest 鑾峰彇鍖哄煙鍐呰妭鐐硅姹?
type GetNodesInAreaRequest struct {
	MinX float64 `json:"min_x" binding:"required"`
	MinY float64 `json:"min_y" binding:"required"`
	MaxX float64 `json:"max_x" binding:"required"`
	MaxY float64 `json:"max_y" binding:"required"`
}

// GetNearbyNodesRequest 鑾峰彇闄勮繎鑺傜偣璇锋眰
type GetNearbyNodesRequest struct {
	Position domain.Position `json:"position" binding:"required"`
	Radius   float64         `json:"radius" binding:"required,gt=0"`
	Limit    int             `json:"limit"`
}

// NodeMove 鑺傜偣绉诲姩璇锋眰
type NodeMove struct {
	NodeID      domain.NodeID   `json:"node_id" binding:"required"`
	NewPosition domain.Position `json:"new_position" binding:"required"`
}

// nodeService 鑺傜偣鏈嶅姟瀹炵幇
type nodeService struct {
	nodeRepo repositories.NodeRepository
}

// NewNodeService 鍒涘缓鑺傜偣鏈嶅姟瀹炰緥
func NewNodeService(nodeRepo repositories.NodeRepository) NodeService {
	return &nodeService{
		nodeRepo: nodeRepo,
	}
}

// CreateNode 鍒涘缓鑺傜偣
func (s *nodeService) CreateNode(ctx context.Context, req CreateNodeRequest) (*domain.Node, error) {
	// 1. 鍒涘缓鑺傜偣瀹炰綋
	node := domain.NewNode(req.Name, req.Position)

	// 2. 璁剧疆鍙€夊睘鎬?
	if req.Type != "" {
		node.Type = req.Type
	}

	if req.RobotCoords != nil {
		node.RobotCoords = req.RobotCoords
	}

	if req.Properties != nil {
		node.Properties = req.Properties
	}

	if req.Style.Shape != "" {
		node.Style = req.Style
	}

	if req.Labels != nil {
		node.Metadata.Labels = req.Labels
	}

	if req.Annotations != nil {
		node.Metadata.Annotations = req.Annotations
	}

	// 3. 涓氬姟瑙勫垯楠岃瘉
	if err := s.ValidateNode(ctx, node); err != nil {
		return nil, fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
	}

	// 4. 鎸佷箙鍖?
	if err := s.nodeRepo.Create(ctx, node); err != nil {
		return nil, fmt.Errorf("鍒涘缓鑺傜偣澶辫触: %w", err)
	}

	return node, nil
}

// GetNode 鑾峰彇鑺傜偣
func (s *nodeService) GetNode(ctx context.Context, id domain.NodeID) (*domain.Node, error) {
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣澶辫触: %w", err)
	}

	return node, nil
}

// UpdateNode 鏇存柊鑺傜偣
func (s *nodeService) UpdateNode(ctx context.Context, req UpdateNodeRequest) (*domain.Node, error) {
	// 1. 鑾峰彇鐜版湁鑺傜偣
	node, err := s.nodeRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("鑺傜偣涓嶅瓨鍦? %w", err)
	}

	// 2. 搴旂敤鏇存柊
	if req.Name != nil {
		node.Name = *req.Name
	}

	if req.Type != nil {
		node.Type = *req.Type
	}

	if req.Status != nil {
		node.Status = *req.Status
	}

	if req.Position != nil {
		node.Position = *req.Position
	}

	if req.RobotCoords != nil {
		node.RobotCoords = req.RobotCoords
	}

	if req.Properties != nil {
		node.Properties = req.Properties
	}

	if req.Style != nil {
		node.Style = *req.Style
	}

	if req.Labels != nil {
		node.Metadata.Labels = req.Labels
	}

	if req.Annotations != nil {
		node.Metadata.Annotations = req.Annotations
	}

	// 3. 楠岃瘉鏇存柊鍚庣殑鑺傜偣
	if err := s.ValidateNode(ctx, node); err != nil {
		return nil, fmt.Errorf("鑺傜偣楠岃瘉澶辫触: %w", err)
	}

	// 4. 鎸佷箙鍖栨洿鏂?
	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return nil, fmt.Errorf("鏇存柊鑺傜偣澶辫触: %w", err)
	}

	return node, nil
}

// DeleteNode 鍒犻櫎鑺傜偣
func (s *nodeService) DeleteNode(ctx context.Context, id domain.NodeID) error {
	// 1. 妫€鏌ヨ妭鐐规槸鍚﹀瓨鍦?
	_, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %w", err)
	}

	// 2. 妫€鏌ユ槸鍚︽湁璺緞杩炴帴 (涓氬姟瑙勫垯: 涓嶈兘鍒犻櫎鏈夎繛鎺ョ殑鑺傜偣)
	connectedNodes, err := s.nodeRepo.GetConnectedNodes(ctx, id)
	if err != nil {
		return fmt.Errorf("妫€鏌ヨ妭鐐硅繛鎺ュけ璐? %w", err)
	}

	if len(connectedNodes) > 0 {
		return fmt.Errorf("涓嶈兘鍒犻櫎鏈夎矾寰勮繛鎺ョ殑鑺傜偣锛岃鍏堝垹闄ょ浉鍏宠矾寰?)
	}

	// 3. 鎵ц鍒犻櫎
	if err := s.nodeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("鍒犻櫎鑺傜偣澶辫触: %w", err)
	}

	return nil
}

// CreateNodes 鎵归噺鍒涘缓鑺傜偣
func (s *nodeService) CreateNodes(ctx context.Context, req CreateNodesRequest) ([]*domain.Node, error) {
	nodes := make([]*domain.Node, 0, len(req.Nodes))

	// 1. 鍒涘缓鎵€鏈夎妭鐐瑰疄浣?
	for _, nodeReq := range req.Nodes {
		node := domain.NewNode(nodeReq.Name, nodeReq.Position)

		// 璁剧疆灞炴€?
		if nodeReq.Type != "" {
			node.Type = nodeReq.Type
		}
		if nodeReq.RobotCoords != nil {
			node.RobotCoords = nodeReq.RobotCoords
		}
		if nodeReq.Properties != nil {
			node.Properties = nodeReq.Properties
		}
		if nodeReq.Style.Shape != "" {
			node.Style = nodeReq.Style
		}
		if nodeReq.Labels != nil {
			node.Metadata.Labels = nodeReq.Labels
		}
		if nodeReq.Annotations != nil {
			node.Metadata.Annotations = nodeReq.Annotations
		}

		// 楠岃瘉鑺傜偣
		if err := s.ValidateNode(ctx, node); err != nil {
			return nil, fmt.Errorf("鑺傜偣 %s 楠岃瘉澶辫触: %w", node.Name, err)
		}

		nodes = append(nodes, node)
	}

	// 2. 鎵归噺鍒涘缓
	if err := s.nodeRepo.CreateBatch(ctx, nodes); err != nil {
		return nil, fmt.Errorf("鎵归噺鍒涘缓鑺傜偣澶辫触: %w", err)
	}

	return nodes, nil
}

// GetNodes 鑾峰彇鑺傜偣鍒楄〃
func (s *nodeService) GetNodes(ctx context.Context, req GetNodesRequest) (*GetNodesResponse, error) {
	// 1. 璁剧疆榛樿鍒嗛〉鍙傛暟
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	// 2. 鏋勫缓鏌ヨ閫夐」
	options := repositories.ListOptions{
		Filter:   req.Filter,
		Page:     req.Page,
		PageSize: req.PageSize,
		OrderBy:  req.OrderBy,
		Order:    req.Order,
	}

	// 3. 鏌ヨ鑺傜偣鍜屾€绘暟
	nodes, err := s.nodeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("鏌ヨ鑺傜偣鍒楄〃澶辫触: %w", err)
	}

	total, err := s.nodeRepo.Count(ctx, req.Filter)
	if err != nil {
		return nil, fmt.Errorf("缁熻鑺傜偣鏁伴噺澶辫触: %w", err)
	}

	// 4. 璁＄畻鎬婚〉鏁?
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &GetNodesResponse{
		Nodes:      nodes,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListNodes 鑾峰彇鎵€鏈夎妭鐐瑰垪琛?
func (s *nodeService) ListNodes(ctx context.Context) ([]*domain.Node, error) {
	// 鏋勫缓鏌ヨ閫夐」锛屼笉鍒嗛〉
	options := repositories.ListOptions{
		PageSize: 0, // 0 琛ㄧず涓嶅垎椤?
	}

	nodes, err := s.nodeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鑺傜偣鍒楄〃澶辫触: %w", err)
	}

	return nodes, nil
}

// UpdateNodePosition 鏇存柊鑺傜偣浣嶇疆
func (s *nodeService) UpdateNodePosition(ctx context.Context, id domain.NodeID, position domain.Position) error {
	// 鑾峰彇鑺傜偣
	node, err := s.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("鑺傜偣涓嶅瓨鍦? %w", err)
	}

	// 鏇存柊浣嶇疆
	node.Position = position

	// 淇濆瓨鏇存柊
	if err := s.nodeRepo.Update(ctx, node); err != nil {
		return fmt.Errorf("鏇存柊鑺傜偣浣嶇疆澶辫触: %w", err)
	}

	return nil
}

// MoveNodes 鎵归噺绉诲姩鑺傜偣
func (s *nodeService) MoveNodes(ctx context.Context, moves []NodeMove) error {
	// 鑾峰彇鎵€鏈夐渶瑕佺Щ鍔ㄧ殑鑺傜偣
	nodeIDs := make([]domain.NodeID, len(moves))
	for i, move := range moves {
		nodeIDs[i] = move.NodeID
	}

	nodes, err := s.nodeRepo.GetByIDs(ctx, nodeIDs)
	if err != nil {
		return fmt.Errorf("鑾峰彇鑺傜偣澶辫触: %w", err)
	}

	// 鍒涘缓鑺傜偣ID鍒颁綅缃殑鏄犲皠
	nodePositions := make(map[domain.NodeID]domain.Position)
	for _, move := range moves {
		nodePositions[move.NodeID] = move.NewPosition
	}

	// 鏇存柊鑺傜偣浣嶇疆
	for _, node := range nodes {
		if newPos, exists := nodePositions[node.ID]; exists {
			node.Position = newPos
		}
	}

	// 鎵归噺鏇存柊
	if err := s.nodeRepo.UpdateBatch(ctx, nodes); err != nil {
		return fmt.Errorf("鎵归噺鏇存柊鑺傜偣浣嶇疆澶辫触: %w", err)
	}

	return nil
}

// SearchNodes 鎼滅储鑺傜偣
func (s *nodeService) SearchNodes(ctx context.Context, req SearchNodesRequest) ([]*domain.Node, error) {
	filter := repositories.NodeFilter{
		Name:   req.Query,
		Type:   req.Type,
		Labels: req.Labels,
	}

	options := repositories.ListOptions{
		Filter:   filter,
		PageSize: req.Limit,
	}

	nodes, err := s.nodeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("鎼滅储鑺傜偣澶辫触: %w", err)
	}

	return nodes, nil
}

// GetNodesInArea 鑾峰彇鍖哄煙鍐呯殑鑺傜偣
func (s *nodeService) GetNodesInArea(ctx context.Context, req GetNodesInAreaRequest) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetByArea(ctx, req.MinX, req.MinY, req.MaxX, req.MaxY)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇鍖哄煙鍐呰妭鐐瑰け璐? %w", err)
	}

	return nodes, nil
}

// GetNearbyNodes 鑾峰彇闄勮繎鐨勮妭鐐?
func (s *nodeService) GetNearbyNodes(ctx context.Context, req GetNearbyNodesRequest) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetNearby(ctx, req.Position, req.Radius)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇闄勮繎鑺傜偣澶辫触: %w", err)
	}

	// 濡傛灉璁剧疆浜嗛檺鍒讹紝鎴彇缁撴灉
	if req.Limit > 0 && len(nodes) > req.Limit {
		nodes = nodes[:req.Limit]
	}

	return nodes, nil
}

// GetConnectedNodes 鑾峰彇杩炴帴鐨勮妭鐐?
func (s *nodeService) GetConnectedNodes(ctx context.Context, nodeID domain.NodeID) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetConnectedNodes(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇杩炴帴鑺傜偣澶辫触: %w", err)
	}

	return nodes, nil
}

// GetIsolatedNodes 鑾峰彇瀛ょ珛鑺傜偣
func (s *nodeService) GetIsolatedNodes(ctx context.Context) ([]*domain.Node, error) {
	nodes, err := s.nodeRepo.GetIsolatedNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("鑾峰彇瀛ょ珛鑺傜偣澶辫触: %w", err)
	}

	return nodes, nil
}

// ValidateNode 楠岃瘉鑺傜偣
func (s *nodeService) ValidateNode(ctx context.Context, node *domain.Node) error {
	// 1. 鍩虹楠岃瘉
	if err := node.IsValid(); err != nil {
		return err
	}

	// 2. 涓氬姟瑙勫垯楠岃瘉

	// 鍚嶇О涓嶈兘閲嶅锛堝彲閫夌殑涓氬姟瑙勫垯锛?
	// 杩欓噷鍙互娣诲姞鏇村涓氬姟楠岃瘉閫昏緫

	// 3. 浣嶇疆鍚堢悊鎬ч獙璇?
	// 鍙互楠岃瘉浣嶇疆鏄惁鍦ㄥ厑璁哥殑鑼冨洿鍐?

	return nil
}
