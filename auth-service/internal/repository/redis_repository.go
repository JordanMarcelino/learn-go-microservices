package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisRepositoryImpl struct {
	RDB *redis.ClusterClient
}

func NewRedisRepository(rdb *redis.ClusterClient) RedisRepository {
	return &redisRepositoryImpl{
		RDB: rdb,
	}
}

func (r *redisRepositoryImpl) Get(ctx context.Context, key string) (string, error) {
	return r.RDB.Get(ctx, key).Result()
}

func (r *redisRepositoryImpl) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.RDB.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepositoryImpl) Delete(ctx context.Context, key string) error {
	return r.RDB.Del(ctx, key).Err()
}
