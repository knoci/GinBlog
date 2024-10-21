package redis

import (
	"GinBlog/setting"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(cfg *setting.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	ctx := client.Context()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}
