package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterReports registers reporting endpoints.
func RegisterReports(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/reports")
	g.GET("/shift-revenue", h.getShiftRevenue)
	g.GET("/waiters", h.getWaiterPerformance)
	g.GET("/dishes-availability", h.getDishesAvailability)
}

// getShiftRevenue godoc
// @Summary Shift revenue view
// @Tags reports
// @Produce json
// @Success 200 {array} domain.ShiftRevenue
// @Router /api/reports/shift-revenue [get]
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
// @Router /api/reports/waiters [get]
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
// @Router /api/reports/dishes-availability [get]
func (h *Handler) getDishesAvailability(c *gin.Context) {
	data, err := h.Repo.GetDishesAvailability(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
