package database

import (
	"context"
	"fmt"
	"inventory-app/internal/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func GetRedisCache(conf config.RedisDB) *redis.Client {
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)

	cache := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     conf.Password,
		DB:           conf.DB,
		DialTimeout:  5 * time.Minute,
		ReadTimeout:  3 * time.Minute,
		WriteTimeout: 3 * time.Minute,
	})

	if err := cache.Ping(Ctx).Err(); err != nil {
		log.Fatal("Redis ping err: ", err)
	}

	log.Println("✅ Redis Connected")

	return cache
}
