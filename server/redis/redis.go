package redis

import (
	"context"
	"fmt"

	"github.com/InspectorGadget/realtime-polling-system/config"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Connect(ctx context.Context) error {
	newClient := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.GetConfig("REDIS_HOST"), config.GetConfig("REDIS_PORT")),
			Password: "", // no password set
			DB:       0,  // use default DB
		},
	)

	err := newClient.Ping(ctx).Err()
	if err != nil {
		return err
	}

	// Assign back interface
	rdb = newClient

	return nil
}

func GetRedis() *redis.Client {
	return rdb
}
