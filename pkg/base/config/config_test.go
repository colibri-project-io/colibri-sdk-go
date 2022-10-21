package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	invalid_value           = "XYZ"
	app_name_value          = "TEST APP NAME"
	app_type_value          = "service"
	new_relic_license_value = "123456-abcdef-7890-ghi"
	cloud_value             = "aws"
	cloud_host_value        = "http://my-cloud-host-fake.com"
	cloud_region_value      = "test-region"
	cloud_secret_value      = "test-secret"
	cloud_token_value       = "test-token"
	cloud_disable_ssl_value = "true"
	port_value              = "8081"
	cache_uri_value         = "my-cache-fake:6379"
	db_value                = "sql"
	db_name_value           = "my-db-name"
	db_host_value           = "my-db-host"
	db_port_value           = "1234"
	db_user_value           = "my-db-user"
	db_password_value       = "my-db-password"
	db_ssl_mode_value       = "disable"
)

func TestLoad(t *testing.T) {
	t.Run("Should return error when enviroment is not configured", func(t *testing.T) {
		assert.EqualError(t, Load(), error_enviroment_not_configured)
	})

	t.Run("Should return error when enviroment contains a invalid value", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", invalid_value)

		err := Load()

		assert.Equal(t, ENVIRONMENT, invalid_value)
		assert.EqualError(t, err, error_enviroment_not_configured)
	})

	t.Run("Should return error when app name is not configured", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		err := Load()

		assert.Equal(t, ENVIRONMENT, ENV_PRODUCTION)
		assert.EqualError(t, err, error_app_name_not_configured)
	})

	t.Run("Should return error when app_type is not configured", func(t *testing.T) {
		os.Setenv("APP_NAME", app_name_value)
		err := Load()

		assert.EqualError(t, err, error_app_type_not_configured)
	})

	t.Run("Should return error when app_type contains a invalid value", func(t *testing.T) {
		os.Setenv("APP_TYPE", invalid_value)

		err := Load()

		assert.Equal(t, APP_TYPE, invalid_value)
		assert.EqualError(t, err, error_app_type_not_configured)
	})

	t.Run("Should return error when cloud is not configured", func(t *testing.T) {
		os.Setenv("APP_TYPE", app_type_value)
		err := Load()

		assert.EqualError(t, err, error_cloud_not_configured)
	})

	t.Run("Should return error when cloud contains a invalid value", func(t *testing.T) {
		os.Setenv("CLOUD", invalid_value)

		err := Load()

		assert.EqualError(t, err, error_cloud_not_configured)
	})

	t.Run("Should return error when production environment contains required params is not configured", func(t *testing.T) {
		os.Setenv("CLOUD", cloud_value)

		err := Load()

		assert.EqualError(t, err, error_production_required_params_not_configured)
	})

	t.Run("Should return error when db contains a invalid value", func(t *testing.T) {
		os.Setenv("NEW_RELIC_LICENSE", new_relic_license_value)
		os.Setenv("DB", invalid_value)

		err := Load()

		assert.EqualError(t, err, error_database_misconfigured)
	})

	t.Run("Should load configurations with success and return nil error", func(t *testing.T) {
		os.Setenv("DB", db_value)
		os.Setenv("CLOUD_HOST", cloud_host_value)
		os.Setenv("CLOUD_REGION", cloud_region_value)
		os.Setenv("CLOUD_SECRET", cloud_secret_value)
		os.Setenv("CLOUD_TOKEN", cloud_token_value)
		os.Setenv("CLOUD_DISABLE_SSL", cloud_disable_ssl_value)
		os.Setenv("PORT", port_value)
		os.Setenv("CACHE_URI", cache_uri_value)
		os.Setenv("DB_HOST", db_host_value)
		os.Setenv("DB_PORT", db_port_value)
		os.Setenv("DB_USER", db_user_value)
		os.Setenv("DB_NAME", db_name_value)
		os.Setenv("DB_PASSWORD", db_password_value)
		os.Setenv("DB_SSL_MODE", db_ssl_mode_value)

		cloudDisableSsl, _ := strconv.ParseBool(cloud_disable_ssl_value)
		serverPort, _ := strconv.Atoi(port_value)
		dbConnectionUri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name='%s' sslmode=%s",
			db_host_value,
			db_port_value,
			db_user_value,
			db_password_value,
			db_name_value,
			app_name_value,
			db_ssl_mode_value)

		err := Load()

		assert.Equal(t, ENVIRONMENT, ENV_PRODUCTION)
		assert.Equal(t, APP_NAME, app_name_value)
		assert.Equal(t, APP_TYPE, APP_TYPE_SERVICE)
		assert.Equal(t, CLOUD, cloud_value)
		assert.Equal(t, NEW_RELIC_LICENSE, new_relic_license_value)
		assert.Equal(t, CLOUD_HOST, cloud_host_value)
		assert.Equal(t, CLOUD_REGION, cloud_region_value)
		assert.Equal(t, CLOUD_SECRET, cloud_secret_value)
		assert.Equal(t, CLOUD_TOKEN, cloud_token_value)
		assert.Equal(t, CLOUD_DISABLE_SSL, cloudDisableSsl)
		assert.Equal(t, PORT, serverPort)
		assert.Equal(t, CACHE_URI, cache_uri_value)
		assert.Equal(t, DB_CONNECTION_URI, dbConnectionUri)
		assert.Nil(t, err)
	})
}

func TestIsProductionEnviroment(t *testing.T) {
	t.Run("Should return false when enviroment is not production", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_TEST)

		Load()

		assert.Equal(t, IsProductionEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not production", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_DEVELOPMENT)

		Load()

		assert.Equal(t, IsProductionEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not production", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_SANDBOX)

		Load()

		assert.Equal(t, IsProductionEnvironment(), false)
	})

	t.Run("Should return true when enviroment is production", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		Load()

		assert.Equal(t, IsProductionEnvironment(), true)
	})
}

func TestIsSandboxEnviroment(t *testing.T) {
	t.Run("Should return false when enviroment is not sandbox", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		Load()

		assert.Equal(t, IsSandboxEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not sandbox", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_TEST)

		Load()

		assert.Equal(t, IsSandboxEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not sandbox", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_DEVELOPMENT)

		Load()

		assert.Equal(t, IsSandboxEnvironment(), false)
	})

	t.Run("Should return true when enviroment is sandbox", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_SANDBOX)

		Load()

		assert.Equal(t, IsSandboxEnvironment(), true)
	})
}

func TestIsDevelopmentEnviroment(t *testing.T) {
	t.Run("Should return false when enviroment is not development", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_TEST)

		Load()

		assert.Equal(t, IsDevelopmentEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not development", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		Load()

		assert.Equal(t, IsDevelopmentEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not development", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_SANDBOX)

		Load()

		assert.Equal(t, IsDevelopmentEnvironment(), false)
	})

	t.Run("Should return true when enviroment is development", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_DEVELOPMENT)

		Load()

		assert.Equal(t, IsDevelopmentEnvironment(), true)
	})
}

func TestIsTestEnviroment(t *testing.T) {
	t.Run("Should return false when enviroment is not test", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_DEVELOPMENT)

		Load()

		assert.Equal(t, IsTestEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not test", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_SANDBOX)

		Load()

		assert.Equal(t, IsTestEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not test", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		Load()

		assert.Equal(t, IsTestEnvironment(), false)
	})

	t.Run("Should return true when enviroment is test", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_TEST)

		Load()

		assert.Equal(t, IsTestEnvironment(), true)
	})
}

func TestIsCloudEnviroment(t *testing.T) {
	t.Run("Should return false when enviroment is not cloud environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_DEVELOPMENT)

		Load()

		assert.Equal(t, IsCloudEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not cloud environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_TEST)

		Load()

		assert.Equal(t, IsCloudEnvironment(), false)
	})

	t.Run("Should return true when enviroment is cloud environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		Load()

		assert.Equal(t, IsCloudEnvironment(), true)
	})

	t.Run("Should return true when enviroment is cloud environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_SANDBOX)

		Load()

		assert.Equal(t, IsCloudEnvironment(), true)
	})
}

func TestIsLocalEnvironment(t *testing.T) {
	t.Run("Should return false when enviroment is not local environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_PRODUCTION)

		Load()

		assert.Equal(t, IsLocalEnvironment(), false)
	})

	t.Run("Should return false when enviroment is not local environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_SANDBOX)

		Load()

		assert.Equal(t, IsLocalEnvironment(), false)
	})

	t.Run("Should return true when enviroment is local environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_DEVELOPMENT)

		Load()

		assert.Equal(t, IsLocalEnvironment(), true)
	})

	t.Run("Should return true when enviroment is local environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", ENV_TEST)

		Load()

		assert.Equal(t, IsLocalEnvironment(), true)
	})
}
