// Package domain 定义核心业务领域模型
//
// 设计参考：
// - DDD (Domain-Driven Design) 的实体设计
// - Kubernetes的资源模型设计
// - Grafana的数据模型结构
//
// 设计原则：
// 1. 领域纯净：不依赖外部框架
// 2. 不变性：重要字段不可变
// 3. 聚合根：明确聚合边界
// 4. 值对象：封装业务规则
package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// Node 表示图中的一个节�?点位
// 这是一个聚合根，包含了点位的所有业务逻辑
//
// 设计参考：
// - Kubernetes Pod的元数据结构
// - CAD软件中的几何点表�?
type Node struct {
	// 基础标识信息
	ID     NodeID     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name   string     `json:"name" gorm:"type:varchar(100);not null"`
	Type   NodeType   `json:"type" gorm:"type:varchar(20);not null;default:'point'"`
	Status NodeStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	// 位置信息 - 支持2D�?D坐标
	Position Position `json:"position" gorm:"embedded;embeddedPrefix:pos_"`

	// 机器人相关的6轴坐标信�?
	RobotCoords *RobotCoordinates `json:"robot_coords,omitempty" gorm:"embedded;embeddedPrefix:robot_"`

	// 扩展属�?- 支持动态字段，类似Kubernetes的Labels
	Properties map[string]interface{} `json:"properties,omitempty" gorm:"serializer:json"`

	// 样式配置
	Style NodeStyle `json:"style" gorm:"embedded;embeddedPrefix:style_"`

	// 元数�?- 参考Kubernetes的ObjectMeta
	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// Path 表示两个节点之间的路�?连接
// 聚合根，管理路径的完整生命周�?
type Path struct {
	// 基础标识信息
	ID     PathID     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name   string     `json:"name" gorm:"type:varchar(100)"`
	Type   PathType   `json:"type" gorm:"type:varchar(20);not null;default:'direct'"`
	Status PathStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	// 连接信息
	StartNodeID NodeID        `json:"start_node_id" gorm:"type:varchar(36);not null;index"`
	EndNodeID   NodeID        `json:"end_node_id" gorm:"type:varchar(36);not null;index"`
	Direction   PathDirection `json:"direction" gorm:"type:varchar(20);not null;default:'bidirectional'"`

	// 路径属�?
	Weight   float64 `json:"weight" gorm:"type:decimal(10,2);default:1.0"` // 权重/代价
	Length   float64 `json:"length,omitempty" gorm:"type:decimal(10,2)"`   // 实际长度
	MaxSpeed float64 `json:"max_speed,omitempty" gorm:"type:decimal(8,2)"` // 最大速度

	// 路径几何信息
	Waypoints []Position `json:"waypoints,omitempty" gorm:"serializer:json"`          // 中间�?
	CurveType CurveType  `json:"curve_type" gorm:"type:varchar(20);default:'linear'"` // 曲线类型

	// 扩展属�?
	Properties map[string]interface{} `json:"properties,omitempty" gorm:"serializer:json"`

	// 样式配置
	Style PathStyle `json:"style" gorm:"embedded;embeddedPrefix:style_"`

	// 元数�?
	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// DatabaseConnection 表示数据库连接配�?
// 值对象，封装数据库连接的所有信�?
type DatabaseConnection struct {
	ID       string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name     string            `json:"name" gorm:"type:varchar(100);not null"`
	Type     string            `json:"type" gorm:"type:varchar(20);not null"` // sqlite, mysql, postgres
	DSN      string            `json:"dsn" gorm:"type:text;not null"`
	Options  map[string]string `json:"options,omitempty" gorm:"serializer:json"`
	Metadata ObjectMeta        `json:"metadata" gorm:"embedded"`
}

// TableMapping 表示表字段映射配�?
// 值对象，定义如何将通用表映射到Node和Path
type TableMapping struct {
	ID           string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ConnectionID string `json:"connection_id" gorm:"type:varchar(36);not null;index"`
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	Type         string `json:"type" gorm:"type:varchar(20);not null"` // node, path

	// 表信�?
	TableName string `json:"table_name" gorm:"type:varchar(100);not null"`

	// 字段映射
	IDField   string `json:"id_field" gorm:"type:varchar(100);not null"`    // 主键字段
	NameField string `json:"name_field,omitempty" gorm:"type:varchar(100)"` // 名称字段

	// 位置字段映射
	XField string `json:"x_field,omitempty" gorm:"type:varchar(100)"` // X坐标字段
	YField string `json:"y_field,omitempty" gorm:"type:varchar(100)"` // Y坐标字段
	ZField string `json:"z_field,omitempty" gorm:"type:varchar(100)"` // Z坐标字段

	// 路径特有字段
	StartNodeField string `json:"start_node_field,omitempty" gorm:"type:varchar(100)"` // 起始节点字段
	EndNodeField   string `json:"end_node_field,omitempty" gorm:"type:varchar(100)"`   // 结束节点字段

	// 扩展字段映射
	FieldMappings map[string]string `json:"field_mappings,omitempty" gorm:"serializer:json"`

	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// === 值对象定�?===

// NodeID 节点唯一标识�?
type NodeID string

func NewNodeID() NodeID {
	return NodeID(uuid.New().String())
}

func (id NodeID) String() string {
	return string(id)
}

// PathID 路径唯一标识�?
type PathID string

func NewPathID() PathID {
	return PathID(uuid.New().String())
}

func (id PathID) String() string {
	return string(id)
}

// Position 位置信息 - 值对�?
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

// RobotCoordinates 机器人六轴坐�?- 值对�?
type RobotCoordinates struct {
	X     float64 `json:"x" gorm:"type:decimal(12,6)"`    // X轴位�?
	Y     float64 `json:"y" gorm:"type:decimal(12,6)"`    // Y轴位�?
	Z     float64 `json:"z" gorm:"type:decimal(12,6)"`    // Z轴位�?
	Roll  float64 `json:"roll" gorm:"type:decimal(8,3)"`  // 翻滚�?
	Pitch float64 `json:"pitch" gorm:"type:decimal(8,3)"` // 俯仰�?
	Yaw   float64 `json:"yaw" gorm:"type:decimal(8,3)"`   // 偏航�?
}

// NodeStyle 节点样式配置 - 值对�?
type NodeStyle struct {
	Shape       string  `json:"shape" gorm:"type:varchar(20);default:'circle'"`         // 形状
	Radius      int     `json:"radius" gorm:"type:int;default:20"`                      // 半径
	Color       string  `json:"color" gorm:"type:varchar(20);default:'#3498db'"`        // 颜色
	BorderColor string  `json:"border_color" gorm:"type:varchar(20);default:'#2980b9'"` // 边框颜色
	BorderWidth int     `json:"border_width" gorm:"type:int;default:2"`                 // 边框宽度
	Opacity     float64 `json:"opacity" gorm:"type:decimal(3,2);default:1.0"`           // 透明�?
}

// PathStyle 路径样式配置 - 值对�?
type PathStyle struct {
	LineType  string  `json:"line_type" gorm:"type:varchar(20);default:'solid'"` // 线型
	Width     int     `json:"width" gorm:"type:int;default:2"`                   // 线宽
	Color     string  `json:"color" gorm:"type:varchar(20);default:'#34495e'"`   // 颜色
	ArrowSize int     `json:"arrow_size" gorm:"type:int;default:8"`              // 箭头大小
	Opacity   float64 `json:"opacity" gorm:"type:decimal(3,2);default:1.0"`      // 透明�?
}

// ObjectMeta 对象元数�?- 参考Kubernetes ObjectMeta
type ObjectMeta struct {
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	Version     int               `json:"version" gorm:"type:int;default:1"`            // 乐观锁版�?
	Labels      map[string]string `json:"labels,omitempty" gorm:"serializer:json"`      // 标签
	Annotations map[string]string `json:"annotations,omitempty" gorm:"serializer:json"` // 注解
}

// === 枚举类型定义 ===

type NodeType string

const (
	NodeTypePoint    NodeType = "point"    // 普通点�?
	NodeTypeWaypoint NodeType = "waypoint" // 路径�?
	NodeTypeStation  NodeType = "station"  // 工作�?
	NodeTypeCharging NodeType = "charging" // 充电�?
)

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"   // 激�?
	NodeStatusInactive NodeStatus = "inactive" // 非激�?
	NodeStatusDeleted  NodeStatus = "deleted"  // 已删�?
)

type PathType string

const (
	PathTypeDirect PathType = "direct" // 直线路径
	PathTypeCurved PathType = "curved" // 曲线路径
	PathTypeSpline PathType = "spline" // 样条曲线
)

type PathStatus string

const (
	PathStatusActive   PathStatus = "active"   // 激�?
	PathStatusInactive PathStatus = "inactive" // 非激�?
	PathStatusBlocked  PathStatus = "blocked"  // 阻塞
	PathStatusDeleted  PathStatus = "deleted"  // 已删�?
)

type PathDirection string

const (
	PathDirectionUnidirectional PathDirection = "unidirectional" // 单向
	PathDirectionBidirectional  PathDirection = "bidirectional"  // 双向
)

type CurveType string

const (
	CurveTypeLinear CurveType = "linear" // 线�?
	CurveTypeBezier CurveType = "bezier" // 贝塞尔曲�?
	CurveTypeSpline CurveType = "spline" // 样条曲线
)

// === 业务方法 ===

// NewNode 创建新节�?
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

// NewPath 创建新路�?
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

// IsValid 验证节点是否有效
func (n *Node) IsValid() error {
	if n.ID == "" {
		return fmt.Errorf("节点ID不能为空")
	}
	if n.Name == "" {
		return fmt.Errorf("节点名称不能为空")
	}
	return nil
}

// IsValid 验证路径是否有效
func (p *Path) IsValid() error {
	if p.ID == "" {
		return fmt.Errorf("路径ID不能为空")
	}
	if p.StartNodeID == "" || p.EndNodeID == "" {
		return fmt.Errorf("路径的起始节点和结束节点不能为空")
	}
	if p.StartNodeID == p.EndNodeID {
		return fmt.Errorf("路径的起始节点和结束节点不能相同")
	}
	if p.Weight < 0 {
		return fmt.Errorf("路径权重不能为负�?)
	}
	return nil
}
