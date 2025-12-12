package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterMenuCategories registers menu categories endpoints.
func RegisterMenuCategories(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/menu-categories")
	g.GET("", h.listMenuCategories)
	g.POST("", h.upsertMenuCategory)
	g.PUT("/:id", h.upsertMenuCategory)
	g.DELETE("/:id", h.deleteMenuCategory)
}

// listMenuCategories godoc
// @Summary List menu categories
// @Tags menu-categories
// @Produce json
// @Success 200 {array} domain.MenuCategory
// @Router /menu-categories [get]
func (h *Handler) listMenuCategories(c *gin.Context) {
	cats, err := h.Repo.ListMenuCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// upsertMenuCategory godoc
// @Summary Create or update menu category
// @Tags menu-categories
// @Accept json
// @Produce json
// @Param category body domain.MenuCategory true "category"
// @Success 200 {object} domain.MenuCategory
// @Router /menu-categories [post]
func (h *Handler) upsertMenuCategory(c *gin.Context) {
	var req domain.MenuCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if err := h.Repo.UpsertMenuCategory(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// deleteMenuCategory godoc
// @Summary Delete menu category
// @Tags menu-categories
// @Param id path int true "category id"
// @Success 204
// @Router /menu-categories/{id} [delete]
func (h *Handler) deleteMenuCategory(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteMenuCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
