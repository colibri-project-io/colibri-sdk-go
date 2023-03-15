package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

const (
	// Environments
	ENV_ENVIRONMENT           string = "ENVIRONMENT"
	ENV_APP_NAME              string = "APP_NAME"
	ENV_APP_TYPE              string = "APP_TYPE"
	ENV_CLOUD                 string = "CLOUD"
	ENV_NEW_RELIC_LICENSE     string = "NEW_RELIC_LICENSE"
	ENV_PORT                  string = "PORT"
	ENV_DEBUG                 string = "DEBUG"
	ENV_SQL_DB_MIGRATION      string = "SQL_DB_MIGRATION"
	ENV_CLOUD_HOST            string = "CLOUD_HOST"
	ENV_CLOUD_REGION          string = "CLOUD_REGION"
	ENV_CLOUD_SECRET          string = "CLOUD_SECRET"
	ENV_CLOUD_TOKEN           string = "CLOUD_TOKEN"
	ENV_CLOUD_DISABLE_SSL     string = "CLOUD_DISABLE_SSL"
	ENV_CACHE_URI             string = "CACHE_URI"
	ENV_CACHE_PASSWORD        string = "CACHE_PASSWORD"
	ENV_SQL_DB_NAME           string = "SQL_DB_NAME"
	ENV_SQL_DB_HOST           string = "SQL_DB_HOST"
	ENV_SQL_DB_PORT           string = "SQL_DB_PORT"
	ENV_SQL_DB_USER           string = "SQL_DB_USER"
	ENV_SQL_DB_PASSWORD       string = "SQL_DB_PASSWORD"
	ENV_SQL_DB_SSL_MODE       string = "SQL_DB_SSL_MODE"
	ENV_SQL_DB_MAX_OPEN_CONNS string = "SQL_DB_MAX_OPEN_CONNS"
	ENV_SQL_DB_MAX_IDLE_CONNS string = "SQL_DB_MAX_IDLE_CONNS"

	// Environment values
	ENVIRONMENT_PRODUCTION        string = "production"
	ENVIRONMENT_SANDBOX           string = "sandbox"
	ENVIRONMENT_DEVELOPMENT       string = "development"
	ENVIRONMENT_TEST              string = "test"
	APP_TYPE_SERVICE              string = "service"
	APP_TYPE_SERVERLESS           string = "serverless"
	CLOUD_AWS                     string = "aws"
	CLOUD_AZURE                   string = "azure"
	CLOUD_GCP                     string = "gcp"
	CLOUD_FIREBASE                string = "firebase"
	SQL_DB_DRIVER                 string = "nrpostgres"
	SQL_DB_CONNECTION_URI_DEFAULT string = "host=%s port=%s user=%s password=%s dbname=%s application_name='%s' sslmode=%s"

	// Errors
	error_enviroment_not_configured                 string = "environment is not configured. Set production, sandbox, development or test"
	error_app_name_not_configured                   string = "app name is not configured"
	error_app_type_not_configured                   string = "app type is not configured. Set service or serverless"
	error_cloud_not_configured                      string = "cloud is not configured. Set aws, azure, gcp or firebase"
	error_production_required_params_not_configured string = "production required params not configured. Set NEW_RELIC_LICENSE"
	error_integer_parse                             string = "could not parse %s, permited int value, got %v: %w"
	error_boolean_parse                             string = "could not parse %s, permited 'true' or 'false', got %v: %w"
)

var (
	ENVIRONMENT       = ""
	APP_NAME          = ""
	APP_TYPE          = ""
	NEW_RELIC_LICENSE = ""
	DEBUG             = false
	PORT              = 8080

	CLOUD             = ""
	CLOUD_HOST        = ""
	CLOUD_REGION      = ""
	CLOUD_SECRET      = ""
	CLOUD_TOKEN       = ""
	CLOUD_DISABLE_SSL = true

	SQL_DB_NAME           = ""
	SQL_DB_CONNECTION_URI = ""
	SQL_DB_MIGRATION      = false
	SQL_DB_MAX_OPEN_CONNS = 10
	SQL_DB_MAX_IDLE_CONNS = 3

	CACHE_URI      = ""
	CACHE_PASSWORD = ""
)

// Load loads and validates all environment variables. It's used in app initialization.
func Load() error {
	godotenv.Load()

	ENVIRONMENT = os.Getenv(ENV_ENVIRONMENT)
	if !slices.Contains([]string{ENVIRONMENT_PRODUCTION, ENVIRONMENT_SANDBOX, ENVIRONMENT_DEVELOPMENT, ENVIRONMENT_TEST}, ENVIRONMENT) {
		return errors.New(error_enviroment_not_configured)
	}

	APP_NAME = os.Getenv(ENV_APP_NAME)
	if APP_NAME == "" {
		return errors.New(error_app_name_not_configured)
	}

	APP_TYPE = os.Getenv(ENV_APP_TYPE)
	if !slices.Contains([]string{APP_TYPE_SERVICE, APP_TYPE_SERVERLESS}, APP_TYPE) {
		return errors.New(error_app_type_not_configured)
	}

	CLOUD = os.Getenv(ENV_CLOUD)
	if !slices.Contains([]string{CLOUD_AWS, CLOUD_AZURE, CLOUD_GCP, CLOUD_FIREBASE}, CLOUD) {
		return errors.New(error_cloud_not_configured)
	}

	NEW_RELIC_LICENSE = os.Getenv(ENV_NEW_RELIC_LICENSE)
	if (ENVIRONMENT == ENVIRONMENT_PRODUCTION) && (NEW_RELIC_LICENSE == "") {
		return errors.New(error_production_required_params_not_configured)
	}

	if err := convertIntEnv(&PORT, ENV_PORT); err != nil {
		return err
	}

	if err := convertBoolEnv(&DEBUG, ENV_DEBUG); err != nil {
		return err
	}

	if err := convertIntEnv(&SQL_DB_MAX_OPEN_CONNS, ENV_SQL_DB_MAX_OPEN_CONNS); err != nil {
		return err
	}

	if err := convertIntEnv(&SQL_DB_MAX_IDLE_CONNS, ENV_SQL_DB_MAX_IDLE_CONNS); err != nil {
		return err
	}

	if err := convertBoolEnv(&SQL_DB_MIGRATION, ENV_SQL_DB_MIGRATION); err != nil {
		return err
	}

	if err := convertBoolEnv(&CLOUD_DISABLE_SSL, ENV_CLOUD_DISABLE_SSL); err != nil {
		return err
	}

	CLOUD_HOST = os.Getenv(ENV_CLOUD_HOST)
	CLOUD_REGION = os.Getenv(ENV_CLOUD_REGION)
	CLOUD_SECRET = os.Getenv(ENV_CLOUD_SECRET)
	CLOUD_TOKEN = os.Getenv(ENV_CLOUD_TOKEN)

	CACHE_URI = os.Getenv(ENV_CACHE_URI)
	CACHE_PASSWORD = os.Getenv(ENV_CACHE_PASSWORD)

	SQL_DB_NAME = os.Getenv(ENV_SQL_DB_NAME)
	SQL_DB_CONNECTION_URI = fmt.Sprintf(SQL_DB_CONNECTION_URI_DEFAULT,
		os.Getenv(ENV_SQL_DB_HOST),
		os.Getenv(ENV_SQL_DB_PORT),
		os.Getenv(ENV_SQL_DB_USER),
		os.Getenv(ENV_SQL_DB_PASSWORD),
		SQL_DB_NAME,
		APP_NAME,
		os.Getenv(ENV_SQL_DB_SSL_MODE))

	return nil
}

// convertBoolEnv loads the value of an environment variable, converts it to boolean and insert the result into a pointer.
func convertBoolEnv(env *bool, envName string) error {
	if envString := os.Getenv(envName); envString != "" {
		var err error
		if *env, err = strconv.ParseBool(envString); err != nil {
			return fmt.Errorf(error_boolean_parse, envName, envString, err)
		}
	}
	return nil
}

// convertIntEnv loads the value of an environment variable, converts it to interger and insert the result into a pointer.
func convertIntEnv(env *int, envName string) error {
	if envString := os.Getenv(envName); envString != "" {
		var err error
		if *env, err = strconv.Atoi(envString); err != nil {
			return fmt.Errorf(error_integer_parse, envName, envString, err)
		}
	}
	return nil
}

// IsProductionEnvironment returns a boolean if is production environment.
func IsProductionEnvironment() bool {
	return ENVIRONMENT == ENVIRONMENT_PRODUCTION
}

// IsSandboxEnvironment returns a boolean if is sandbox environment.
func IsSandboxEnvironment() bool {
	return ENVIRONMENT == ENVIRONMENT_SANDBOX
}

// IsDevelopmentEnvironment returns a boolean if is development environment.
func IsDevelopmentEnvironment() bool {
	return ENVIRONMENT == ENVIRONMENT_DEVELOPMENT
}

// IsTestEnvironment returns a boolean if is test environment.
func IsTestEnvironment() bool {
	return ENVIRONMENT == ENVIRONMENT_TEST
}

// IsCloudEnvironment returns a boolean if is production or sandbox environment.
func IsCloudEnvironment() bool {
	return IsProductionEnvironment() || IsSandboxEnvironment()
}

// IsLocalEnvironment returns a boolean if is development or test environment.
func IsLocalEnvironment() bool {
	return IsDevelopmentEnvironment() || IsTestEnvironment()
}
