package app

import (
	"net/http"

	"github.com/mmedinam1600/product-comparison-api/internal/adapters/in/http/handlers"
	"github.com/mmedinam1600/product-comparison-api/internal/adapters/in/http/router"
	"github.com/mmedinam1600/product-comparison-api/internal/cache"
	"github.com/mmedinam1600/product-comparison-api/internal/data"
	"github.com/mmedinam1600/product-comparison-api/internal/service"
	"github.com/mmedinam1600/product-comparison-api/internal/shared/config"
	"go.uber.org/zap"
)

type App struct {
	HTTPServer       *http.Server
	Logger           *zap.Logger
	RequestCache     *cache.RequestCache
	IdempotencyCache *cache.IdempotencyCache
}

// Bootstrap initializes all the components of the application
func Bootstrap(cfg config.Config) (*App, error) {
	// === 1. Initialize Logger ===
	var logger *zap.Logger
	var err error

	if cfg.GinMode == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	logger.Info("application starting",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.Port),
	)

	// === 2. Initialize Catalog Repository ===
	catalogRepo, err := data.NewFileCatalogRepo(cfg.DataFile, logger)
	if err != nil {
		logger.Fatal("failed to initialize catalog repository", zap.Error(err))
		return nil, err
	}

	// === 3. Initialize Request Cache ===
	requestCache, err := cache.NewRequestCache(cfg.CacheSize, cfg.CacheTTL, logger)
	if err != nil {
		logger.Fatal("failed to initialize request cache", zap.Error(err))
		return nil, err
	}

	// === 4. Initialize Idempotency Cache ===
	idempotencyCache, err := cache.NewIdempotencyCache(cfg.IdempotencySize, cfg.IdempotencyTTL, logger)
	if err != nil {
		logger.Fatal("failed to initialize idempotency cache", zap.Error(err))
		return nil, err
	}

	// === 5. Initialize Compare Service ===
	compareService := service.NewCompareService(catalogRepo, logger)

	// === 6. Inicializar Handler ===
	compareHandler := handlers.NewCompareHandler(
		compareService,
		requestCache,
		idempotencyCache,
		logger,
	)

	// === 7. Create HTTP Engine ===
	engine := router.NewEngine(router.Options{
		Mode:             cfg.GinMode,
		CompareHandler:   compareHandler,
		IdempotencyCache: idempotencyCache,
		Logger:           logger,
	})

	// === 8. Configure HTTP Server with timeouts ===
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	logger.Info("application initialized successfully")

	return &App{
		HTTPServer:       httpServer,
		Logger:           logger,
		RequestCache:     requestCache,
		IdempotencyCache: idempotencyCache,
	}, nil
}

// Shutdown cleans up resources of the application
func (a *App) Shutdown() {
	a.Logger.Info("shutting down application")

	if a.RequestCache != nil {
		a.RequestCache.Close()
	}

	if a.IdempotencyCache != nil {
		a.IdempotencyCache.Close()
	}

	_ = a.Logger.Sync()
}
