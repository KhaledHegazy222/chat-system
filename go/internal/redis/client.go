package redis

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisClient defines the interface for the methods used by Redis operations
type RedisClient interface {
	KeyExists(ctx context.Context, hashKey, field string) (bool, error)
	IncrementField(ctx context.Context, hashKey, field string) (int64, error)
	PushToQueue(ctx context.Context, queueName string, data interface{}) error
}

// RedisClientWrapper wraps the *redis.Client and implements RedisClient interface
type RedisClientWrapper struct {
	client *redis.Client
}

// NewRedisClientWrapper creates and returns a new RedisClientWrapper
func NewRedisClientWrapper() *RedisClientWrapper {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})

	return &RedisClientWrapper{client: client}
}

// KeyExists checks if a key exists in the Redis hash
func (r *RedisClientWrapper) KeyExists(ctx context.Context, hashKey, field string) (bool, error) {
	return r.client.HExists(ctx, hashKey, field).Result()
}

// IncrementField increments a field in a Redis hash by 1
func (r *RedisClientWrapper) IncrementField(ctx context.Context, hashKey, field string) (int64, error) {
	return r.client.HIncrBy(ctx, hashKey, field, 1).Result()
}

// PushToQueue pushes an item into a Redis list
func (r *RedisClientWrapper) PushToQueue(ctx context.Context, queueName string, data interface{}) error {
	queueMessageData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.RPush(ctx, queueName, queueMessageData).Err()
}
