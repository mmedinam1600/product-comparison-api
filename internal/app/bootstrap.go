package app

import (
	"net/http"
	"time"

	"github.com/mmedinam1600/product-comparison-api/internal/adapters/in/http/router"
	"github.com/mmedinam1600/product-comparison-api/internal/shared/config"
)

func BuildHTTPServer(cfg config.Config) *http.Server {
	// Engine HTTP
	engine := router.NewEngine(router.Options{Mode: cfg.GinMode})

	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
