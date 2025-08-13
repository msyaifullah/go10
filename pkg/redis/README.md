# Redis Client Package

This package provides a comprehensive Redis client implementation for caching, pub/sub messaging, and other Redis operations.

## Features

- **Redis Client**: Low-level Redis operations with connection pooling and error handling
- **Cache Service**: High-level caching operations with JSON serialization
- **Pub/Sub Service**: Publish/subscribe messaging with message handlers
- **Connection Management**: Automatic connection pooling and health checks
- **Error Handling**: Comprehensive error handling and logging

## Configuration

Add Redis configuration to your TOML config file:

```toml
[redis]
host = "localhost"
port = 6379
password = ""
db = 0
pool_size = 10
min_idle_conns = 2
max_retries = 3
dial_timeout = "5s"
read_timeout = "3s"
write_timeout = "3s"
pool_timeout = "4s"
idle_timeout = "5m"
```

## Usage

### Basic Redis Client

```go
import "loan-service/pkg/redis"

// Initialize Redis client
redisClient, err := redis.NewConnection(cfg.Redis, logger)
if err != nil {
    log.Fatalf("Failed to connect to Redis: %v", err)
}
defer redisClient.Close()

ctx := context.Background()

// Basic operations
err = redisClient.Set(ctx, "key", "value", time.Hour)
value, err := redisClient.Get(ctx, "key")
err = redisClient.Del(ctx, "key")
```

### Cache Service

```go
// Initialize cache service
cacheService := redis.NewCacheService(redisClient, logger)

// Set cache with expiration
err = cacheService.SetCache(ctx, "user:123", userData, 30*time.Minute)

// Get cache
var user User
err = cacheService.GetCacheAs(ctx, "user:123", &user)

// Check if key exists
exists, err := cacheService.Exists(ctx, "user:123")

// Get TTL
ttl, err := cacheService.TTL(ctx, "user:123")
```

### Pub/Sub Service

```go
// Initialize pub/sub service
pubSubService := redis.NewPubSubService(redisClient, logger)

// Subscribe to channels
pubSubService.Subscribe("notifications", func(channel, payload string) error {
    fmt.Printf("Received: %s on %s\n", payload, channel)
    return nil
})

// Publish messages
err = pubSubService.Publish(ctx, "notifications", map[string]interface{}{
    "type": "email",
    "message": "Welcome!",
})

// Start listening
go func() {
    if err := pubSubService.StartListening(ctx); err != nil {
        log.Printf("Failed to start listening: %v", err)
    }
}()
```

## Redis Operations

### String Operations
- `Set(key, value, expiration)` - Set key with expiration
- `Get(key)` - Get string value
- `GetBytes(key)` - Get bytes value
- `Del(keys...)` - Delete keys
- `Exists(keys...)` - Check if keys exist
- `Expire(key, duration)` - Set expiration
- `TTL(key)` - Get time to live

### Hash Operations
- `HSet(key, field, value)` - Set hash field
- `HGet(key, field)` - Get hash field
- `HGetAll(key)` - Get all hash fields
- `HDel(key, fields...)` - Delete hash fields

### List Operations
- `LPush(key, values...)` - Push to left of list
- `RPush(key, values...)` - Push to right of list
- `LPop(key)` - Pop from left of list
- `RPop(key)` - Pop from right of list
- `LLen(key)` - Get list length

### Set Operations
- `SAdd(key, members...)` - Add set members
- `SRem(key, members...)` - Remove set members
- `SMembers(key)` - Get all set members
- `SIsMember(key, member)` - Check membership

### Sorted Set Operations
- `ZAdd(key, score, member)` - Add with score
- `ZRange(key, start, stop)` - Get range by rank
- `ZRem(key, members...)` - Remove members

### Pub/Sub Operations
- `Publish(channel, message)` - Publish to channel
- `Subscribe(channels...)` - Subscribe to channels
- `Listen(channels, handler)` - Listen with handler

## Advanced Features

### Connection Pooling
The Redis client automatically manages connection pooling with configurable settings:
- Pool size
- Minimum idle connections
- Connection timeouts
- Retry policies

### Error Handling
All operations include comprehensive error handling with detailed error messages and logging.

### Message Serialization
The pub/sub service automatically handles JSON serialization/deserialization of messages.

### Delayed Publishing
Support for delayed message publishing using Redis expiration.

## Examples

See `example.go` for comprehensive usage examples including:
- Basic caching
- Pub/Sub messaging
- Advanced caching patterns
- Hash operations
- List operations
- Delayed publishing

## Dependencies

- `github.com/redis/go-redis/v9` - Redis client library
- `loan-service/pkg/config` - Configuration management
- `loan-service/pkg/logger` - Logging

## Running Examples

```bash
# Make sure Redis is running
redis-server

# Run the example
go run pkg/redis/example.go
```

## Best Practices

1. **Connection Management**: Always close Redis connections when done
2. **Error Handling**: Check errors from all Redis operations
3. **Context Usage**: Use context for cancellation and timeouts
4. **Key Naming**: Use descriptive key names with colons (e.g., `user:123:profile`)
5. **Expiration**: Set appropriate TTL for cached data
6. **Message Handlers**: Keep message handlers lightweight and handle errors gracefully
7. **Connection Pooling**: Configure pool size based on your application needs
