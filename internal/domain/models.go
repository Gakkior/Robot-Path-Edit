// Package domain 瀹氫箟鏍稿績涓氬姟棰嗗煙妯″瀷
//
// 璁捐鍙傝€冿細
// - DDD (Domain-Driven Design) 鐨勫疄浣撹璁?
// - Kubernetes鐨勮祫婧愭ā鍨嬭璁?
// - Grafana鐨勬暟鎹ā鍨嬬粨鏋?
//
// 璁捐鍘熷垯锛?
// 1. 棰嗗煙绾噣锛氫笉渚濊禆澶栭儴妗嗘灦
// 2. 涓嶅彉鎬э細閲嶈瀛楁涓嶅彲鍙?
// 3. 鑱氬悎鏍癸細鏄庣‘鑱氬悎杈圭晫
// 4. 鍊煎璞★細灏佽涓氬姟瑙勫垯
package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// Node 琛ㄧず鍥句腑鐨勪竴涓妭鐐?鐐逛綅
// 杩欐槸涓€涓仛鍚堟牴锛屽寘鍚簡鐐逛綅鐨勬墍鏈変笟鍔￠€昏緫
//
// 璁捐鍙傝€冿細
// - Kubernetes Pod鐨勫厓鏁版嵁缁撴瀯
// - CAD杞欢涓殑鍑犱綍鐐硅〃绀?
type Node struct {
	// 鍩虹鏍囪瘑淇℃伅
	ID     NodeID     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name   string     `json:"name" gorm:"type:varchar(100);not null"`
	Type   NodeType   `json:"type" gorm:"type:varchar(20);not null;default:'point'"`
	Status NodeStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	// 浣嶇疆淇℃伅 - 鏀寔2D鍜?D鍧愭爣
	Position Position `json:"position" gorm:"embedded;embeddedPrefix:pos_"`

	// 鏈哄櫒浜虹浉鍏崇殑6杞村潗鏍囦俊鎭?
	RobotCoords *RobotCoordinates `json:"robot_coords,omitempty" gorm:"embedded;embeddedPrefix:robot_"`

	// 鎵╁睍灞炴€?- 鏀寔鍔ㄦ€佸瓧娈碉紝绫讳技Kubernetes鐨凩abels
	Properties map[string]interface{} `json:"properties,omitempty" gorm:"serializer:json"`

	// 鏍峰紡閰嶇疆
	Style NodeStyle `json:"style" gorm:"embedded;embeddedPrefix:style_"`

	// 鍏冩暟鎹?- 鍙傝€僈ubernetes鐨凮bjectMeta
	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// Path 琛ㄧず涓や釜鑺傜偣涔嬮棿鐨勮矾寰?杩炴帴
// 鑱氬悎鏍癸紝绠＄悊璺緞鐨勫畬鏁寸敓鍛藉懆鏈?
type Path struct {
	// 鍩虹鏍囪瘑淇℃伅
	ID     PathID     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name   string     `json:"name" gorm:"type:varchar(100)"`
	Type   PathType   `json:"type" gorm:"type:varchar(20);not null;default:'direct'"`
	Status PathStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	// 杩炴帴淇℃伅
	StartNodeID NodeID        `json:"start_node_id" gorm:"type:varchar(36);not null;index"`
	EndNodeID   NodeID        `json:"end_node_id" gorm:"type:varchar(36);not null;index"`
	Direction   PathDirection `json:"direction" gorm:"type:varchar(20);not null;default:'bidirectional'"`

	// 璺緞灞炴€?
	Weight   float64 `json:"weight" gorm:"type:decimal(10,2);default:1.0"` // 鏉冮噸/浠ｄ环
	Length   float64 `json:"length,omitempty" gorm:"type:decimal(10,2)"`   // 瀹為檯闀垮害
	MaxSpeed float64 `json:"max_speed,omitempty" gorm:"type:decimal(8,2)"` // 鏈€澶ч€熷害

	// 璺緞鍑犱綍淇℃伅
	Waypoints []Position `json:"waypoints,omitempty" gorm:"serializer:json"`          // 涓棿鐐?
	CurveType CurveType  `json:"curve_type" gorm:"type:varchar(20);default:'linear'"` // 鏇茬嚎绫诲瀷

	// 鎵╁睍灞炴€?
	Properties map[string]interface{} `json:"properties,omitempty" gorm:"serializer:json"`

	// 鏍峰紡閰嶇疆
	Style PathStyle `json:"style" gorm:"embedded;embeddedPrefix:style_"`

	// 鍏冩暟鎹?
	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// DatabaseConnection 琛ㄧず鏁版嵁搴撹繛鎺ラ厤缃?
// 鍊煎璞★紝灏佽鏁版嵁搴撹繛鎺ョ殑鎵€鏈変俊鎭?
type DatabaseConnection struct {
	ID       string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name     string            `json:"name" gorm:"type:varchar(100);not null"`
	Type     string            `json:"type" gorm:"type:varchar(20);not null"` // sqlite, mysql, postgres
	DSN      string            `json:"dsn" gorm:"type:text;not null"`
	Options  map[string]string `json:"options,omitempty" gorm:"serializer:json"`
	Metadata ObjectMeta        `json:"metadata" gorm:"embedded"`
}

// TableMapping 琛ㄧず琛ㄥ瓧娈垫槧灏勯厤缃?
// 鍊煎璞★紝瀹氫箟濡備綍灏嗛€氱敤琛ㄦ槧灏勫埌Node鍜孭ath
type TableMapping struct {
	ID           string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ConnectionID string `json:"connection_id" gorm:"type:varchar(36);not null;index"`
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	Type         string `json:"type" gorm:"type:varchar(20);not null"` // node, path

	// 琛ㄤ俊鎭?
	TableName string `json:"table_name" gorm:"type:varchar(100);not null"`

	// 瀛楁鏄犲皠
	IDField   string `json:"id_field" gorm:"type:varchar(100);not null"`    // 涓婚敭瀛楁
	NameField string `json:"name_field,omitempty" gorm:"type:varchar(100)"` // 鍚嶇О瀛楁

	// 浣嶇疆瀛楁鏄犲皠
	XField string `json:"x_field,omitempty" gorm:"type:varchar(100)"` // X鍧愭爣瀛楁
	YField string `json:"y_field,omitempty" gorm:"type:varchar(100)"` // Y鍧愭爣瀛楁
	ZField string `json:"z_field,omitempty" gorm:"type:varchar(100)"` // Z鍧愭爣瀛楁

	// 璺緞鐗规湁瀛楁
	StartNodeField string `json:"start_node_field,omitempty" gorm:"type:varchar(100)"` // 璧峰鑺傜偣瀛楁
	EndNodeField   string `json:"end_node_field,omitempty" gorm:"type:varchar(100)"`   // 缁撴潫鑺傜偣瀛楁

	// 鎵╁睍瀛楁鏄犲皠
	FieldMappings map[string]string `json:"field_mappings,omitempty" gorm:"serializer:json"`

	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// === 鍊煎璞″畾涔?===

// NodeID 鑺傜偣鍞竴鏍囪瘑绗?
type NodeID string

func NewNodeID() NodeID {
	return NodeID(uuid.New().String())
}

func (id NodeID) String() string {
	return string(id)
}

// PathID 璺緞鍞竴鏍囪瘑绗?
type PathID string

func NewPathID() PathID {
	return PathID(uuid.New().String())
}

func (id PathID) String() string {
	return string(id)
}

// Position 浣嶇疆淇℃伅 - 鍊煎璞?
type Position struct {
	X float64 `json:"x" gorm:"type:decimal(12,6);not null;default:0"`
	Y float64 `json:"y" gorm:"type:decimal(12,6);not null;default:0"`
	Z float64 `json:"z" gorm:"type:decimal(12,6);not null;default:0"`
}

func (p Position) Distance(other Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	dz := p.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// RobotCoordinates 鏈哄櫒浜哄叚杞村潗鏍?- 鍊煎璞?
type RobotCoordinates struct {
	X     float64 `json:"x" gorm:"type:decimal(12,6)"`    // X杞翠綅缃?
	Y     float64 `json:"y" gorm:"type:decimal(12,6)"`    // Y杞翠綅缃?
	Z     float64 `json:"z" gorm:"type:decimal(12,6)"`    // Z杞翠綅缃?
	Roll  float64 `json:"roll" gorm:"type:decimal(8,3)"`  // 缈绘粴瑙?
	Pitch float64 `json:"pitch" gorm:"type:decimal(8,3)"` // 淇话瑙?
	Yaw   float64 `json:"yaw" gorm:"type:decimal(8,3)"`   // 鍋忚埅瑙?
}

// NodeStyle 鑺傜偣鏍峰紡閰嶇疆 - 鍊煎璞?
type NodeStyle struct {
	Shape       string  `json:"shape" gorm:"type:varchar(20);default:'circle'"`         // 褰㈢姸
	Radius      int     `json:"radius" gorm:"type:int;default:20"`                      // 鍗婂緞
	Color       string  `json:"color" gorm:"type:varchar(20);default:'#3498db'"`        // 棰滆壊
	BorderColor string  `json:"border_color" gorm:"type:varchar(20);default:'#2980b9'"` // 杈规棰滆壊
	BorderWidth int     `json:"border_width" gorm:"type:int;default:2"`                 // 杈规瀹藉害
	Opacity     float64 `json:"opacity" gorm:"type:decimal(3,2);default:1.0"`           // 閫忔槑搴?
}

// PathStyle 璺緞鏍峰紡閰嶇疆 - 鍊煎璞?
type PathStyle struct {
	LineType  string  `json:"line_type" gorm:"type:varchar(20);default:'solid'"` // 绾垮瀷
	Width     int     `json:"width" gorm:"type:int;default:2"`                   // 绾垮
	Color     string  `json:"color" gorm:"type:varchar(20);default:'#34495e'"`   // 棰滆壊
	ArrowSize int     `json:"arrow_size" gorm:"type:int;default:8"`              // 绠ご澶у皬
	Opacity   float64 `json:"opacity" gorm:"type:decimal(3,2);default:1.0"`      // 閫忔槑搴?
}

// ObjectMeta 瀵硅薄鍏冩暟鎹?- 鍙傝€僈ubernetes ObjectMeta
type ObjectMeta struct {
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	Version     int               `json:"version" gorm:"type:int;default:1"`            // 涔愯閿佺増鏈?
	Labels      map[string]string `json:"labels,omitempty" gorm:"serializer:json"`      // 鏍囩
	Annotations map[string]string `json:"annotations,omitempty" gorm:"serializer:json"` // 娉ㄨВ
}

// === 鏋氫妇绫诲瀷瀹氫箟 ===

type NodeType string

const (
	NodeTypePoint    NodeType = "point"    // 鏅€氱偣浣?
	NodeTypeWaypoint NodeType = "waypoint" // 璺緞鐐?
	NodeTypeStation  NodeType = "station"  // 宸ヤ綔绔?
	NodeTypeCharging NodeType = "charging" // 鍏呯數妗?
)

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"   // 婵€娲?
	NodeStatusInactive NodeStatus = "inactive" // 闈炴縺娲?
	NodeStatusDeleted  NodeStatus = "deleted"  // 宸插垹闄?
)

type PathType string

const (
	PathTypeDirect PathType = "direct" // 鐩寸嚎璺緞
	PathTypeCurved PathType = "curved" // 鏇茬嚎璺緞
	PathTypeSpline PathType = "spline" // 鏍锋潯鏇茬嚎
)

type PathStatus string

const (
	PathStatusActive   PathStatus = "active"   // 婵€娲?
	PathStatusInactive PathStatus = "inactive" // 闈炴縺娲?
	PathStatusBlocked  PathStatus = "blocked"  // 闃诲
	PathStatusDeleted  PathStatus = "deleted"  // 宸插垹闄?
)

type PathDirection string

const (
	PathDirectionUnidirectional PathDirection = "unidirectional" // 鍗曞悜
	PathDirectionBidirectional  PathDirection = "bidirectional"  // 鍙屽悜
)

type CurveType string

const (
	CurveTypeLinear CurveType = "linear" // 绾挎€?
	CurveTypeBezier CurveType = "bezier" // 璐濆灏旀洸绾?
	CurveTypeSpline CurveType = "spline" // 鏍锋潯鏇茬嚎
)

// === 涓氬姟鏂规硶 ===

// NewNode 鍒涘缓鏂拌妭鐐?
func NewNode(name string, position Position) *Node {
	return &Node{
		ID:         NewNodeID(),
		Name:       name,
		Type:       NodeTypePoint,
		Status:     NodeStatusActive,
		Position:   position,
		Style:      NodeStyle{},
		Properties: make(map[string]interface{}),
		Metadata: ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
	}
}

// NewPath 鍒涘缓鏂拌矾寰?
func NewPath(name string, startNodeID, endNodeID NodeID) *Path {
	return &Path{
		ID:          NewPathID(),
		Name:        name,
		Type:        PathTypeDirect,
		Status:      PathStatusActive,
		StartNodeID: startNodeID,
		EndNodeID:   endNodeID,
		Direction:   PathDirectionBidirectional,
		Weight:      1.0,
		CurveType:   CurveTypeLinear,
		Style:       PathStyle{},
		Properties:  make(map[string]interface{}),
		Metadata: ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
	}
}

// IsValid 楠岃瘉鑺傜偣鏄惁鏈夋晥
func (n *Node) IsValid() error {
	if n.ID == "" {
		return fmt.Errorf("鑺傜偣ID涓嶈兘涓虹┖")
	}
	if n.Name == "" {
		return fmt.Errorf("鑺傜偣鍚嶇О涓嶈兘涓虹┖")
	}
	return nil
}

// IsValid 楠岃瘉璺緞鏄惁鏈夋晥
func (p *Path) IsValid() error {
	if p.ID == "" {
		return fmt.Errorf("璺緞ID涓嶈兘涓虹┖")
	}
	if p.StartNodeID == "" || p.EndNodeID == "" {
		return fmt.Errorf("璺緞鐨勮捣濮嬭妭鐐瑰拰缁撴潫鑺傜偣涓嶈兘涓虹┖")
	}
	if p.StartNodeID == p.EndNodeID {
		return fmt.Errorf("璺緞鐨勮捣濮嬭妭鐐瑰拰缁撴潫鑺傜偣涓嶈兘鐩稿悓")
	}
	if p.Weight < 0 {
		return fmt.Errorf("璺緞鏉冮噸涓嶈兘涓鸿礋鏁?)
	}
	return nil
}
