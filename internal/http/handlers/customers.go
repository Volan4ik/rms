package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterCustomers registers customer endpoints.
func RegisterCustomers(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/customers")
	g.GET("", h.listCustomers)
	g.POST("", h.createCustomer)
	g.PUT("/:id", h.updateCustomer)
	g.DELETE("/:id", h.deleteCustomer)
}

// listCustomers godoc
// @Summary List customers
// @Tags customers
// @Produce json
// @Success 200 {array} domain.Customer
// @Router /customers [get]
func (h *Handler) listCustomers(c *gin.Context) {
	customers, err := h.Repo.ListCustomers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

// createCustomer godoc
// @Summary Create customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body domain.Customer true "customer"
// @Success 201 {object} domain.Customer
// @Router /customers [post]
func (h *Handler) createCustomer(c *gin.Context) {
	var req domain.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FullName == "" || req.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "full_name and phone are required"})
		return
	}
	if err := h.Repo.CreateCustomer(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// updateCustomer godoc
// @Summary Update customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "customer id"
// @Param customer body domain.Customer true "customer"
// @Success 200 {object} domain.Customer
// @Router /customers/{id} [put]
func (h *Handler) updateCustomer(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req domain.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.UpdateCustomer(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	c.JSON(http.StatusOK, req)
}

// deleteCustomer godoc
// @Summary Delete customer
// @Tags customers
// @Param id path int true "customer id"
// @Success 204
// @Router /customers/{id} [delete]
func (h *Handler) deleteCustomer(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteCustomer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
