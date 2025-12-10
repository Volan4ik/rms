package handlers

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/example/rms/internal/domain"
)

// RegisterBatchImport registers batch-import endpoints.
func RegisterBatchImport(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/batch-import")
	g.POST("/products", h.batchImportProducts)
}

// batchImportProducts godoc
// @Summary Batch import products from JSON array or CSV (name,unit,cost_price,is_available)
// @Tags batch-import
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Router /api/batch-import/products [post]
func (h *Handler) batchImportProducts(c *gin.Context) {
	var products []domain.Product
	ct := c.GetHeader("Content-Type")
	switch {
	case strings.Contains(ct, "application/json"):
		if err := c.ShouldBindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	default:
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "provide JSON array or multipart file"})
			return
		}
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()
		reader := csv.NewReader(f)
		reader.TrimLeadingSpace = true
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if len(rec) < 4 {
				continue
			}
			p := domain.Product{
				Name: rec[0],
				Unit: rec[1],
			}
			if v, err := strconv.ParseFloat(rec[2], 64); err == nil {
				p.CostPrice = &v
			}
			p.IsAvailable = strings.ToLower(rec[3]) == "true"
			products = append(products, p)
		}
	}

	inserted, err := h.Repo.BatchImportProducts(c.Request.Context(), products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"inserted": inserted, "total": len(products)})
}
