package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterTables registers restaurant tables endpoints.
func RegisterTables(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/tables")
	g.GET("", h.listTables)
	g.POST("", h.upsertTable)
	g.PUT("/:id", h.upsertTable) // table_number is unique, so PUT will act similar
	g.DELETE("/:id", h.deleteTable)
}

// listTables godoc
// @Summary List restaurant tables
// @Tags tables
// @Produce json
// @Success 200 {array} domain.RestaurantTable
// @Router /api/tables [get]
func (h *Handler) listTables(c *gin.Context) {
	tables, err := h.Repo.ListTables(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tables)
}

// upsertTable godoc
// @Summary Create or update restaurant table
// @Tags tables
// @Accept json
// @Produce json
// @Param table body domain.RestaurantTable true "table"
// @Success 200 {object} domain.RestaurantTable
// @Router /api/tables [post]
func (h *Handler) upsertTable(c *gin.Context) {
	var req domain.RestaurantTable
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.TableNumber == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_number is required"})
		return
	}
	if err := h.Repo.UpsertTable(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// deleteTable godoc
// @Summary Delete restaurant table
// @Tags tables
// @Param id path int true "table id"
// @Success 204
// @Router /api/tables/{id} [delete]
func (h *Handler) deleteTable(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteTable(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
