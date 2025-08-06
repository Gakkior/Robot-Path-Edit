// Package domain 模板领域模型
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Template 表示一个模板配置
type Template struct {
	// 基础标识信息
	ID          TemplateID `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null"`
	Description string     `json:"description" gorm:"type:text"`
	Category    string     `json:"category" gorm:"type:varchar(50);default:'custom'"`
	Tags        []string   `json:"tags" gorm:"serializer:json"`

	// 布局类型和配置
	LayoutType   LayoutType             `json:"layout_type" gorm:"type:varchar(30);not null"`
	LayoutConfig map[string]interface{} `json:"layout_config" gorm:"serializer:json"`

	// 模板数据
	TemplateData TemplateData `json:"template_data" gorm:"serializer:json"`

	// 预览信息
	Preview TemplatePreview `json:"preview" gorm:"serializer:json"`

	// 使用统计
	UsageCount int `json:"usage_count" gorm:"type:int;default:0"`

	// 状态信息
	Status   TemplateStatus `json:"status" gorm:"type:varchar(20);default:'active'"`
	IsPublic bool           `json:"is_public" gorm:"type:boolean;default:false"`

	// 元数据
	Metadata TemplateMetadata `json:"metadata" gorm:"embedded"`
}

// TemplateData 模板的实际数据内容
type TemplateData struct {
	Nodes []TemplateNode `json:"nodes"`
	Paths []TemplatePath `json:"paths"`

	// 画布配置
	CanvasConfig CanvasConfig `json:"canvas_config"`

	// 布局参数
	LayoutParams map[string]interface{} `json:"layout_params,omitempty"`
}

// TemplateNode 模板中的节点定义
type TemplateNode struct {
	// 使用相对ID，避免与实际数据冲突
	TemplateID string   `json:"template_id"`
	Name       string   `json:"name"`
	Type       NodeType `json:"type"`

	// 相对位置（0-1之间的比例）
	RelativePosition RelativePosition `json:"relative_position"`

	// 样式配置
	Style NodeStyle `json:"style"`

	// 额外属性
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// TemplatePath 模板中的路径定义
type TemplatePath struct {
	TemplateID      string                 `json:"template_id"`
	Name            string                 `json:"name"`
	Type            PathType               `json:"type"`
	StartNodeTempID string                 `json:"start_node_temp_id"` // 关联模板节点ID
	EndNodeTempID   string                 `json:"end_node_temp_id"`   // 关联模板节点ID
	Direction       string                 `json:"direction"`
	CurveType       CurveType              `json:"curve_type"`
	Style           PathStyle              `json:"style"`
	Properties      map[string]interface{} `json:"properties,omitempty"`
}

// RelativePosition 相对位置（0-1之间的比例坐标）
type RelativePosition struct {
	X float64 `json:"x"` // 0-1之间
	Y float64 `json:"y"` // 0-1之间
	Z float64 `json:"z"` // 0-1之间
}

// CanvasConfig 画布配置
type CanvasConfig struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Zoom        float64 `json:"zoom"`
	CenterX     float64 `json:"center_x"`
	CenterY     float64 `json:"center_y"`
	GridEnabled bool    `json:"grid_enabled"`
	GridSize    int     `json:"grid_size"`
}

// TemplatePreview 模板预览信息
type TemplatePreview struct {
	Thumbnail  string `json:"thumbnail"`  // Base64编码的缩略图
	NodeCount  int    `json:"node_count"` // 节点数量
	PathCount  int    `json:"path_count"` // 路径数量
	Complexity string `json:"complexity"` // 复杂度: simple, medium, complex
	Dimensions string `json:"dimensions"` // 推荐尺寸
}

// TemplateMetadata 模板元数据
type TemplateMetadata struct {
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy   string            `json:"created_by" gorm:"type:varchar(100)"`
	Version     int               `json:"version" gorm:"type:int;default:1"`
	Labels      map[string]string `json:"labels,omitempty" gorm:"serializer:json"`
	Annotations map[string]string `json:"annotations,omitempty" gorm:"serializer:json"`
}

// === 枚举类型定义 ===

// TemplateID 模板唯一标识符
type TemplateID string

// NewTemplateID 创建新的模板ID
func NewTemplateID() TemplateID {
	return TemplateID(uuid.New().String())
}

// String 返回字符串表示
func (id TemplateID) String() string {
	return string(id)
}

// LayoutType 布局类型
type LayoutType string

const (
	LayoutTypeTree      LayoutType = "tree"      // 树形布局
	LayoutTypeGrid      LayoutType = "grid"      // 网格布局
	LayoutTypeCircular  LayoutType = "circular"  // 圆形布局
	LayoutTypeForce     LayoutType = "force"     // 力导向布局
	LayoutTypePipeline  LayoutType = "pipeline"  // 管道布局
	LayoutTypeHierarchy LayoutType = "hierarchy" // 层次布局
	LayoutTypeRadial    LayoutType = "radial"    // 径向布局
	LayoutTypeCustom    LayoutType = "custom"    // 自定义布局
)

// TemplateStatus 模板状态
type TemplateStatus string

const (
	TemplateStatusActive   TemplateStatus = "active"   // 激活
	TemplateStatusInactive TemplateStatus = "inactive" // 非激活
	TemplateStatusDraft    TemplateStatus = "draft"    // 草稿
	TemplateStatusArchived TemplateStatus = "archived" // 已归档
)

// === 工厂方法 ===

// NewTemplate 创建新模板
func NewTemplate(name, description string, layoutType LayoutType) *Template {
	return &Template{
		ID:          NewTemplateID(),
		Name:        name,
		Description: description,
		Category:    "custom",
		LayoutType:  layoutType,
		Status:      TemplateStatusDraft,
		IsPublic:    false,
		TemplateData: TemplateData{
			Nodes: []TemplateNode{},
			Paths: []TemplatePath{},
			CanvasConfig: CanvasConfig{
				Width:       1920,
				Height:      1080,
				Zoom:        1.0,
				GridEnabled: true,
				GridSize:    20,
			},
		},
		Preview: TemplatePreview{
			Complexity: "simple",
			Dimensions: "1920x1080",
		},
		Metadata: TemplateMetadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}
}

// === 业务方法 ===

// IsValid 验证模板有效性
func (t *Template) IsValid() error {
	if t.Name == "" {
		return fmt.Errorf("模板名称不能为空")
	}
	if t.LayoutType == "" {
		return fmt.Errorf("布局类型不能为空")
	}
	return nil
}

// UpdatePreview 更新模板预览信息
func (t *Template) UpdatePreview() {
	t.Preview.NodeCount = len(t.TemplateData.Nodes)
	t.Preview.PathCount = len(t.TemplateData.Paths)

	// 根据节点和路径数量确定复杂度
	totalElements := t.Preview.NodeCount + t.Preview.PathCount
	if totalElements <= 10 {
		t.Preview.Complexity = "simple"
	} else if totalElements <= 50 {
		t.Preview.Complexity = "medium"
	} else {
		t.Preview.Complexity = "complex"
	}
}

// AddNode 向模板添加节点
func (t *Template) AddNode(node TemplateNode) {
	if node.TemplateID == "" {
		node.TemplateID = uuid.New().String()
	}
	t.TemplateData.Nodes = append(t.TemplateData.Nodes, node)
	t.UpdatePreview()
	t.Metadata.UpdatedAt = time.Now()
	t.Metadata.Version++
}

// AddPath 向模板添加路径
func (t *Template) AddPath(path TemplatePath) {
	if path.TemplateID == "" {
		path.TemplateID = uuid.New().String()
	}
	t.TemplateData.Paths = append(t.TemplateData.Paths, path)
	t.UpdatePreview()
	t.Metadata.UpdatedAt = time.Now()
	t.Metadata.Version++
}

// IncrementUsage 增加使用次数
func (t *Template) IncrementUsage() {
	t.UsageCount++
	t.Metadata.UpdatedAt = time.Now()
}

// ToAbsolutePosition 将相对位置转换为绝对位置
func (rp RelativePosition) ToAbsolutePosition(canvasWidth, canvasHeight int) Position {
	return Position{
		X: rp.X * float64(canvasWidth),
		Y: rp.Y * float64(canvasHeight),
		Z: rp.Z * 100, // Z轴使用固定缩放
	}
}

// FromAbsolutePosition 将绝对位置转换为相对位置
func NewRelativePosition(pos Position, canvasWidth, canvasHeight int) RelativePosition {
	return RelativePosition{
		X: pos.X / float64(canvasWidth),
		Y: pos.Y / float64(canvasHeight),
		Z: pos.Z / 100,
	}
}
