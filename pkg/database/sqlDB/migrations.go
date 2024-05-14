package sqlDB

import (
	"database/sql"
	"errors"
	"os"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrationSourceURLEnv       string = "MIGRATION_SOURCE_URL"
	migrationWithPwdDefaultPath string = "${PWD}/migrations"
	migrationDefaultPath        string = "./migrations"

	migrationIgnoringMsg              string = "Ignoring migration because env variable SQL_DB_MIGRATION is set to false"
	migrationEnvNotSetUsingDefaultMsg string = "Migration env variable %s is not set, using default value %s"
	migrationStartingMsg              string = "Starting migration execution"
	migrationCouldNotConnectDBMsg     string = "Could not connect to database for migration: %v"
	migrationExecutingInfoMsg         string = "Executing migration on path: %s"
	migrationExecutionWithErrorMsg    string = "An error when executing database migration: %v"
	migrationFinalizedMsg             string = "Migration finalized successfully"
)

// executeDatabaseMigration performs database migrations based on the provided source URL.
//
// It checks if the SQL_DB_MIGRATION environment variable is set to true before proceeding.
// It uses the MIGRATION_SOURCE_URL environment variable for migration source. If not set, it defaults to "./migrations".
// Returns an error if there is a failure during migration execution.
func executeDatabaseMigration(instance *sql.DB) error {
	if !config.SQL_DB_MIGRATION {
		logging.Info(migrationIgnoringMsg)
		return nil
	}

	sourceUrl := os.Getenv(migrationSourceURLEnv)
	if sourceUrl == "" {
		logging.Warn(migrationEnvNotSetUsingDefaultMsg, migrationSourceURLEnv, migrationWithPwdDefaultPath)
		sourceUrl = migrationDefaultPath
	}

	logging.Info(migrationStartingMsg)
	driver, err := postgres.WithInstance(instance, &postgres.Config{})
	if err != nil {
		logging.Error(migrationCouldNotConnectDBMsg, err)
		return err
	}

	logging.Info(migrationExecutingInfoMsg, sourceUrl)
	migrateDatabaseInstance, _ := migrate.NewWithDatabaseInstance("file://"+sourceUrl, config.SQL_DB_NAME, driver)
	if migrateDatabaseInstance != nil {
		if err = migrateDatabaseInstance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			logging.Error(migrationExecutionWithErrorMsg, err)
			return err
		}
	}

	logging.Info(migrationFinalizedMsg)
	return nil
}
