package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmedinam1600/product-comparison-api/internal/app"
	"github.com/mmedinam1600/product-comparison-api/internal/shared/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	application, err := app.Bootstrap(cfg)
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}

	defer application.Shutdown()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine to avoid blocking the main thread
	go func() {
		application.Logger.Info("starting HTTP server",
			zap.String("addr", application.HTTPServer.Addr),
		)
		if err := application.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			application.Logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-quit
	application.Logger.Info("shutdown signal received, gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.HTTPServer.Shutdown(ctx); err != nil {
		application.Logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	application.Logger.Info("server stopped gracefully")
}
