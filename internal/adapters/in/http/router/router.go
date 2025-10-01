package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Mode string // "debug" | "release" | "test"
}

func NewEngine(opts Options) *gin.Engine {
	gin.SetMode(opts.Mode)

	router := gin.New()
	// Logger every request and add error recovery handler
	router.Use(gin.Logger(), gin.Recovery())
	api := router.Group("/api")

	// Health check endpoint
	api.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	return router
}
