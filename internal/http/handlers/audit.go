package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

var (
	_ domain.AuditLog
	_ domain.ImportError
)

// RegisterAudit registers audit and import error endpoints.
func RegisterAudit(r *gin.RouterGroup, h *Handler) {
	r.GET("/audit", h.listAuditLog)
	r.GET("/import-errors", h.listImportErrors)
}

// listAuditLog godoc
// @Summary List audit log
// @Tags audit
// @Produce json
// @Param table_name query string false "table name"
// @Param record_id query int false "record id"
// @Param from query string false "from datetime (RFC3339)"
// @Param to query string false "to datetime (RFC3339)"
// @Param limit query int false "limit"
// @Success 200 {array} domain.AuditLog
// @Router /audit [get]
func (h *Handler) listAuditLog(c *gin.Context) {
	tableName := c.Query("table_name")
	var recordID *int64
	if v := c.Query("record_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			recordID = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
			return
		}
	}
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			from = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
			return
		}
	}
	if v := c.Query("to"); v != "" {
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			to = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
			return
		}
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	data, err := h.Repo.ListAuditLog(c.Request.Context(), tableName, recordID, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// listImportErrors godoc
// @Summary List import errors
// @Tags import-errors
// @Produce json
// @Param entity query string false "entity"
// @Param from query string false "from datetime (RFC3339)"
// @Param to query string false "to datetime (RFC3339)"
// @Param limit query int false "limit"
// @Success 200 {array} domain.ImportError
// @Router /import-errors [get]
func (h *Handler) listImportErrors(c *gin.Context) {
	entity := c.Query("entity")
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			from = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
			return
		}
	}
	if v := c.Query("to"); v != "" {
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			to = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to"})
			return
		}
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	data, err := h.Repo.ListImportErrors(c.Request.Context(), entity, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
