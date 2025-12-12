package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterEmployees registers employee endpoints.
func RegisterEmployees(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/employees")
	g.GET("", h.listEmployees)
	g.POST("", h.createEmployee)
	g.PUT("/:id", h.updateEmployee)
	g.DELETE("/:id", h.deleteEmployee)
}

// listEmployees godoc
// @Summary List employees
// @Tags employees
// @Produce json
// @Success 200 {array} domain.Employee
// @Router /employees [get]
func (h *Handler) listEmployees(c *gin.Context) {
	employees, err := h.Repo.ListEmployees(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// createEmployee godoc
// @Summary Create employee
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body domain.Employee true "employee"
// @Success 201 {object} domain.Employee
// @Router /employees [post]
func (h *Handler) createEmployee(c *gin.Context) {
	var req domain.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FullName == "" || req.Phone == "" || req.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "full_name, phone and role_id are required"})
		return
	}
	if err := h.Repo.CreateEmployee(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// updateEmployee godoc
// @Summary Update employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "employee id"
// @Param employee body domain.Employee true "employee"
// @Success 200 {object} domain.Employee
// @Router /employees/{id} [put]
func (h *Handler) updateEmployee(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req domain.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.UpdateEmployee(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	c.JSON(http.StatusOK, req)
}

// deleteEmployee godoc
// @Summary Delete employee
// @Tags employees
// @Param id path int true "employee id"
// @Success 204
// @Router /employees/{id} [delete]
func (h *Handler) deleteEmployee(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.Repo.DeleteEmployee(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
