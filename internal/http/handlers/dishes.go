package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterDishes registers dishes endpoints.
func RegisterDishes(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/dishes")
	g.GET("", h.listDishes)
	g.POST("", h.upsertDish)
	g.PUT("/:id", h.upsertDish)
	g.DELETE("/:id", h.deleteDish)
}

// listDishes godoc
// @Summary List dishes
// @Tags dishes
// @Produce json
// @Param limit query int false "limit"
// @Success 200 {array} domain.Dish
// @Router /dishes [get]
func (h *Handler) listDishes(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	dishes, err := h.Repo.ListDishes(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

// upsertDish godoc
// @Summary Create or update dish
// @Tags dishes
// @Accept json
// @Produce json
// @Param dish body domain.Dish true "dish"
// @Success 200 {object} domain.Dish
// @Router /dishes [post]
func (h *Handler) upsertDish(c *gin.Context) {
	var req domain.Dish
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || req.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and category_id are required"})
		return
	}
	if err := h.Repo.UpsertDish(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// deleteDish godoc
// @Summary Delete dish
// @Tags dishes
// @Param id path int true "dish id"
// @Success 204
// @Router /dishes/{id} [delete]
func (h *Handler) deleteDish(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteDish(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
