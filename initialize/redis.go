package initialize

import (
	"context"
	"fx-vote-server/common/app"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func Redis() *redis.Client {

	client := redis.NewClient(getRedisOptions())

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		app.Log.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		app.Log.Info("redis connect ping response:", zap.String("pong", pong))
	}
	return client
}

func getRedisOptions() *redis.Options {
	r := app.Config.Redis
	o := &redis.Options{
		Addr:     r.Address,
		Username: r.Username,
		Password: r.Password,
		DB:       r.DB,
	}

	if r.PoolSize != 0 {
		o.PoolSize = r.PoolSize
	}

	if r.MinIdleConns != 0 {
		o.MinIdleConns = r.MinIdleConns
	}

	if r.DialTimeout != 0 {
		o.DialTimeout = time.Duration(r.DialTimeout) * time.Second
	}

	if r.ReadTimeout != 0 {
		o.ReadTimeout = time.Duration(r.ReadTimeout) * time.Second
	}

	if r.WriteTimeout != 0 {
		o.WriteTimeout = time.Duration(r.WriteTimeout) * time.Second
	}

	return o
}
