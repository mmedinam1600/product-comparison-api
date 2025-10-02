package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmedinam1600/product-comparison-api/internal/adapters/in/http/handlers"
	"github.com/mmedinam1600/product-comparison-api/internal/adapters/in/http/middleware"
	"github.com/mmedinam1600/product-comparison-api/internal/cache"
	"go.uber.org/zap"
)

type Options struct {
	Mode             string                   // "debug" | "release" | "test"
	CompareHandler   *handlers.CompareHandler // Handler for comparison
	IdempotencyCache *cache.IdempotencyCache  // Idempotency cache
	Logger           *zap.Logger              // Logger
}

func NewEngine(opts Options) *gin.Engine {
	gin.SetMode(opts.Mode)

	router := gin.New()
	// Logger every request and add error recovery handler
	router.Use(gin.Logger(), gin.Recovery())

	// API base group
	api := router.Group("/api")

	// Health check endpoint
	api.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// V1 API group
	v1 := api.Group("/v1")
	{
		// Items group w middleware idempotency
		items := v1.Group("/items")
		items.Use(middleware.IdempotencyMiddleware(opts.IdempotencyCache, opts.Logger))
		{
			// POST /api/v1/items/compare
			items.POST("/compare", opts.CompareHandler.Compare)
		}
	}

	return router
}
