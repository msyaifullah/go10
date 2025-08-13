// pkg/redis/pubsub.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"loan-service/pkg/logger"
)

type PubSubService struct {
	redis    *RedisClient
	logger   *logger.Logger
	handlers map[string][]MessageHandler
	mu       sync.RWMutex
}

type Message struct {
	Channel   string      `json:"channel"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	ID        string      `json:"id,omitempty"`
}

func NewPubSubService(redis *RedisClient, logger *logger.Logger) *PubSubService {
	return &PubSubService{
		redis:    redis,
		logger:   logger,
		handlers: make(map[string][]MessageHandler),
	}
}

// Publish publishes a message to a channel
func (p *PubSubService) Publish(ctx context.Context, channel string, payload interface{}) error {
	message := Message{
		Channel:   channel,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := p.redis.Publish(ctx, channel, data); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.logger.Info("Message published successfully", map[string]interface{}{
		"channel": channel,
		"payload": payload,
	})

	return nil
}

// PublishWithID publishes a message with a specific ID
func (p *PubSubService) PublishWithID(ctx context.Context, channel string, payload interface{}, id string) error {
	message := Message{
		Channel:   channel,
		Payload:   payload,
		Timestamp: time.Now(),
		ID:        id,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := p.redis.Publish(ctx, channel, data); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.logger.Info("Message published successfully", map[string]interface{}{
		"channel": channel,
		"payload": payload,
		"id":      id,
	})

	return nil
}

// Subscribe subscribes to a channel and registers a message handler
func (p *PubSubService) Subscribe(channel string, handler MessageHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.handlers[channel] == nil {
		p.handlers[channel] = make([]MessageHandler, 0)
	}

	p.handlers[channel] = append(p.handlers[channel], handler)

	p.logger.Info("Handler registered for channel", map[string]interface{}{
		"channel": channel,
	})
}

// Unsubscribe removes a message handler from a channel
func (p *PubSubService) Unsubscribe(channel string, handler MessageHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if handlers, exists := p.handlers[channel]; exists {
		for i, h := range handlers {
			if &h == &handler {
				p.handlers[channel] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}

	p.logger.Info("Handler unregistered from channel", map[string]interface{}{
		"channel": channel,
	})
}

// StartListening starts listening to all subscribed channels
func (p *PubSubService) StartListening(ctx context.Context) error {
	p.mu.RLock()
	channels := make([]string, 0, len(p.handlers))
	for channel := range p.handlers {
		channels = append(channels, channel)
	}
	p.mu.RUnlock()

	if len(channels) == 0 {
		p.logger.Warn("No channels to listen to", map[string]interface{}{})
		return nil
	}

	p.logger.Info("Starting to listen to channels", map[string]interface{}{
		"channels": channels,
	})

	return p.redis.Listen(ctx, channels, p.handleMessage)
}

// handleMessage processes incoming messages and calls registered handlers
func (p *PubSubService) handleMessage(channel, payload string) error {
	p.mu.RLock()
	handlers, exists := p.handlers[channel]
	p.mu.RUnlock()

	if !exists {
		p.logger.Warn("No handlers registered for channel", map[string]interface{}{
			"channel": channel,
		})
		return nil
	}

	var message Message
	if err := json.Unmarshal([]byte(payload), &message); err != nil {
		p.logger.Error("Failed to unmarshal message", map[string]interface{}{
			"channel": channel,
			"payload": payload,
			"error":   err.Error(),
		})
		return err
	}

	// Call all registered handlers for this channel
	for _, handler := range handlers {
		if err := handler(channel, payload); err != nil {
			p.logger.Error("Handler error", map[string]interface{}{
				"channel": channel,
				"error":   err.Error(),
			})
		}
	}

	p.logger.Debug("Message processed successfully", map[string]interface{}{
		"channel": channel,
		"id":      message.ID,
	})

	return nil
}

// GetSubscriberCount returns the number of handlers for a channel
func (p *PubSubService) GetSubscriberCount(channel string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if handlers, exists := p.handlers[channel]; exists {
		return len(handlers)
	}

	return 0
}

// GetChannels returns all channels with registered handlers
func (p *PubSubService) GetChannels() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	channels := make([]string, 0, len(p.handlers))
	for channel := range p.handlers {
		channels = append(channels, channel)
	}

	return channels
}

// Broadcast publishes a message to multiple channels
func (p *PubSubService) Broadcast(ctx context.Context, channels []string, payload interface{}) error {
	for _, channel := range channels {
		if err := p.Publish(ctx, channel, payload); err != nil {
			p.logger.Error("Failed to broadcast to channel", map[string]interface{}{
				"channel": channel,
				"error":   err.Error(),
			})
			// Continue with other channels
		}
	}

	p.logger.Info("Broadcast completed", map[string]interface{}{
		"channels": channels,
		"payload":  payload,
	})

	return nil
}

// PublishDelayed publishes a message with a delay using Redis
func (p *PubSubService) PublishDelayed(ctx context.Context, channel string, payload interface{}, delay time.Duration) error {
	// Create a delayed message key
	delayedKey := fmt.Sprintf("delayed:%s:%d", channel, time.Now().Add(delay).UnixNano())

	message := Message{
		Channel:   channel,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal delayed message: %w", err)
	}

	// Store the message with expiration
	if err := p.redis.Set(ctx, delayedKey, data, delay); err != nil {
		return fmt.Errorf("failed to store delayed message: %w", err)
	}

	p.logger.Info("Delayed message scheduled", map[string]interface{}{
		"channel": channel,
		"delay":   delay.String(),
		"key":     delayedKey,
	})

	return nil
}
