package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"go.uber.org/zap"
)

// IdempotencyCache manages the idempotency cache
type IdempotencyCache struct {
	cache  *ristretto.Cache[string, IdempotentEntry]
	logger *zap.Logger
	ttl    time.Duration
}

// IdempotentEntry represents an idempotency entry
type IdempotentEntry struct {
	BodyHash string      // Hash of the original request body
	Response interface{} // Complete cached response
}

// NewIdempotencyCache creates a new instance of the idempotency cache
func NewIdempotencyCache(maxSize int64, ttl time.Duration, logger *zap.Logger) (*IdempotencyCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, IdempotentEntry]{
		NumCounters: maxSize * 10,
		MaxCost:     maxSize,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	logger.Info("idempotency cache initialized",
		zap.Int64("max_size", maxSize),
		zap.Duration("ttl", ttl),
	)

	return &IdempotencyCache{
		cache:  cache,
		logger: logger,
		ttl:    ttl,
	}, nil
}

// Get gets an idempotency entry
func (c *IdempotencyCache) Get(ctx context.Context, key string) (IdempotentEntry, bool) {
	if value, found := c.cache.Get(key); found {
		c.logger.Debug("idempotency hit", zap.String("key", key))
		return value, true
	}
	c.logger.Debug("idempotency miss", zap.String("key", key))
	return IdempotentEntry{}, false
}

// Set store an idempotency entry
func (c *IdempotencyCache) Set(ctx context.Context, key string, entry IdempotentEntry) {
	c.cache.SetWithTTL(key, entry, 1, c.ttl)
	c.logger.Debug("idempotency set", zap.String("key", key), zap.Duration("ttl", c.ttl))
}

// HashBody generates a SHA-256 hash of the body
func HashBody(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}

// Close closes the cache
func (c *IdempotencyCache) Close() {
	c.cache.Close()
	c.logger.Info("idempotency cache closed")
}
