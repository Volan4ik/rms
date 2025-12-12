package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterOrders registers order endpoints.
func RegisterOrders(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/orders")
	g.GET("", h.listOrders)
	g.POST("", h.createOrder)
	g.PUT("/:id/status", h.updateOrderStatus)
	g.GET("/:id/items", h.listOrderItems)
	g.POST("/:id/items", h.addOrderItem)
	g.DELETE("/:id/items/:itemId", h.deleteOrderItem)
}

// listOrders godoc
// @Summary List orders
// @Tags orders
// @Produce json
// @Param status query string false "status filter"
// @Param limit query int false "limit"
// @Success 200 {array} domain.Order
// @Router /orders [get]
func (h *Handler) listOrders(c *gin.Context) {
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	orders, err := h.Repo.ListOrders(c.Request.Context(), status, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

type orderRequest struct {
	TableID       int64              `json:"table_id"`
	CustomerID    *int64             `json:"customer_id"`
	WaiterID      int64              `json:"waiter_id"`
	ReservationID *int64             `json:"reservation_id"`
	ShiftID       *int64             `json:"shift_id"`
	Status        string             `json:"status"`
	Items         []domain.OrderItem `json:"items"`
}

// createOrder godoc
// @Summary Create order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body orderRequest true "order"
// @Success 201 {object} domain.Order
// @Router /orders [post]
func (h *Handler) createOrder(c *gin.Context) {
	var req orderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.TableID == 0 || req.WaiterID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_id and waiter_id are required"})
		return
	}
	order := domain.Order{
		TableID:       req.TableID,
		CustomerID:    req.CustomerID,
		WaiterID:      req.WaiterID,
		ReservationID: req.ReservationID,
		ShiftID:       req.ShiftID,
		Status:        req.Status,
	}
	if order.Status == "" {
		order.Status = "new"
	}
	if err := h.Repo.CreateOrder(c.Request.Context(), &order, req.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// updateOrderStatus godoc
// @Summary Update order status
// @Tags orders
// @Param id path int true "order id"
// @Param status query string true "new status"
// @Success 200
// @Router /orders/{id}/status [put]
func (h *Handler) updateOrderStatus(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	if err := h.Repo.UpdateOrderStatus(c.Request.Context(), id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// listOrderItems godoc
// @Summary List items for order
// @Tags orders
// @Param id path int true "order id"
// @Produce json
// @Success 200 {array} domain.OrderItem
// @Router /orders/{id}/items [get]
func (h *Handler) listOrderItems(c *gin.Context) {
	orderID, ok := parseID(c, "id")
	if !ok {
		return
	}
	items, err := h.Repo.ListOrderItems(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// addOrderItem godoc
// @Summary Add or update order item
// @Tags orders
// @Param id path int true "order id"
// @Accept json
// @Produce json
// @Param item body domain.OrderItem true "item"
// @Success 200
// @Router /orders/{id}/items [post]
func (h *Handler) addOrderItem(c *gin.Context) {
	orderID, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req domain.OrderItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.DishID == 0 || req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dish_id and quantity are required"})
		return
	}
	if err := h.Repo.AddOrderItem(c.Request.Context(), orderID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// deleteOrderItem godoc
// @Summary Delete order item
// @Tags orders
// @Param id path int true "order id"
// @Param itemId path int true "item id"
// @Success 204
// @Router /orders/{id}/items/{itemId} [delete]
func (h *Handler) deleteOrderItem(c *gin.Context) {
	itemID, ok := parseID(c, "itemId")
	if !ok {
		return
	}
	if err := h.Repo.DeleteOrderItem(c.Request.Context(), itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
