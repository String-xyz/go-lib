package database

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/go-redis/redis/v8"
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

type redisStore struct {
	client RedisRepresentable
}

type RedisConfigOptions struct {
	ClusterMode bool
	Port        string
	Host        string
	Password    string
}

const REDIS_NOT_FOUND_ERROR = "redis: nil"

func redisTLSConf() *tls.Config {
	var tlsCf *tls.Config
	if !common.IsLocalEnv() {
		tlsCf = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return tlsCf
}

func redisOptions(host, port, password string) *redis.Options {
	url := host + ":" + port
	var tlsCf *tls.Config
	if !common.IsLocalEnv() {
		tlsCf = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	op := &redis.Options{
		Addr:      url,
		TLSConfig: tlsCf,
		Password:  password,
		DB:        0,
	}
	return op
}

func cluster(host, port, password string) *redis.ClusterClient {
	url := host + ":" + port
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          []string{url},
		Password:       password,
		PoolSize:       10,
		MinIdleConns:   10,
		TLSConfig:      redisTLSConf(),
		ReadOnly:       false,
		RouteRandomly:  false,
		RouteByLatency: false,
	})
}

func NewRedisStore(options RedisConfigOptions) RedisStore {
	ctx := context.Background()
	var client RedisRepresentable
	if !options.ClusterMode {
		client = redis.NewClient(redisOptions(options.Host, options.Port, options.Password))
	} else {
		client = cluster(options.Host, options.Port, options.Password)
	}
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	return &redisStore{
		client: client,
	}
}

func (r redisStore) Delete(id string) error {
	ctx := context.Background()
	_, err := r.client.Del(ctx, id).Result()
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func (r redisStore) Get(id string) ([]byte, error) {
	ctx := context.Background()
	bytes, err := r.client.Get(ctx, id).Bytes()
	if err != nil {
		return nil, common.StringError(err)
	}
	return bytes, nil
}

func (r redisStore) Set(id string, value any, expire time.Duration) error {
	ctx := context.Background()
	if err := r.client.Set(ctx, id, value, expire).Err(); err != nil {
		return common.StringError(err)
	}
	return nil
}

func (r redisStore) HSet(key string, data map[string]interface{}) error {
	ctx := context.Background()
	if err := r.client.HSet(ctx, key, data).Err(); err != nil {
		return common.StringError(err, "failed to save array to redis")
		// return common.UnexpectedError(err)
	}

	return nil
}

func (r redisStore) HGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return data, common.StringError(err)
	}
	return data, nil
}

func (r redisStore) HMLen(key string) int64 {
	ctx := context.Background()
	data := r.client.HLen(ctx, key)
	return data.Val()
}

func (r redisStore) HDel(key, val string) int64 {
	ctx := context.Background()
	data := r.client.HDel(ctx, key, val)
	return data.Val()
}
