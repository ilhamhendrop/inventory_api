package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache struct {
	rdb        *redis.Client
	prefix     string
	compress   bool
	defaultTTL time.Duration
}

func New(rdb *redis.Client, prefix string, defaultTTL time.Duration, compress bool) *Cache {
	return &Cache{
		rdb:        rdb,
		prefix:     prefix,
		defaultTTL: defaultTTL,
		compress:   compress,
	}
}

func (c *Cache) key(k string) string {
	return c.prefix + ":" + k
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	w.Close()
	return b.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	defer r.Close()
	return io.ReadAll(r)
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl ...time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if c.compress {
		b, err = compress(b)
		if err != nil {
			return err
		}
	}

	exp := c.defaultTTL
	if len(ttl) > 0 {
		exp = ttl[0]
	}

	if err := c.rdb.Set(ctx, c.key(key), b, exp).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	b, err := c.rdb.Get(ctx, c.key(key)).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}

	if c.compress {
		b, err = decompress(b)
		if err != nil {
			return err
		}
	}

	return json.Unmarshal(b, dest)
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	var fullKeys []string
	for _, k := range keys {
		fullKeys = append(fullKeys, c.key(k))
	}
	return c.rdb.Del(ctx, fullKeys...).Err()
}

func (c *Cache) DeleteByPrefix(ctx context.Context, prefix string) error {
	pattern := c.prefix + ":" + prefix + "*"

	iter := c.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}
