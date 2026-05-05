package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	// Проверяем соединение
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
