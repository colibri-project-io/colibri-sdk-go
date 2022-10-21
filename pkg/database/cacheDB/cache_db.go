package cacheDB

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
)

var instance *redis.Client

func Initialize() {
	opts := &redis.Options{Addr: config.CACHE_URI, Password: config.CACHE_PASSWORD}

	redisClient := redis.NewClient(opts)
	redisClient.AddHook(nrredis.NewHook(opts))
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logging.Fatal("An error occurred while trying to connect to the cache database. Error: %s", err)
	}

	instance = redisClient
	logging.Info("Cache database connected")
}
