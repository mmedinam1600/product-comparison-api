package main

import (
	"log"

	"github.com/mmedinam1600/product-comparison-api/internal/adapters/in/http/router"
	"github.com/mmedinam1600/product-comparison-api/internal/shared/config"
)

func main() {
	cfg := config.Load()
	mainRouter := router.NewEngine(router.Options{Mode: cfg.GinMode})

	if err := mainRouter.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
