package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterPayments registers payment endpoints.
func RegisterPayments(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/payments")
	g.POST("", h.upsertPayment)
	g.DELETE("/:orderId", h.deletePayment)
}

// upsertPayment godoc
// @Summary Create or update payment
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body domain.Payment true "payment"
// @Success 200 {object} domain.Payment
// @Router /payments [post]
func (h *Handler) upsertPayment(c *gin.Context) {
	var req domain.Payment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.OrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}
	if req.PaidAt.IsZero() {
		req.PaidAt = time.Now()
	}
	if req.Status == "" {
		req.Status = "paid"
	}
	if err := h.Repo.UpsertPayment(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// deletePayment godoc
// @Summary Delete payment by order id
// @Tags payments
// @Param orderId path int true "order id"
// @Success 204
// @Router /payments/{orderId} [delete]
func (h *Handler) deletePayment(c *gin.Context) {
	orderID, ok := parseID(c, "orderId")
	if !ok {
		return
	}
	if err := h.Repo.DeletePayment(c.Request.Context(), orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
