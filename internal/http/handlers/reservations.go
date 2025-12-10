package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterReservations registers reservation endpoints.
func RegisterReservations(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/reservations")
	g.GET("", h.listReservations)
	g.POST("", h.createReservation)
	g.PUT("/:id/status", h.updateReservationStatus)
	g.DELETE("/:id", h.deleteReservation)
}

// listReservations godoc
// @Summary List reservations
// @Tags reservations
// @Produce json
// @Param status query string false "status filter"
// @Success 200 {array} domain.Reservation
// @Router /api/reservations [get]
func (h *Handler) listReservations(c *gin.Context) {
	status := c.Query("status")
	reservations, err := h.Repo.ListReservations(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

// createReservation godoc
// @Summary Create reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body domain.Reservation true "reservation"
// @Success 201 {object} domain.Reservation
// @Router /api/reservations [post]
func (h *Handler) createReservation(c *gin.Context) {
	var req domain.Reservation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.CustomerID == 0 || req.TableID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id and table_id are required"})
		return
	}
	if err := h.Repo.CreateReservation(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// updateReservationStatus godoc
// @Summary Update reservation status
// @Tags reservations
// @Param id path int true "reservation id"
// @Param status query string true "new status"
// @Success 200
// @Router /api/reservations/{id}/status [put]
func (h *Handler) updateReservationStatus(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	if err := h.Repo.UpdateReservationStatus(c.Request.Context(), id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// deleteReservation godoc
// @Summary Delete reservation
// @Tags reservations
// @Param id path int true "reservation id"
// @Success 204
// @Router /api/reservations/{id} [delete]
func (h *Handler) deleteReservation(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteReservation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
