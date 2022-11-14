package sqlDB

import (
	"errors"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

func migrations() error {
	if !config.EXEC_MIGRATION {
		logging.Info("Ignoring migration because env variable EXEC_MIGRATION is set to false")
		return nil
	}

	sourceUrl := os.Getenv("MIGRATION_SOURCE_URL")
	if sourceUrl == "" {
		logging.Warn("Migration env variable MIGRATION_SOURCE_URL is empty, using default value ${PWD}/migrations")
		sourceUrl = "./migrations"
	}

	logging.Info("Starting migration execution")
	driver, err := postgres.WithInstance(instance, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("database: could not create migration connection: %v", err)
	}

	logging.Info("Executing migrations on path: %s", sourceUrl)
	m, _ := migrate.NewWithDatabaseInstance(
		"file://"+sourceUrl,
		config.DB_NAME, driver,
	)

	if m != nil {
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("database: error when executing database migration: %v", err)
		}
	}

	logging.Info("Finalized migrations execution.")
	return nil
}
