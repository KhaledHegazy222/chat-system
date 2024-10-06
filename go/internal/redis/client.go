package redis

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})
}

func KeyExists(ctx context.Context, hashKey, field string) (bool, error) {
	return redisClient.HExists(ctx, hashKey, field).Result()
}

func IncrementField(ctx context.Context, hashKey, field string) (int64, error) {
	return redisClient.HIncrBy(ctx, hashKey, field, 1).Result()
}

func PushToQueue(ctx context.Context, queueName string, data interface{}) error {
	queueMessageData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redisClient.RPush(ctx, queueName, queueMessageData).Err()
}
