package test

import (
	"context"
	"fmt"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	wiremockDockerImage = "wiremock/wiremock:2.32.0-alpine"
	wiremockSvcPort     = "8080"
)

var wiremockContainerInstance *WiremockContainer

type WiremockContainer struct {
	wContainerRequest *testcontainers.ContainerRequest
	wContainer        testcontainers.Container
	configPath        string
	instancePort      int
}

func UseWiremockContainer(configPath string) *WiremockContainer {
	if wiremockContainerInstance == nil {
		wiremockContainerInstance = newWiremockContainer(configPath)
		wiremockContainerInstance.start()
	}
	return wiremockContainerInstance
}

func newWiremockContainer(configPath string) *WiremockContainer {
	req := testcontainers.ContainerRequest{
		Image:        wiremockDockerImage,
		Name:         fmt.Sprintf("colibri-project-test-wiremock-%s", uuid.New().String()),
		ExposedPorts: []string{wiremockSvcPort},
		Env:          map[string]string{},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{Source: testcontainers.GenericBindMountSource{HostPath: configPath}, Target: "/home/wiremock"},
		},
		Cmd:        []string{"--local-response-templating"},
		WaitingFor: wait.ForListeningPort(wiremockSvcPort),
	}

	return &WiremockContainer{wContainerRequest: &req, configPath: configPath}
}

func (c *WiremockContainer) start() {
	var err error
	ctx := context.Background()
	c.wContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: *c.wContainerRequest,
		Started:          true,
	})
	if err != nil {
		logging.Fatal(err.Error())
	}

	runningPort, _ := c.wContainer.MappedPort(ctx, wiremockSvcPort)
	c.instancePort = runningPort.Int()
	logging.Info("Wiremock container exposed port: %s", runningPort.Port())
}

func (c *WiremockContainer) Port() int {
	return c.instancePort
}
