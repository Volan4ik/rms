package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterProducts registers product endpoints.
func RegisterProducts(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/products")
	g.GET("", h.listProducts)
	g.POST("", h.upsertProduct)
	g.PUT("/:id", h.upsertProduct)
	g.DELETE("/:id", h.deleteProduct)
}

// listProducts godoc
// @Summary List products
// @Tags products
// @Produce json
// @Param limit query int false "limit"
// @Success 200 {array} domain.Product
// @Router /api/products [get]
func (h *Handler) listProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "200"))
	products, err := h.Repo.ListProducts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// upsertProduct godoc
// @Summary Create or update product
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "product"
// @Success 200 {object} domain.Product
// @Router /api/products [post]
func (h *Handler) upsertProduct(c *gin.Context) {
	var req domain.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || req.Unit == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and unit are required"})
		return
	}
	if err := h.Repo.UpsertProduct(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// deleteProduct godoc
// @Summary Delete product
// @Tags products
// @Param id path int true "product id"
// @Success 204
// @Router /api/products/{id} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
