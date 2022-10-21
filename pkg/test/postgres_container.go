package test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/google/uuid"

	"github.com/docker/go-connections/nat"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresDockerImage = "postgres:14-alpine"
	testDbHost          = "localhost"
	testDbName          = "test_db"
	testDbUser          = "test_user"
	testDbPassword      = "test_password"
	testPostgresSvcPort = "5432"
)

var postgresContainerInstance *PostgresContainer

type PostgresContainer struct {
	pgContainerRequest *testcontainers.ContainerRequest
	pgContainer        testcontainers.Container
	pgDB               *sql.DB
}

func UsePostgresContainer() *PostgresContainer {
	if postgresContainerInstance == nil {
		postgresContainerInstance = newPostgresContainer()
		postgresContainerInstance.start()
	}
	return postgresContainerInstance
}

func newPostgresContainer() *PostgresContainer {
	req := &testcontainers.ContainerRequest{
		Image:        postgresDockerImage,
		ExposedPorts: []string{testPostgresSvcPort},
		Name:         fmt.Sprintf("colibri-project-test-postgres-%s", uuid.New().String()),
		Env: map[string]string{
			"POSTGRES_DB":       testDbName,
			"POSTGRES_USER":     testDbUser,
			"POSTGRES_PASSWORD": testDbPassword,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(testPostgresSvcPort),
			wait.ForSQL(testPostgresSvcPort, config.SQL_DB_DRIVER, func(port nat.Port) string {
				return fmt.Sprintf(
					"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
					testDbHost,
					port.Port(),
					testDbUser,
					testDbPassword,
					testDbName,
				)
			}),
		),
	}

	return &PostgresContainer{pgContainerRequest: req}
}

func (c *PostgresContainer) start() {
	var err error
	ctx := context.Background()
	c.pgContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: *c.pgContainerRequest,
		Started:          true,
	})
	if err != nil {
		logging.Fatal(err.Error())
	}

	testDbPort, err := c.pgContainer.MappedPort(ctx, testPostgresSvcPort)
	if err != nil {
		logging.Fatal(err.Error())
	}

	log.Printf("Test database started at port: %s", testDbPort)
	c.setDatabaseEnv(testDbPort)
	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name='%s' sslmode=%s",
		os.Getenv("DB_HOST"),
		testDbPort.Port(),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		"test-app",
		os.Getenv("DB_SSL_MODE"))
	if c.pgDB, err = sql.Open("postgres", databaseURL); err != nil {
		logging.Fatal(err.Error())
	}
}

func (c PostgresContainer) Dataset(basePath string, scripts ...string) error {
	for _, s := range scripts {
		script, err := c.loadScript(basePath, s)
		if err != nil {
			return err
		}

		if err = c.execScript(script); err != nil {
			return err
		}
	}

	return nil
}

func (c PostgresContainer) loadScript(basePath string, fileName string) (string, error) {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	filePath := fmt.Sprintf("%s%s", basePath, fileName)
	script, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read script file: %v", err)
	}

	return string(script), nil
}

func (c PostgresContainer) execScript(script string) error {
	if _, err := c.pgDB.Exec(script); err != nil {
		return fmt.Errorf("could not execute script: %v", err)
	}

	return nil
}

func (c PostgresContainer) setDatabaseEnv(testDbPort nat.Port) {
	_ = os.Setenv("DB_HOST", testDbHost)
	_ = os.Setenv("DB_PORT", testDbPort.Port())
	_ = os.Setenv("DB_NAME", testDbName)
	_ = os.Setenv("DB_USER", testDbUser)
	_ = os.Setenv("DB_PASSWORD", testDbPassword)
	_ = os.Setenv("DB_SSL_MODE", "disable")
}
