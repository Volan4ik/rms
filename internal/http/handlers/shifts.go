package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
	"github.com/example/rms/internal/repository"
)

// RegisterShifts registers shift endpoints.
func RegisterShifts(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/shifts")
	g.POST("/open", h.openShift)
	g.POST("/close", h.closeShift)
	g.GET("/current/orders", h.listCurrentShiftOrders)
}

type openShiftRequest struct {
	OpenedBy        int64    `json:"opened_by"`
	Note            string   `json:"note"`
	ExpectedRevenue *float64 `json:"expected_revenue"`
}

// openShift godoc
// @Summary Open a new shift
// @Tags shifts
// @Accept json
// @Produce json
// @Param request body openShiftRequest true "shift data"
// @Success 201 {object} domain.Shift
// @Router /shifts/open [post]
func (h *Handler) openShift(c *gin.Context) {
	var req openShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.OpenedBy == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "opened_by is required"})
		return
	}
	shift := domain.Shift{
		OpenedBy:        req.OpenedBy,
		Note:            req.Note,
		ExpectedRevenue: req.ExpectedRevenue,
		Status:          "opened",
	}
	if err := h.Repo.OpenShift(c.Request.Context(), &shift); err != nil {
		if err == repository.ErrOpenShiftExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "shift already opened"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shift)
}

type closeShiftRequest struct {
	ClosedBy int64  `json:"closed_by"`
	Note     string `json:"note"`
}

// closeShift godoc
// @Summary Close current shift
// @Tags shifts
// @Accept json
// @Produce json
// @Param request body closeShiftRequest true "close payload"
// @Success 200 {object} domain.Shift
// @Router /shifts/close [post]
func (h *Handler) closeShift(c *gin.Context) {
	var req closeShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ClosedBy == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "closed_by is required"})
		return
	}
	shift, err := h.Repo.CloseShift(c.Request.Context(), req.ClosedBy, req.Note)
	if err != nil {
		if err == repository.ErrNoOpenShift {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no opened shift"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shift)
}

// listCurrentShiftOrders godoc
// @Summary List orders in current shift
// @Tags shifts
// @Produce json
// @Param status query string false "order status"
// @Param waiter_id query int false "waiter id"
// @Success 200 {array} domain.OrderWithTotal
// @Router /shifts/current/orders [get]
func (h *Handler) listCurrentShiftOrders(c *gin.Context) {
	status := c.Query("status")
	var waiterID *int64
	if v := c.Query("waiter_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			waiterID = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid waiter_id"})
			return
		}
	}
	data, err := h.Repo.GetCurrentShiftOrders(c.Request.Context(), status, waiterID)
	if err != nil {
		if err == repository.ErrNoOpenShift {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no opened shift"})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "orders not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
