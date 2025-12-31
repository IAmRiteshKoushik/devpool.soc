package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var Valkey *redis.Client

func InitValkey() error {
	host := App.RedisHost
	port := App.RedisPort
	uname := App.RedisUsername
	passwd := App.RedisPassword
	resp := 3

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Username: uname,
		Password: passwd,
		DB:       0,
		Protocol: resp,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		Log.Fatal("[FAIL]: Heal-check failed for Redis", err)
		return err
	}
	Log.Info(
		fmt.Sprintf("[INFO]: Health-check successful for Redis. Response: %s", pong),
	)

	Valkey = rdb
	return nil
}
