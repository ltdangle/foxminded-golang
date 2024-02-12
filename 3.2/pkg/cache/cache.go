package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheInterface interface {
	Set(key string, val interface{}, duration time.Duration) error
	Get(key string, dest interface{}) error
}

type RedisCache struct {
	redis *redis.Client
}

func NewRedisCache(addr string, pass string, db int) *RedisCache {
	return &RedisCache{
		redis: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       db,
		}),
	}
}

func (r *RedisCache) Set(key string, val interface{}, duration time.Duration) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	err := r.redis.Set(context.Background(), key, buf.Bytes(), duration).Err()
	return err
}

func (r *RedisCache) Get(key string, dest interface{}) error {
    data, err := r.redis.Get(context.Background(), key).Bytes()
    if err != nil {
        return err
    }

    // Check if data is not empty
    if len(data) == 0 {
        return errors.New("no data found for the key")
    }

    buf := bytes.NewBuffer(data)
    dec := gob.NewDecoder(buf)
    return dec.Decode(dest)
}
