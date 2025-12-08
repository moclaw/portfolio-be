package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client with additional functionality
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg RedisConfig) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("✓ Successfully connected to Redis at %s", addr)

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Set stores a value with expiration
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// Get retrieves a value and unmarshals it into the target
func (r *RedisClient) Get(key string, target interface{}) error {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get value: %w", err)
	}
	return json.Unmarshal(data, target)
}

// Delete removes a key
func (r *RedisClient) Delete(keys ...string) error {
	return r.client.Del(r.ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetNX sets a value only if the key doesn't exist (for distributed locking)
func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.client.SetNX(r.ctx, key, data, expiration).Result()
}

// GetString retrieves a string value
func (r *RedisClient) GetString(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// SetString stores a string value
func (r *RedisClient) SetString(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Increment increments a counter
func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Expire sets expiration on a key
func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// Keys returns all keys matching the pattern
func (r *RedisClient) Keys(pattern string) ([]string, error) {
	return r.client.Keys(r.ctx, pattern).Result()
}

// FlushDB clears the current database
func (r *RedisClient) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}

// HSet sets hash field
func (r *RedisClient) HSet(key string, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.client.HSet(r.ctx, key, field, data).Err()
}

// HGet gets hash field
func (r *RedisClient) HGet(key string, field string, target interface{}) error {
	data, err := r.client.HGet(r.ctx, key, field).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("field not found: %s.%s", key, field)
		}
		return fmt.Errorf("failed to get hash field: %w", err)
	}
	return json.Unmarshal(data, target)
}

// HGetAll gets all hash fields
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// HDel deletes hash fields
func (r *RedisClient) HDel(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// TTL gets the time to live for a key
func (r *RedisClient) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// GetClient returns the underlying Redis client for advanced operations
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Context returns the context
func (r *RedisClient) Context() context.Context {
	return r.ctx
}
