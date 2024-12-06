package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

const (
	_refreshTokenLifeTime = 24 * 30 * time.Hour
)

type RefreshTokenRepository interface {
	Get(ctx context.Context, refreshToken string) (string, error)
	Set(ctx context.Context, refreshToken string) error
	Delete(ctx context.Context, refreshToken string) error
}

type RefreshTokenRepositoryImpl struct {
	redisClient *redis.Client
}

func NewRefreshTokenRepository(redisClient *redis.Client) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{
		redisClient: redisClient,
	}
}

func (r *RefreshTokenRepositoryImpl) Get(ctx context.Context, refreshToken string) (string, error) {
	userKey := r.getUserKey(refreshToken)

	token, err := r.redisClient.Get(ctx, userKey).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *RefreshTokenRepositoryImpl) Set(ctx context.Context, refreshToken string) error {
	userKey := r.getUserKey(refreshToken)

	refreshTokenTTl := r.redisClient.TTL(ctx, userKey).Val()

	if refreshTokenTTl <= 0 {
		refreshTokenTTl = _refreshTokenLifeTime
	}

	return r.redisClient.Set(ctx, userKey, refreshToken, refreshTokenTTl).Err()
}

func (r *RefreshTokenRepositoryImpl) Delete(ctx context.Context, refreshToken string) error {
	return r.redisClient.Del(ctx, r.getUserKey(refreshToken)).Err()
}

func (r *RefreshTokenRepositoryImpl) getUserKey(token string) string {
	return strings.Split(token, ".")[0]
}
