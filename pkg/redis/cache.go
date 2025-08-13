// pkg/redis/cache.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"loan-service/pkg/logger"
)

type CacheService struct {
	redis  *RedisClient
	logger *logger.Logger
}

func NewCacheService(redis *RedisClient, logger *logger.Logger) *CacheService {
	return &CacheService{
		redis:  redis,
		logger: logger,
	}
}

// SetCache sets a value in cache with expiration
func (c *CacheService) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	if err := c.redis.Set(ctx, key, data, expiration); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	c.logger.Debug("Cache set successfully", map[string]interface{}{
		"key":        key,
		"expiration": expiration.String(),
	})

	return nil
}

// GetCache retrieves a value from cache
func (c *CacheService) GetCache(ctx context.Context, key string) ([]byte, error) {
	data, err := c.redis.GetBytes(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	c.logger.Debug("Cache retrieved successfully", map[string]interface{}{
		"key": key,
	})

	return data, nil
}

// GetCacheAs retrieves a value from cache and unmarshals it to the target type
func (c *CacheService) GetCacheAs(ctx context.Context, key string, target interface{}) error {
	data, err := c.GetCache(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

// DeleteCache removes a key from cache
func (c *CacheService) DeleteCache(ctx context.Context, key string) error {
	if err := c.redis.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	c.logger.Debug("Cache deleted successfully", map[string]interface{}{
		"key": key,
	})

	return nil
}

// DeleteCachePattern removes multiple keys matching a pattern
func (c *CacheService) DeleteCachePattern(ctx context.Context, pattern string) error {
	// Note: This is a simplified implementation
	// In production, you might want to use SCAN command for large datasets
	if err := c.redis.Del(ctx, pattern); err != nil {
		return fmt.Errorf("failed to delete cache pattern: %w", err)
	}

	c.logger.Debug("Cache pattern deleted successfully", map[string]interface{}{
		"pattern": pattern,
	})

	return nil
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.redis.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}

	return count > 0, nil
}

// TTL gets the time to live for a key
func (c *CacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.redis.TTL(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

// SetExpiration sets expiration for an existing key
func (c *CacheService) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	if err := c.redis.Expire(ctx, key, expiration); err != nil {
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	c.logger.Debug("Cache expiration set successfully", map[string]interface{}{
		"key":        key,
		"expiration": expiration.String(),
	})

	return nil
}

// Increment increments a numeric value in cache
func (c *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	// This would need to be implemented in the Redis client
	// For now, we'll use a simple approach
	val, err := c.redis.Get(ctx, key)
	if err != nil {
		// Key doesn't exist, set it to 1
		if err := c.redis.Set(ctx, key, "1", 24*time.Hour); err != nil {
			return 0, fmt.Errorf("failed to set initial value: %w", err)
		}
		return 1, nil
	}

	// Parse current value and increment
	var current int64
	if _, err := fmt.Sscanf(val, "%d", &current); err != nil {
		return 0, fmt.Errorf("failed to parse current value: %w", err)
	}

	current++
	if err := c.redis.Set(ctx, key, fmt.Sprintf("%d", current), 24*time.Hour); err != nil {
		return 0, fmt.Errorf("failed to set incremented value: %w", err)
	}

	return current, nil
}
