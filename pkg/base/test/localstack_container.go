package test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	localstackDockerImage = "localstack/localstack:1.4"
	localstackSvcPort     = "4566"
)

var (
	localstackContainerInstance *LocalstackContainer
)

type LocalstackContainer struct {
	lsContainerRequest *testcontainers.ContainerRequest
	lsContainer        testcontainers.Container
	ctx                context.Context
}

func UseLocalstackContainer(ctx context.Context, configPath string) *LocalstackContainer {
	if localstackContainerInstance == nil {
		localstackContainerInstance = newLocalstackContainer(ctx, configPath)
		localstackContainerInstance.start()
	}
	return localstackContainerInstance
}

func newLocalstackContainer(ctx context.Context, configPath string) *LocalstackContainer {
	req := &testcontainers.ContainerRequest{
		Image:        localstackDockerImage,
		ExposedPorts: []string{localstackSvcPort},
		Name:         fmt.Sprintf("colibri-project-test-localstack-%s", uuid.New().String()),
		Env: map[string]string{
			"DEBUG":           "${DEBUG-}",
			"SERVICES":        "sns,sqs,s3,dynamodb",
			"DATA_DIR":        "${DATA_DIR-}",
			"LAMBDA_EXECUTOR": "${LAMBDA_EXECUTOR-}",
			"HOST_TMP_FOLDER": "${TMPDIR:-/tmp/}localstack",
			"DOCKER_HOST":     "unix:///var/run/docker.sock",
		},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(configPath, "/docker-entrypoint-initaws.d/"),
			testcontainers.BindMount("/var/run/docker.sock", "/var/run/docker.sock"),
		),
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(localstackSvcPort),
			wait.ForLog("localstack topics and queues started"),
		),
	}

	return &LocalstackContainer{lsContainerRequest: req, ctx: ctx}
}

func (c *LocalstackContainer) start() {
	var err error
	c.lsContainer, err = testcontainers.GenericContainer(c.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: *c.lsContainerRequest,
		Started:          true,
	})
	if err != nil {
		logging.Fatal(err.Error())
	}
	localstackPort, err := c.lsContainer.MappedPort(c.ctx, localstackSvcPort)
	if err != nil {
		logging.Fatal(err.Error())
	}
	log.Printf("Test localstack started at port: %s", localstackPort)
	c.setEnv(localstackPort)
}

func (c *LocalstackContainer) setEnv(port nat.Port) {
	os.Setenv(config.ENV_CLOUD_HOST, fmt.Sprintf("http://localhost:%s", port.Port()))
	os.Setenv(config.ENV_CLOUD_REGION, "us-east-1")
	os.Setenv(config.ENV_CLOUD_SECRET, "no_secret")
	os.Setenv(config.ENV_CLOUD_TOKEN, "no_token")
	os.Setenv(config.ENV_CLOUD_DISABLE_SSL, "true")
}
