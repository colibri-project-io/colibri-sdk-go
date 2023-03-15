package test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	redisDockerImage = "redis:7.0-alpine"
	testRedisSvcPort = "6379"
)

var redisContainerInstance *RedisContainer

type RedisContainer struct {
	redisContainerRequest *testcontainers.ContainerRequest
	redisContainer        testcontainers.Container
	redisClient           *redis.Client
}

func UseRedisContainer() *RedisContainer {
	if redisContainerInstance == nil {
		redisContainerInstance = newRedisContainer()
		redisContainerInstance.start()
	}
	return redisContainerInstance
}

func newRedisContainer() *RedisContainer {
	req := &testcontainers.ContainerRequest{
		Image:        redisDockerImage,
		ExposedPorts: []string{testRedisSvcPort},
		Name:         fmt.Sprintf("colibri-project-test-redis-%s", uuid.New().String()),
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(testRedisSvcPort),
		),
	}

	return &RedisContainer{redisContainerRequest: req}
}

func (c *RedisContainer) start() {
	var err error
	ctx := context.Background()
	c.redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: *c.redisContainerRequest,
		Started:          true,
	})
	if err != nil {
		logging.Fatal(err.Error())
	}

	testDbPort, err := c.redisContainer.MappedPort(ctx, testRedisSvcPort)
	if err != nil {
		logging.Fatal(err.Error())
	}

	log.Printf("Test redis started at port: %s", testDbPort)
	c.setRedisEnv(testDbPort)
	opts := &redis.Options{Addr: fmt.Sprintf("localhost:%s", testDbPort.Port())}
	c.redisClient = redis.NewClient(opts)
}

func (c RedisContainer) setRedisEnv(port nat.Port) {
	_ = os.Setenv(config.ENV_CACHE_URI, fmt.Sprintf("localhost:%s", port.Port()))
}
