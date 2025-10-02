package cache

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"go.uber.org/zap"
)

// RequestCache manages the cache of comparison responses
type RequestCache struct {
	cache  *ristretto.Cache[string, CachedResponse]
	logger *zap.Logger
	ttl    time.Duration
}

// CachedResponse represents a cached response
type CachedResponse struct {
	Data     interface{}
	Metadata interface{}
}

// NewRequestCache creates a new instance of the cache
func NewRequestCache(maxSize int64, ttl time.Duration, logger *zap.Logger) (*RequestCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, CachedResponse]{
		NumCounters: maxSize * 10, // Number of keys to track
		MaxCost:     maxSize,      // Maximum cost (number of entries)
		BufferItems: 64,           // Number of keys per buffer
	})
	if err != nil {
		return nil, err
	}

	logger.Info("request cache initialized",
		zap.Int64("max_size", maxSize),
		zap.Duration("ttl", ttl),
	)

	return &RequestCache{
		cache:  cache,
		logger: logger,
		ttl:    ttl,
	}, nil
}

// Get gets a response from the cache
func (c *RequestCache) Get(ctx context.Context, key string) (CachedResponse, bool) {
	if value, found := c.cache.Get(key); found {
		c.logger.Debug("cache hit", zap.String("key", key))
		return value, true
	}
	c.logger.Debug("cache miss", zap.String("key", key))
	return CachedResponse{}, false
}

// Set stores a response in the cache
func (c *RequestCache) Set(ctx context.Context, key string, response CachedResponse) {
	// Cost = 1 (each entry counts as 1)
	c.cache.SetWithTTL(key, response, 1, c.ttl)
	c.logger.Debug("cache set", zap.String("key", key), zap.Duration("ttl", c.ttl))
}

// Close closes the cache
func (c *RequestCache) Close() {
	c.cache.Close()
	c.logger.Info("request cache closed")
}
