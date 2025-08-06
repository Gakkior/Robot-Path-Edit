// Package handlers 模板相关的HTTP处理器
package handlers

import (
	"net/http"

	"robot-path-editor/internal/domain"
	"robot-path-editor/internal/services"

	"github.com/gin-gonic/gin"
)

// 模板相关处理器
func (h *Handlers) ListTemplates(c *gin.Context) {
	var req services.ListTemplatesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.templateService.ListTemplates(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) CreateTemplate(c *gin.Context) {
	var req services.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.templateService.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

func (h *Handlers) GetTemplate(c *gin.Context) {
	id := c.Param("id")
	template, err := h.templateService.GetTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"template": template})
}

func (h *Handlers) UpdateTemplate(c *gin.Context) {
	var req services.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = c.Param("id")
	template, err := h.templateService.UpdateTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

func (h *Handlers) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	err := h.templateService.DeleteTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}

func (h *Handlers) SearchTemplates(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "查询参数不能为空"})
		return
	}

	templates, err := h.templateService.SearchTemplates(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *Handlers) GetPublicTemplates(c *gin.Context) {
	templates, err := h.templateService.GetPublicTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *Handlers) GetTemplatesByCategory(c *gin.Context) {
	category := c.Param("category")
	templates, err := h.templateService.GetTemplatesByCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *Handlers) ApplyTemplate(c *gin.Context) {
	templateID := c.Param("id")

	var canvasConfig struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	// 获取画布配置，如果没有提供则使用默认值
	if err := c.ShouldBindJSON(&canvasConfig); err != nil {
		canvasConfig.Width = 1920
		canvasConfig.Height = 1080
	}

	response, err := h.templateService.ApplyTemplate(c.Request.Context(), templateID, domain.CanvasConfig{
		Width:  canvasConfig.Width,
		Height: canvasConfig.Height,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": response})
}

func (h *Handlers) SaveAsTemplate(c *gin.Context) {
	var req services.SaveAsTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.templateService.SaveAsTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

func (h *Handlers) CloneTemplate(c *gin.Context) {
	templateID := c.Param("id")

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.templateService.CloneTemplate(c.Request.Context(), templateID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

func (h *Handlers) ExportTemplate(c *gin.Context) {
	templateID := c.Param("id")

	response, err := h.templateService.ExportTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"export": response})
}

func (h *Handlers) ImportTemplate(c *gin.Context) {
	var req services.ImportTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.templateService.ImportTemplate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

// 获取模板统计信息
func (h *Handlers) GetTemplateStats(c *gin.Context) {
	// 这里可以实现模板统计功能
	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total_templates":   0,
			"public_templates":  0,
			"categories":        []string{},
			"popular_templates": []gin.H{},
		},
	})
}
