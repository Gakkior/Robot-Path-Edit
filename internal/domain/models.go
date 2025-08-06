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

// Node 表示图中的一个节点点位
// 这是一个聚合根，包含了点位的所有业务逻辑
//
// 设计参考：
// - Kubernetes Pod的元数据结构
// - CAD软件中的几何点表示
type Node struct {
	// 基础标识信息
	ID     NodeID     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name   string     `json:"name" gorm:"type:varchar(100);not null"`
	Type   NodeType   `json:"type" gorm:"type:varchar(20);not null;default:'point'"`
	Status NodeStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	// 位置信息 - 支持2D/3D坐标
	Position Position `json:"position" gorm:"embedded;embeddedPrefix:pos_"`

	// 机器人相关的6轴坐标信息
	RobotCoords *RobotCoordinates `json:"robot_coords,omitempty" gorm:"embedded;embeddedPrefix:robot_"`

	// 扩展属性 - 支持动态字段，类似Kubernetes的Labels
	Properties map[string]interface{} `json:"properties,omitempty" gorm:"serializer:json"`

	// 样式配置
	Style NodeStyle `json:"style" gorm:"embedded;embeddedPrefix:style_"`

	// 元数据 - 参考Kubernetes的ObjectMeta
	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// Path 表示两个节点之间的路径连接
// 聚合根，管理路径的完整生命周期
//
// 设计参考：
// - 图论中的边（Edge）概念
// - 导航系统中的路径规划
// - 网络拓扑中的连接关系
type Path struct {
	// 基础标识信息
	ID     PathID     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name   string     `json:"name" gorm:"type:varchar(100);not null"`
	Type   PathType   `json:"type" gorm:"type:varchar(20);not null;default:'normal'"`
	Status PathStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	// 连接关系
	StartNodeID NodeID `json:"start_node_id" gorm:"type:varchar(36);not null;index"`
	EndNodeID   NodeID `json:"end_node_id" gorm:"type:varchar(36);not null;index"`

	// 路径属性
	Weight    float64   `json:"weight" gorm:"type:decimal(12,6);default:0"`
	Length    float64   `json:"length,omitempty" gorm:"type:decimal(12,6)"`
	Direction string    `json:"direction,omitempty" gorm:"type:varchar(20);default:'bidirectional'"`
	CurveType CurveType `json:"curve_type" gorm:"type:varchar(20);default:'linear'"`

	// 路径关键点
	Waypoints []Position `json:"waypoints,omitempty" gorm:"serializer:json"` // 中间点

	// 样式配置
	Style PathStyle `json:"style" gorm:"embedded;embeddedPrefix:style_"`

	// 扩展属性
	Properties map[string]interface{} `json:"properties,omitempty" gorm:"serializer:json"`

	// 元数据
	Metadata ObjectMeta `json:"metadata" gorm:"embedded"`
}

// DatabaseConnection 表示数据库连接配置
// 值对象，封装数据库连接的所有信息
type DatabaseConnection struct {
	ID         string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name       string            `json:"name" gorm:"type:varchar(100);not null"`
	Type       string            `json:"type" gorm:"type:varchar(20);not null"`
	Host       string            `json:"host" gorm:"type:varchar(255);not null"`
	Port       int               `json:"port" gorm:"type:int;not null"`
	Database   string            `json:"database" gorm:"type:varchar(100);not null"`
	Username   string            `json:"username" gorm:"type:varchar(100);not null"`
	Password   string            `json:"password" gorm:"type:varchar(255);not null"`
	Properties map[string]string `json:"properties,omitempty" gorm:"serializer:json"`
}

// TableMapping 表示表字段映射配置
type TableMapping struct {
	ID           string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ConnectionID string            `json:"connection_id" gorm:"type:varchar(36);not null"`
	TableName    string            `json:"table_name" gorm:"type:varchar(100);not null"`
	NodeMapping  *NodeTableMapping `json:"node_mapping,omitempty" gorm:"serializer:json"`
	PathMapping  *PathTableMapping `json:"path_mapping,omitempty" gorm:"serializer:json"`
}

// NodeTableMapping 节点表映射
type NodeTableMapping struct {
	IDField   string `json:"id_field"`
	NameField string `json:"name_field"`
	TypeField string `json:"type_field"`
	XField    string `json:"x_field"`
	YField    string `json:"y_field"`
	ZField    string `json:"z_field"`
}

// PathTableMapping 路径表映射
type PathTableMapping struct {
	IDField        string `json:"id_field"`
	NameField      string `json:"name_field"`
	StartNodeField string `json:"start_node_field"`
	EndNodeField   string `json:"end_node_field"`
	WeightField    string `json:"weight_field"`
}

// === 值对象定义 ===

// NodeID 节点唯一标识符
type NodeID string

// NewNodeID 创建新的节点ID
func NewNodeID() NodeID {
	return NodeID(uuid.New().String())
}

// String 返回字符串表示
func (id NodeID) String() string {
	return string(id)
}

// PathID 路径唯一标识符
type PathID string

// NewPathID 创建新的路径ID
func NewPathID() PathID {
	return PathID(uuid.New().String())
}

// String 返回字符串表示
func (id PathID) String() string {
	return string(id)
}

// Position 位置信息 - 值对象
type Position struct {
	X float64 `json:"x" gorm:"type:decimal(12,6);default:0"`
	Y float64 `json:"y" gorm:"type:decimal(12,6);default:0"`
	Z float64 `json:"z" gorm:"type:decimal(12,6);default:0"`
}

// DistanceTo 计算到另一个位置的距离
func (p Position) DistanceTo(other Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	dz := p.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// RobotCoordinates 机器人六轴坐标 - 值对象
type RobotCoordinates struct {
	X     float64 `json:"x" gorm:"type:decimal(12,6)"`    // X轴位置
	Y     float64 `json:"y" gorm:"type:decimal(12,6)"`    // Y轴位置
	Z     float64 `json:"z" gorm:"type:decimal(12,6)"`    // Z轴位置
	Roll  float64 `json:"roll" gorm:"type:decimal(8,3)"`  // 翻滚角
	Pitch float64 `json:"pitch" gorm:"type:decimal(8,3)"` // 俯仰角
	Yaw   float64 `json:"yaw" gorm:"type:decimal(8,3)"`   // 偏航角
}

// NodeStyle 节点样式配置 - 值对象
type NodeStyle struct {
	Color       string  `json:"color" gorm:"type:varchar(20);default:'#007bff'"`        // 颜色
	Size        float64 `json:"size" gorm:"type:decimal(5,2);default:10.0"`             // 大小
	Shape       string  `json:"shape" gorm:"type:varchar(20);default:'circle'"`         // 形状
	BorderColor string  `json:"border_color" gorm:"type:varchar(20);default:'#000000'"` // 边框颜色
	BorderWidth float64 `json:"border_width" gorm:"type:decimal(3,1);default:1.0"`      // 边框宽度
	Opacity     float64 `json:"opacity" gorm:"type:decimal(3,2);default:1.0"`           // 透明度
}

// PathStyle 路径样式配置 - 值对象
type PathStyle struct {
	Color   string  `json:"color" gorm:"type:varchar(20);default:'#6c757d'"` // 颜色
	Width   float64 `json:"width" gorm:"type:decimal(3,1);default:2.0"`      // 宽度
	Style   string  `json:"style" gorm:"type:varchar(20);default:'solid'"`   // 样式
	Opacity float64 `json:"opacity" gorm:"type:decimal(3,2);default:1.0"`    // 透明度
}

// ObjectMeta 对象元数据 - 参考Kubernetes ObjectMeta
type ObjectMeta struct {
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	Version     int               `json:"version" gorm:"type:int;default:1"`            // 乐观锁版本
	Labels      map[string]string `json:"labels,omitempty" gorm:"serializer:json"`      // 标签
	Annotations map[string]string `json:"annotations,omitempty" gorm:"serializer:json"` // 注解
}

// === 枚举类型定义 ===

// NodeType 节点类型
type NodeType string

const (
	NodeTypePoint    NodeType = "point"    // 普通点位
	NodeTypeWaypoint NodeType = "waypoint" // 路径点
	NodeTypeStation  NodeType = "station"  // 工作站
	NodeTypeCharging NodeType = "charging" // 充电站
)

// NodeStatus 节点状态
type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"   // 激活
	NodeStatusInactive NodeStatus = "inactive" // 非激活
	NodeStatusDeleted  NodeStatus = "deleted"  // 已删除
	NodeStatusError    NodeStatus = "error"    // 错误状态
)

// PathType 路径类型
type PathType string

const (
	PathTypeNormal     PathType = "normal"     // 普通路径
	PathTypeCurved     PathType = "curved"     // 曲线路径
	PathTypeRestricted PathType = "restricted" // 受限路径
)

// PathStatus 路径状态
type PathStatus string

const (
	PathStatusActive   PathStatus = "active"   // 激活
	PathStatusInactive PathStatus = "inactive" // 非激活
	PathStatusBlocked  PathStatus = "blocked"  // 阻塞
	PathStatusDeleted  PathStatus = "deleted"  // 已删除
)

// CurveType 曲线类型
type CurveType string

const (
	CurveTypeLinear CurveType = "linear" // 线性
	CurveTypeBezier CurveType = "bezier" // 贝塞尔曲线
	CurveTypeSpline CurveType = "spline" // 样条曲线
)

// === 工厂方法 ===

// NewNode 创建新节点
func NewNode(name, nodeType string) *Node {
	return &Node{
		ID:     NewNodeID(),
		Name:   name,
		Type:   NodeType(nodeType),
		Status: NodeStatusActive,
		Style: NodeStyle{
			Color:       "#007bff",
			Size:        10.0,
			Shape:       "circle",
			BorderColor: "#000000",
			BorderWidth: 1.0,
			Opacity:     1.0,
		},
		Metadata: ObjectMeta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}
}

// NewPath 创建新路径
func NewPath(name string, startNodeID, endNodeID NodeID) *Path {
	return &Path{
		ID:          NewPathID(),
		Name:        name,
		Type:        PathTypeNormal,
		Status:      PathStatusActive,
		StartNodeID: startNodeID,
		EndNodeID:   endNodeID,
		Direction:   "bidirectional",
		CurveType:   CurveTypeLinear,
		Style: PathStyle{
			Color:   "#6c757d",
			Width:   2.0,
			Style:   "solid",
			Opacity: 1.0,
		},
		Metadata: ObjectMeta{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}
}

// === 业务方法 ===

// IsValid 验证节点有效性
func (n *Node) IsValid() error {
	if n.Name == "" {
		return fmt.Errorf("节点名称不能为空")
	}
	return nil
}

// UpdatedAt 更新时间戳
func (n *Node) UpdatedAt() {
	n.Metadata.UpdatedAt = time.Now()
	n.Metadata.Version++
}

// IsValid 验证路径有效性
func (p *Path) IsValid() error {
	if p.Name == "" {
		return fmt.Errorf("路径名称不能为空")
	}
	if p.StartNodeID == "" || p.EndNodeID == "" {
		return fmt.Errorf("路径的起始节点和结束节点不能为空")
	}
	if p.StartNodeID == p.EndNodeID {
		return fmt.Errorf("路径的起始节点和结束节点不能相同")
	}
	if p.Weight < 0 {
		return fmt.Errorf("路径权重不能为负数")
	}
	return nil
}

// UpdatedAt 更新时间戳
func (p *Path) UpdatedAt() {
	p.Metadata.UpdatedAt = time.Now()
	p.Metadata.Version++
}
