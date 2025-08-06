// Package services 模板业务服务实现
package services

import (
	"context"
	"fmt"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/repositories"
)

// TemplateService 模板业务服务接口
type TemplateService interface {
	// 模板管理
	CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*domain.Template, error)
	GetTemplate(ctx context.Context, id string) (*domain.Template, error)
	UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*domain.Template, error)
	DeleteTemplate(ctx context.Context, id string) error

	// 模板列表和搜索
	ListTemplates(ctx context.Context, req ListTemplatesRequest) (*ListTemplatesResponse, error)
	SearchTemplates(ctx context.Context, query string) ([]*domain.Template, error)
	GetPublicTemplates(ctx context.Context) ([]*domain.Template, error)
	GetTemplatesByCategory(ctx context.Context, category string) ([]*domain.Template, error)

	// 模板应用
	ApplyTemplate(ctx context.Context, templateID string, canvasConfig domain.CanvasConfig) (*ApplyTemplateResponse, error)

	// 从当前画布保存为模板
	SaveAsTemplate(ctx context.Context, req SaveAsTemplateRequest) (*domain.Template, error)

	// 模板复制和导入导出
	CloneTemplate(ctx context.Context, templateID string, newName string) (*domain.Template, error)
	ExportTemplate(ctx context.Context, templateID string) (*ExportTemplateResponse, error)
	ImportTemplate(ctx context.Context, req ImportTemplateRequest) (*domain.Template, error)
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	LayoutType   domain.LayoutType      `json:"layout_type" binding:"required"`
	LayoutConfig map[string]interface{} `json:"layout_config,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	IsPublic     bool                   `json:"is_public"`
	TemplateData domain.TemplateData    `json:"template_data"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	ID           string                 `json:"id" binding:"required"`
	Name         *string                `json:"name,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Category     *string                `json:"category,omitempty"`
	LayoutConfig map[string]interface{} `json:"layout_config,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	IsPublic     *bool                  `json:"is_public,omitempty"`
	Status       *domain.TemplateStatus `json:"status,omitempty"`
	TemplateData *domain.TemplateData   `json:"template_data,omitempty"`
}

// ListTemplatesRequest 列出模板请求
type ListTemplatesRequest struct {
	Category   string                `json:"category,omitempty"`
	LayoutType domain.LayoutType     `json:"layout_type,omitempty"`
	Status     domain.TemplateStatus `json:"status,omitempty"`
	IsPublic   *bool                 `json:"is_public,omitempty"`
	Tags       []string              `json:"tags,omitempty"`
	CreatedBy  string                `json:"created_by,omitempty"`
	Page       int                   `json:"page,omitempty"`
	PageSize   int                   `json:"page_size,omitempty"`
	SortBy     string                `json:"sort_by,omitempty"`
	SortOrder  string                `json:"sort_order,omitempty"`
}

// ListTemplatesResponse 列出模板响应
type ListTemplatesResponse struct {
	Templates  []*domain.Template `json:"templates"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// ApplyTemplateResponse 应用模板响应
type ApplyTemplateResponse struct {
	Nodes []domain.Node `json:"nodes"`
	Paths []domain.Path `json:"paths"`

	// ID映射表，用于前端更新引用
	NodeIDMapping map[string]string `json:"node_id_mapping"`
	PathIDMapping map[string]string `json:"path_id_mapping"`

	Message string `json:"message"`
}

// SaveAsTemplateRequest 保存为模板请求
type SaveAsTemplateRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	LayoutType  domain.LayoutType `json:"layout_type" binding:"required"`
	Tags        []string          `json:"tags,omitempty"`
	IsPublic    bool              `json:"is_public"`

	// 当前画布数据
	Nodes        []domain.Node       `json:"nodes"`
	Paths        []domain.Path       `json:"paths"`
	CanvasConfig domain.CanvasConfig `json:"canvas_config"`
}

// ExportTemplateResponse 导出模板响应
type ExportTemplateResponse struct {
	Template     *domain.Template `json:"template"`
	ExportFormat string           `json:"export_format"` // json, yaml
	Content      string           `json:"content"`
	Filename     string           `json:"filename"`
}

// ImportTemplateRequest 导入模板请求
type ImportTemplateRequest struct {
	Content string `json:"content" binding:"required"`
	Format  string `json:"format"` // json, yaml
	Name    string `json:"name,omitempty"`
}

// templateService 模板服务实现
type templateService struct {
	templateRepo repositories.TemplateRepository
	nodeRepo     repositories.NodeRepository
	pathRepo     repositories.PathRepository
}

// NewTemplateService 创建新的模板服务实例
func NewTemplateService(
	templateRepo repositories.TemplateRepository,
	nodeRepo repositories.NodeRepository,
	pathRepo repositories.PathRepository,
) TemplateService {
	return &templateService{
		templateRepo: templateRepo,
		nodeRepo:     nodeRepo,
		pathRepo:     pathRepo,
	}
}

// CreateTemplate 创建模板
func (s *templateService) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*domain.Template, error) {
	template := domain.NewTemplate(req.Name, req.Description, req.LayoutType)

	if req.Category != "" {
		template.Category = req.Category
	}
	template.Tags = req.Tags
	template.IsPublic = req.IsPublic
	template.LayoutConfig = req.LayoutConfig
	template.TemplateData = req.TemplateData

	// 更新预览信息
	template.UpdatePreview()

	// 验证模板
	if err := template.IsValid(); err != nil {
		return nil, fmt.Errorf("模板验证失败: %w", err)
	}

	// 保存到数据库
	err := s.templateRepo.Create(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("创建模板失败: %w", err)
	}

	return template, nil
}

// GetTemplate 获取模板
func (s *templateService) GetTemplate(ctx context.Context, id string) (*domain.Template, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// UpdateTemplate 更新模板
func (s *templateService) UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*domain.Template, error) {
	// 获取现有模板
	template, err := s.templateRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Category != nil {
		template.Category = *req.Category
	}
	if req.LayoutConfig != nil {
		template.LayoutConfig = req.LayoutConfig
	}
	if req.Tags != nil {
		template.Tags = req.Tags
	}
	if req.IsPublic != nil {
		template.IsPublic = *req.IsPublic
	}
	if req.Status != nil {
		template.Status = *req.Status
	}
	if req.TemplateData != nil {
		template.TemplateData = *req.TemplateData
		template.UpdatePreview()
	}

	// 更新元数据
	template.Metadata.Version++

	// 验证和保存
	if err := template.IsValid(); err != nil {
		return nil, fmt.Errorf("模板验证失败: %w", err)
	}

	err = s.templateRepo.Update(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("更新模板失败: %w", err)
	}

	return template, nil
}

// DeleteTemplate 删除模板
func (s *templateService) DeleteTemplate(ctx context.Context, id string) error {
	return s.templateRepo.Delete(ctx, id)
}

// ListTemplates 列出模板
func (s *templateService) ListTemplates(ctx context.Context, req ListTemplatesRequest) (*ListTemplatesResponse, error) {
	// 设置默认值
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.Page == 0 {
		req.Page = 1
	}

	// 构建查询选项
	opts := repositories.ListTemplatesOptions{
		Category:   req.Category,
		LayoutType: req.LayoutType,
		Status:     req.Status,
		IsPublic:   req.IsPublic,
		Tags:       req.Tags,
		CreatedBy:  req.CreatedBy,
		Limit:      req.PageSize,
		Offset:     (req.Page - 1) * req.PageSize,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	// 获取模板列表
	templates, err := s.templateRepo.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("获取模板列表失败: %w", err)
	}

	// 获取总数
	total, err := s.templateRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取模板总数失败: %w", err)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &ListTemplatesResponse{
		Templates:  templates,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchTemplates 搜索模板
func (s *templateService) SearchTemplates(ctx context.Context, query string) ([]*domain.Template, error) {
	return s.templateRepo.Search(ctx, query)
}

// GetPublicTemplates 获取公开模板
func (s *templateService) GetPublicTemplates(ctx context.Context) ([]*domain.Template, error) {
	return s.templateRepo.GetPublicTemplates(ctx)
}

// GetTemplatesByCategory 根据分类获取模板
func (s *templateService) GetTemplatesByCategory(ctx context.Context, category string) ([]*domain.Template, error) {
	return s.templateRepo.GetByCategory(ctx, category)
}

// ApplyTemplate 应用模板
func (s *templateService) ApplyTemplate(ctx context.Context, templateID string, canvasConfig domain.CanvasConfig) (*ApplyTemplateResponse, error) {
	// 获取模板
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	// 增加使用次数
	template.IncrementUsage()
	s.templateRepo.Update(ctx, template)

	// 创建ID映射表
	nodeIDMapping := make(map[string]string)
	pathIDMapping := make(map[string]string)

	// 转换模板节点为实际节点
	var nodes []domain.Node
	for _, templateNode := range template.TemplateData.Nodes {
		node := domain.NewNode(templateNode.Name, string(templateNode.Type))

		// 转换相对位置为绝对位置
		node.Position = templateNode.RelativePosition.ToAbsolutePosition(
			canvasConfig.Width, canvasConfig.Height,
		)

		node.Style = templateNode.Style
		node.Properties = templateNode.Properties

		// 保存ID映射
		nodeIDMapping[templateNode.TemplateID] = node.ID.String()

		nodes = append(nodes, *node)
	}

	// 转换模板路径为实际路径
	var paths []domain.Path
	for _, templatePath := range template.TemplateData.Paths {
		startNodeID, ok1 := nodeIDMapping[templatePath.StartNodeTempID]
		endNodeID, ok2 := nodeIDMapping[templatePath.EndNodeTempID]

		if !ok1 || !ok2 {
			continue // 跳过无效的路径
		}

		path := domain.NewPath(
			templatePath.Name,
			domain.NodeID(startNodeID),
			domain.NodeID(endNodeID),
		)

		path.Type = templatePath.Type
		path.Direction = templatePath.Direction
		path.CurveType = templatePath.CurveType
		path.Style = templatePath.Style
		path.Properties = templatePath.Properties

		// 保存ID映射
		pathIDMapping[templatePath.TemplateID] = path.ID.String()

		paths = append(paths, *path)
	}

	return &ApplyTemplateResponse{
		Nodes:         nodes,
		Paths:         paths,
		NodeIDMapping: nodeIDMapping,
		PathIDMapping: pathIDMapping,
		Message:       fmt.Sprintf("已应用模板 '%s'，创建了 %d 个节点和 %d 条路径", template.Name, len(nodes), len(paths)),
	}, nil
}

// SaveAsTemplate 保存为模板
func (s *templateService) SaveAsTemplate(ctx context.Context, req SaveAsTemplateRequest) (*domain.Template, error) {
	template := domain.NewTemplate(req.Name, req.Description, req.LayoutType)

	if req.Category != "" {
		template.Category = req.Category
	}
	template.Tags = req.Tags
	template.IsPublic = req.IsPublic

	// 转换当前节点为模板节点
	var templateNodes []domain.TemplateNode
	for i, node := range req.Nodes {
		templateNode := domain.TemplateNode{
			TemplateID: fmt.Sprintf("node_%d", i+1),
			Name:       node.Name,
			Type:       node.Type,
			RelativePosition: domain.NewRelativePosition(
				node.Position, req.CanvasConfig.Width, req.CanvasConfig.Height,
			),
			Style:      node.Style,
			Properties: node.Properties,
		}
		templateNodes = append(templateNodes, templateNode)
	}

	// 创建节点ID映射
	nodeIDMapping := make(map[string]string)
	for i, node := range req.Nodes {
		nodeIDMapping[node.ID.String()] = fmt.Sprintf("node_%d", i+1)
	}

	// 转换当前路径为模板路径
	var templatePaths []domain.TemplatePath
	for i, path := range req.Paths {
		startTempID, ok1 := nodeIDMapping[path.StartNodeID.String()]
		endTempID, ok2 := nodeIDMapping[path.EndNodeID.String()]

		if !ok1 || !ok2 {
			continue // 跳过无效的路径
		}

		templatePath := domain.TemplatePath{
			TemplateID:      fmt.Sprintf("path_%d", i+1),
			Name:            path.Name,
			Type:            path.Type,
			StartNodeTempID: startTempID,
			EndNodeTempID:   endTempID,
			Direction:       path.Direction,
			CurveType:       path.CurveType,
			Style:           path.Style,
			Properties:      path.Properties,
		}
		templatePaths = append(templatePaths, templatePath)
	}

	// 设置模板数据
	template.TemplateData = domain.TemplateData{
		Nodes:        templateNodes,
		Paths:        templatePaths,
		CanvasConfig: req.CanvasConfig,
	}

	template.UpdatePreview()
	template.Status = domain.TemplateStatusActive

	// 保存模板
	err := s.templateRepo.Create(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("保存模板失败: %w", err)
	}

	return template, nil
}

// CloneTemplate 克隆模板
func (s *templateService) CloneTemplate(ctx context.Context, templateID string, newName string) (*domain.Template, error) {
	// 获取原模板
	original, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取原模板失败: %w", err)
	}

	// 创建克隆
	clone := *original
	clone.ID = domain.NewTemplateID()
	clone.Name = newName
	clone.UsageCount = 0
	clone.Metadata.Version = 1

	// 保存克隆
	err = s.templateRepo.Create(ctx, &clone)
	if err != nil {
		return nil, fmt.Errorf("克隆模板失败: %w", err)
	}

	return &clone, nil
}

// ExportTemplate 导出模板
func (s *templateService) ExportTemplate(ctx context.Context, templateID string) (*ExportTemplateResponse, error) {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	// 这里可以实现JSON或YAML格式的导出
	// 简化实现，返回JSON格式
	return &ExportTemplateResponse{
		Template:     template,
		ExportFormat: "json",
		Filename:     fmt.Sprintf("template_%s.json", template.Name),
	}, nil
}

// ImportTemplate 导入模板
func (s *templateService) ImportTemplate(ctx context.Context, req ImportTemplateRequest) (*domain.Template, error) {
	// 这里可以实现从JSON或YAML导入模板的逻辑
	// 简化实现，直接返回错误
	return nil, fmt.Errorf("导入功能待实现")
}
