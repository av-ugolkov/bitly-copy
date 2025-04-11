package redis

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/bitly-copy/internal/config"

	rClient "github.com/redis/go-redis/v9"
)

const (
	CountScanItem = 100
)

type Redis struct {
	client *rClient.Client
}

func New(cfg *config.Config) *Redis {
	r := rClient.NewClient(&rClient.Options{
		Addr:       cfg.Redis.Addr(),
		ClientName: cfg.Redis.Name,
		DB:         0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.ReadTimeout)
	defer cancel()
	res := r.Ping(ctx)
	if res.Err() != nil {
		panic(res.Err())
	}

	return &Redis{
		client: r,
	}
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis.Get - key [%s]: %w", key, err)
	}

	return value, nil
}

func (r *Redis) GetAll(ctx context.Context, key string) ([]string, error) {
	var cursor, nextCursor uint64
	var err error
	var keys []string
	values := make([]string, 0, CountScanItem)

	for {
		keys, nextCursor, err = r.client.Scan(ctx, cursor, fmt.Sprintf("%s:*", key), CountScanItem).Result()
		if err != nil {
			return nil, fmt.Errorf("redis.GetAll - key [%s]: %w", key, err)
		}

		for _, key := range keys {
			val, err := r.client.Get(ctx, key).Result()
			if err == nil {
				values = append(values, val)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return values, nil
}

func (r *Redis) SetNX(ctx context.Context, key string, value string) (bool, error) {
	success, err := r.client.SetNX(ctx, key, value, 0).Result()
	if err != nil {
		return false, fmt.Errorf("redis.SetNX - key [%s]: %w", key, err)
	}
	return success, nil
}

func (r *Redis) Increment(ctx context.Context, key string) error {
	_, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis.Increment - key [%s]: %w", key, err)
	}
	return nil
}

func (r *Redis) Remove(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
