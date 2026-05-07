package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr, url string) (*Redis, error) {
	var opts *redis.Options
	var err error

	// Если есть полный URL (Render) — используем его
	if url != "" {
		opts, err = redis.ParseURL(url)
		if err != nil {
			return nil, fmt.Errorf("invalid redis url: %w", err)
		}
	} else {
		opts = &redis.Options{Addr: addr}
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Redis{client: client}, nil
}

// Set сохраняет значение в кэш с TTL
func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

// Get достаёт значение из кэша и десериализует в dest
func (r *Redis) Get(ctx context.Context, key string, dest any) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err // redis.Nil если ключ не найден
	}
	return json.Unmarshal(data, dest)
}

// Delete удаляет ключ из кэша
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
