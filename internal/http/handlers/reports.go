package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

var (
	_ domain.ShiftRevenue
	_ domain.WaiterPerformance
	_ domain.DishAvailability
	_ domain.PopularDish
	_ domain.InventoryItem
)

// RegisterReports registers reporting endpoints.
func RegisterReports(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/reports")
	g.GET("/shift-revenue", h.getShiftRevenue)
	g.GET("/shifts", h.getShiftReport)
	g.GET("/waiters", h.getWaiterPerformance)
	g.GET("/dishes-availability", h.getDishesAvailability)
	g.GET("/dishes/popular", h.getPopularDishes)
	g.GET("/inventory", h.getInventoryReport)
}

// getShiftRevenue godoc
// @Summary Shift revenue view
// @Tags reports
// @Produce json
// @Success 200 {array} domain.ShiftRevenue
// @Router /reports/shift-revenue [get]
func (h *Handler) getShiftRevenue(c *gin.Context) {
	data, err := h.Repo.GetShiftRevenue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// getWaiterPerformance godoc
// @Summary Waiter performance
// @Tags reports
// @Produce json
// @Success 200 {array} domain.WaiterPerformance
// @Router /reports/waiters [get]
func (h *Handler) getWaiterPerformance(c *gin.Context) {
	data, err := h.Repo.GetWaiterPerformance(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// getDishesAvailability godoc
// @Summary Dishes availability
// @Tags reports
// @Produce json
// @Success 200 {array} domain.DishAvailability
// @Router /reports/dishes-availability [get]
func (h *Handler) getDishesAvailability(c *gin.Context) {
	data, err := h.Repo.GetDishesAvailability(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// getShiftReport godoc
// @Summary Shift report for period
// @Tags reports
// @Produce json
// @Param from query string true "from date (YYYY-MM-DD)"
// @Param to query string true "to date (YYYY-MM-DD)"
// @Success 200 {array} domain.ShiftRevenue
// @Router /reports/shifts [get]
func (h *Handler) getShiftReport(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required"})
		return
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
		return
	}
	data, err := h.Repo.GetShiftReport(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// getPopularDishes godoc
// @Summary Popular dishes
// @Tags reports
// @Produce json
// @Param limit query int false "limit"
// @Success 200 {array} domain.PopularDish
// @Router /reports/dishes/popular [get]
func (h *Handler) getPopularDishes(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.Repo.GetPopularDishes(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// getInventoryReport godoc
// @Summary Inventory report
// @Tags reports
// @Produce json
// @Success 200 {array} domain.InventoryItem
// @Router /reports/inventory [get]
func (h *Handler) getInventoryReport(c *gin.Context) {
	data, err := h.Repo.GetInventoryReport(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
