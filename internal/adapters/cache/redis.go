package cache

import (
	"context"
	"github.com/hasElvin/messenger-svc/config"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
	"log"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func InitRedis(cfg config.Config) *redis.Client {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	client := redis.NewClient(opt)

	// Test Redis connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}

func NewRedisCache(client *redis.Client) ports.CacheService {
	return &redisCache{client: client}
}

func (r *redisCache) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
