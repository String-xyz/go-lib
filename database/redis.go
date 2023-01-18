package database

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisRepresentable interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) *redis.StatusCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd
	HLen(ctx context.Context, key string) *redis.IntCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
}

type RedisStore interface {
	Get(id string) ([]byte, error)
	Set(string, any, time.Duration) error
	HSet(string, map[string]interface{}) error
	HGetAll(string) (map[string]string, error)
	HDel(string, string) int64
	HMLen(string) int64
	Delete(string) error
}
