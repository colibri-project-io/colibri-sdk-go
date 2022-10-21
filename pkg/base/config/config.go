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
	ENV_PRODUCTION      string = "production"
	ENV_SANDBOX         string = "sandbox"
	ENV_DEVELOPMENT     string = "development"
	ENV_TEST            string = "test"
	APP_TYPE_SERVICE    string = "service"
	APP_TYPE_SERVERLESS string = "serverless"
	CLOUD_AWS           string = "aws"
	CLOUD_GCP           string = "gcp"
	CLOUD_FIREBASE      string = "firebase"
	SQL_DB              string = "sql"
	SQL_DB_DRIVER       string = "nrpostgres"
	NOSQL_DB            string = "nosql"

	error_enviroment_not_configured                 string = "environment is not configured. Set production, sandbox, development or test"
	error_app_name_not_configured                   string = "app name is not configured"
	error_app_type_not_configured                   string = "app type is not configured. Set service or serverless"
	error_cloud_not_configured                      string = "cloud is not configured. Set aws, gcp or firebase"
	error_production_required_params_not_configured string = "production required params not configured. Set NEW_RELIC_LICENSE"
	error_database_misconfigured                    string = "database is misconfigured. Set sql or nosql"
)

var (
	ENVIRONMENT       string
	APP_NAME          string
	APP_TYPE          string
	CLOUD             string
	NEW_RELIC_LICENSE string
	CLOUD_HOST        string
	CLOUD_REGION      string
	CLOUD_SECRET      string
	CLOUD_TOKEN       string
	CLOUD_DISABLE_SSL bool
	DB                string
	DB_CONNECTION_URI string
	PORT              int
	CACHE_URI         string
	CACHE_PASSWORD    string
)

func Load() error {
	godotenv.Load()

	ENVIRONMENT = os.Getenv("ENVIRONMENT")
	if !slices.Contains([]string{ENV_PRODUCTION, ENV_SANDBOX, ENV_DEVELOPMENT, ENV_TEST}, ENVIRONMENT) {
		return errors.New(error_enviroment_not_configured)
	}

	APP_NAME = os.Getenv("APP_NAME")
	if APP_NAME == "" {
		return errors.New(error_app_name_not_configured)
	}

	APP_TYPE = os.Getenv("APP_TYPE")
	if !slices.Contains([]string{APP_TYPE_SERVICE, APP_TYPE_SERVERLESS}, APP_TYPE) {
		return errors.New(error_app_type_not_configured)
	}

	CLOUD = os.Getenv("CLOUD")
	if !slices.Contains([]string{CLOUD_AWS, CLOUD_GCP, CLOUD_FIREBASE}, CLOUD) {
		return errors.New(error_cloud_not_configured)
	}

	NEW_RELIC_LICENSE = os.Getenv("NEW_RELIC_LICENSE")
	if (ENVIRONMENT == ENV_PRODUCTION) && (NEW_RELIC_LICENSE == "") {
		return errors.New(error_production_required_params_not_configured)
	}

	CLOUD_HOST = os.Getenv("CLOUD_HOST")
	CLOUD_REGION = os.Getenv("CLOUD_REGION")
	CLOUD_SECRET = os.Getenv("CLOUD_SECRET")
	CLOUD_TOKEN = os.Getenv("CLOUD_TOKEN")
	CLOUD_DISABLE_SSL, _ = strconv.ParseBool(os.Getenv("CLOUD_DISABLE_SSL"))

	PORT, _ = strconv.Atoi(os.Getenv("PORT"))

	CACHE_URI = os.Getenv("CACHE_URI")
	CACHE_PASSWORD = os.Getenv("CACHE_PASSWORD")
	DB = os.Getenv("DB")
	if (DB != "") && !slices.Contains([]string{SQL_DB, NOSQL_DB}, DB) {
		return errors.New(error_database_misconfigured)
	}

	DB_CONNECTION_URI = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name='%s' sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		APP_NAME,
		os.Getenv("DB_SSL_MODE"))

	return nil
}

func IsProductionEnvironment() bool {
	return ENVIRONMENT == ENV_PRODUCTION
}

func IsSandboxEnvironment() bool {
	return ENVIRONMENT == ENV_SANDBOX
}

func IsDevelopmentEnvironment() bool {
	return ENVIRONMENT == ENV_DEVELOPMENT
}

func IsTestEnvironment() bool {
	return ENVIRONMENT == ENV_TEST
}

func IsCloudEnvironment() bool {
	return IsProductionEnvironment() || IsSandboxEnvironment()
}

func IsLocalEnvironment() bool {
	return IsDevelopmentEnvironment() || IsTestEnvironment()
}
