package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmedinam1600/product-comparison-api/internal/cache"
	"github.com/mmedinam1600/product-comparison-api/internal/domain"
	"github.com/mmedinam1600/product-comparison-api/internal/service"
	"go.uber.org/zap"
)

// CompareHandler manages the comparison requests
type CompareHandler struct {
	compareService   service.CompareService
	requestCache     *cache.RequestCache
	idempotencyCache *cache.IdempotencyCache
	logger           *zap.Logger
}

// NewCompareHandler creates a new instance of the handler
func NewCompareHandler(
	compareService service.CompareService,
	requestCache *cache.RequestCache,
	idempotencyCache *cache.IdempotencyCache,
	logger *zap.Logger,
) *CompareHandler {
	return &CompareHandler{
		compareService:   compareService,
		requestCache:     requestCache,
		idempotencyCache: idempotencyCache,
		logger:           logger,
	}
}

// CompareResponse structures the response of the endpoint
type CompareResponse struct {
	Data     *domain.CompareResult `json:"data"`
	Metadata *domain.Metadata      `json:"metadata"`
	Error    *domain.ErrorResponse `json:"error"`
}

// Compare manages POST /api/v1/items/compare
func (h *CompareHandler) Compare(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request
	var req domain.CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, CompareResponse{
			Data:     nil,
			Metadata: nil,
			Error: &domain.ErrorResponse{
				ErrorCode: domain.ErrorCodeMissingField,
				Message:   "Missing mandatory field 'ids'",
			},
		})
		return
	}

	h.logger.Info("compare request received",
		zap.Int("ids_count", len(req.Ids)),
		zap.Bool("has_fields", req.Fields != nil),
	)

	// Generate cache key based on IDs
	cacheKey := h.compareService.GenerateCacheKey(req.Ids)

	// Try to get from the request cache
	if cached, found := h.requestCache.Get(ctx, cacheKey); found {
		h.logger.Info("returning cached response", zap.String("cache_key", cacheKey))
		c.Header("Cache-Status", "hit")
		c.JSON(http.StatusOK, CompareResponse{
			Data:     cached.Data.(*domain.CompareResult),
			Metadata: cached.Metadata.(*domain.Metadata),
			Error:    nil,
		})
		return
	}

	c.Header("Cache-Status", "miss")

	// Execute comparison
	result, metadata, errResp := h.compareService.Compare(ctx, req)

	if errResp != nil {
		// Business error
		statusCode := errResp.ErrorCode.HTTPStatusCode()
		h.logger.Info("comparison failed",
			zap.String("error_code", string(errResp.ErrorCode)),
			zap.Int("status", statusCode),
		)
		c.JSON(statusCode, CompareResponse{
			Data:     nil,
			Metadata: nil,
			Error:    errResp,
		})
		return
	}

	// Success: cache the result
	h.requestCache.Set(ctx, cacheKey, cache.CachedResponse{
		Data:     &result,
		Metadata: &metadata,
	})

	// Prepare final response
	response := CompareResponse{
		Data:     &result,
		Metadata: &metadata,
		Error:    nil,
	}

	// If there is an idempotency key, save in the idempotency cache
	if idempotencyKey, exists := c.Get("idempotency_key"); exists {
		if bodyHash, hashExists := c.Get("body_hash"); hashExists {
			h.idempotencyCache.Set(ctx, idempotencyKey.(string), cache.IdempotentEntry{
				BodyHash: bodyHash.(string),
				Response: response,
			})
			h.logger.Debug("saved idempotent response", zap.String("key", idempotencyKey.(string)))
		}
	}

	c.JSON(http.StatusOK, response)
}
