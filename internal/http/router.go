package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/example/rms/api/docs"
	"github.com/example/rms/internal/http/handlers"
	"github.com/example/rms/internal/repository"
)

// NewRouter wires all routes and handlers.
func NewRouter(repo *repository.Repository) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := &handlers.Handler{Repo: repo}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		handlers.RegisterCustomers(api, h)
		handlers.RegisterEmployees(api, h)
		handlers.RegisterTables(api, h)
		handlers.RegisterMenuCategories(api, h)
		handlers.RegisterDishes(api, h)
		handlers.RegisterProducts(api, h)
		handlers.RegisterReservations(api, h)
		handlers.RegisterOrders(api, h)
		handlers.RegisterPayments(api, h)
		handlers.RegisterReports(api, h)
		handlers.RegisterBatchImport(api, h)
	}

	return r
}
