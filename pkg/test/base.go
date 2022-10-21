package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/cloud"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/google/uuid"
)

var (
	m sync.Mutex
)

func InitializeBaseTest() {
	loadConfig()
}

func InitializeCacheDBTest() {
	UseRedisContainer()
	loadConfig()
}

func InitializeSqlDBTest() {
	UsePostgresContainer()
	loadConfig()
}

func InitializeTestLocalstack() {
	m.Lock()
	basePath := MountAbsolutPath("../../development-environment")
	ctx := context.WithValue(context.Background(), "id", uuid.New().String())
	_ = useLocalstackContainer(ctx, fmt.Sprintf("%s/%s/", basePath, "localstack"))
	loadConfig()
	cloud.Initialize()
	m.Unlock()
}

func loadConfig() {
	_ = os.Setenv("ENVIRONMENT", "test")
	_ = os.Setenv("APP_NAME", "colibri-project-test")
	_ = os.Setenv("APP_TYPE", "service")
	_ = os.Setenv("CLOUD", "aws")

	if err := config.Load(); err != nil {
		log.Fatalf("could not start configuration: %v", err)
		return
	}
	monitoring.Initialize()
}
