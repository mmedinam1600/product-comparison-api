package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmedinam1600/product-comparison-api/internal/cache"
	"github.com/mmedinam1600/product-comparison-api/internal/domain"
	"go.uber.org/zap"
)

const IdempotencyKeyHeader = "Idempotency-Key"

// IdempotencyMiddleware manages the idempotency of requests
func IdempotencyMiddleware(idempotencyCache *cache.IdempotencyCache, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if there is an Idempotency-Key header
		idempotencyKey := c.GetHeader(IdempotencyKeyHeader)
		if idempotencyKey == "" {
			// No idempotency key, continue normally
			c.Next()
			return
		}

		// Read the body of the request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("failed to read request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"data":     nil,
				"metadata": nil,
				"error": domain.ErrorResponse{
					ErrorCode: domain.ErrorCodeInvalidRequest,
					Message:   "Failed to read request body.",
				},
			})
			c.Abort()
			return
		}

		// Restore the body for other handlers to read it
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Generate hash of the body
		bodyHash := cache.HashBody(body)

		// Search in idempotency cache
		if entry, found := idempotencyCache.Get(c.Request.Context(), idempotencyKey); found {
			// There is already a request with this idempotency key
			if entry.BodyHash != bodyHash {
				// Body different → Conflict
				logger.Warn("idempotency conflict",
					zap.String("key", idempotencyKey),
					zap.String("expected_hash", entry.BodyHash),
					zap.String("received_hash", bodyHash),
				)
				c.JSON(http.StatusConflict, gin.H{
					"data":     nil,
					"metadata": nil,
					"error": domain.ErrorResponse{
						ErrorCode: domain.ErrorCodeConflict,
						Message:   "Request with same Idempotency-Key but different body already exists.",
					},
				})
				c.Abort()
				return
			}

			// Body equal → return cached response
			logger.Info("returning cached idempotent response", zap.String("key", idempotencyKey))
			c.JSON(http.StatusOK, entry.Response)
			c.Abort()
			return
		}

		// Save the body hash in the context for later use
		c.Set("idempotency_key", idempotencyKey)
		c.Set("body_hash", bodyHash)

		c.Next()
	}
}
