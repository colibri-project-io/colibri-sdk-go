package sqlDB

import (
	"context"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMigrations(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/database/")
	assert.NoError(t, os.Setenv("EXEC_MIGRATION", "true"))

	t.Run("should execute migration successfully and find users", func(t *testing.T) {
		assert.NoError(t, os.Setenv("MIGRATION_SOURCE_URL", "../../../development-environment/database/migrations"))

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
