package sqlDB

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	basePath := test.MountAbsolutPath(test.DATABASE_ENVIRONMENT_PATH)
	assert.NoError(t, os.Setenv(config.ENV_SQL_DB_MIGRATION, "true"))

	t.Run("should execute migration successfully and find users", func(t *testing.T) {
		assert.NoError(t, os.Setenv("MIGRATION_SOURCE_URL", fmt.Sprintf("%smigrations", test.DATABASE_ENVIRONMENT_PATH)))

		test.InitializeSqlDBTest()
		pc := test.UsePostgresContainer()

		Initialize()

		assert.NoError(t, pc.Dataset(basePath, "clear-database.sql", "add-users.sql"))

		const query = "SELECT u.id, u.name, u.birthday, p.id, p.name FROM users u INNER JOIN profiles p ON p.id = u.profile_id"
		result, err := NewQuery[User](context.Background(), query).Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})
}
