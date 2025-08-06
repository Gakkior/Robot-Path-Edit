// Package repositories 模板仓储实现
package repositories

import (
	"context"

	"robot-path-editor/internal/database"
	"robot-path-editor/internal/domain"
)

// TemplateRepository 模板仓储接口
type TemplateRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, template *domain.Template) error
	GetByID(ctx context.Context, id string) (*domain.Template, error)
	Update(ctx context.Context, template *domain.Template) error
	Delete(ctx context.Context, id string) error

	// 查询操作
	List(ctx context.Context, opts ListTemplatesOptions) ([]*domain.Template, error)
	GetByCategory(ctx context.Context, category string) ([]*domain.Template, error)
	GetByLayoutType(ctx context.Context, layoutType domain.LayoutType) ([]*domain.Template, error)
	GetPublicTemplates(ctx context.Context) ([]*domain.Template, error)
	Search(ctx context.Context, query string) ([]*domain.Template, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByCategory(ctx context.Context, category string) (int64, error)
}

// ListTemplatesOptions 列出模板的选项
type ListTemplatesOptions struct {
	Category   string
	LayoutType domain.LayoutType
	Status     domain.TemplateStatus
	IsPublic   *bool
	Tags       []string
	CreatedBy  string
	Limit      int
	Offset     int
	SortBy     string // name, created_at, updated_at, usage_count
	SortOrder  string // asc, desc
}

// templateRepository GORM实现
type templateRepository struct {
	db database.Database
}

// NewTemplateRepository 创建新的模板仓储实例
func NewTemplateRepository(db database.Database) TemplateRepository {
	return &templateRepository{db: db}
}

// Create 创建模板
func (r *templateRepository) Create(ctx context.Context, template *domain.Template) error {
	return r.db.GORMDB().WithContext(ctx).Create(template).Error
}

// GetByID 根据ID获取模板
func (r *templateRepository) GetByID(ctx context.Context, id string) (*domain.Template, error) {
	var template domain.Template
	err := r.db.GORMDB().WithContext(ctx).Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update 更新模板
func (r *templateRepository) Update(ctx context.Context, template *domain.Template) error {
	return r.db.GORMDB().WithContext(ctx).Save(template).Error
}

// Delete 删除模板
func (r *templateRepository) Delete(ctx context.Context, id string) error {
	return r.db.GORMDB().WithContext(ctx).Delete(&domain.Template{}, "id = ?", id).Error
}

// List 列出模板
func (r *templateRepository) List(ctx context.Context, opts ListTemplatesOptions) ([]*domain.Template, error) {
	var templates []*domain.Template

	query := r.db.GORMDB().WithContext(ctx)

	// 应用过滤条件
	if opts.Category != "" {
		query = query.Where("category = ?", opts.Category)
	}
	if opts.LayoutType != "" {
		query = query.Where("layout_type = ?", opts.LayoutType)
	}
	if opts.Status != "" {
		query = query.Where("status = ?", opts.Status)
	}
	if opts.IsPublic != nil {
		query = query.Where("is_public = ?", *opts.IsPublic)
	}
	if opts.CreatedBy != "" {
		query = query.Where("created_by = ?", opts.CreatedBy)
	}

	// 标签过滤（需要JSON查询）
	if len(opts.Tags) > 0 {
		for _, tag := range opts.Tags {
			query = query.Where("JSON_CONTAINS(tags, ?)", `"`+tag+`"`)
		}
	}

	// 排序
	if opts.SortBy != "" {
		order := opts.SortBy
		if opts.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("updated_at DESC")
	}

	// 分页
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}

	err := query.Find(&templates).Error
	return templates, err
}

// GetByCategory 根据分类获取模板
func (r *templateRepository) GetByCategory(ctx context.Context, category string) ([]*domain.Template, error) {
	var templates []*domain.Template
	err := r.db.GORMDB().WithContext(ctx).Where("category = ? AND status = ?", category, domain.TemplateStatusActive).Find(&templates).Error
	return templates, err
}

// GetByLayoutType 根据布局类型获取模板
func (r *templateRepository) GetByLayoutType(ctx context.Context, layoutType domain.LayoutType) ([]*domain.Template, error) {
	var templates []*domain.Template
	err := r.db.GORMDB().WithContext(ctx).Where("layout_type = ? AND status = ?", layoutType, domain.TemplateStatusActive).Find(&templates).Error
	return templates, err
}

// GetPublicTemplates 获取公开模板
func (r *templateRepository) GetPublicTemplates(ctx context.Context) ([]*domain.Template, error) {
	var templates []*domain.Template
	err := r.db.GORMDB().WithContext(ctx).Where("is_public = ? AND status = ?", true, domain.TemplateStatusActive).Find(&templates).Error
	return templates, err
}

// Search 搜索模板
func (r *templateRepository) Search(ctx context.Context, query string) ([]*domain.Template, error) {
	var templates []*domain.Template
	searchPattern := "%" + query + "%"
	err := r.db.GORMDB().WithContext(ctx).Where(
		"(name LIKE ? OR description LIKE ?) AND status = ?",
		searchPattern, searchPattern, domain.TemplateStatusActive,
	).Find(&templates).Error
	return templates, err
}

// Count 统计模板总数
func (r *templateRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.GORMDB().WithContext(ctx).Model(&domain.Template{}).Count(&count).Error
	return count, err
}

// CountByCategory 按分类统计模板数量
func (r *templateRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	var count int64
	err := r.db.GORMDB().WithContext(ctx).Model(&domain.Template{}).Where("category = ?", category).Count(&count).Error
	return count, err
}
