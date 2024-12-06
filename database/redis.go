package database

import (
	"auth-service/settings"
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, redisCfg settings.Redis) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.Database,
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
