package sqlDB

import (
	"context"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/stretchr/testify/assert"
)

func TestStatementWithoutInitialize(t *testing.T) {
	basePath := test.MountAbsolutPath(test.DATABASE_ENVIRONMENT_PATH)

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	instance = nil

	t.Run("Statement", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		birth, _ := time.Parse("2006-01-02", "2021-11-22")
		user := User{123, "Usuário teste stmt", birth, Profile{100, "ADMIN"}}
		err = NewStatement(ctx, "INSERT INTO users VALUES ($1, $2, $3, $4)", user.Id, user.Name, user.Birthday, user.Profile.Id).Execute()
		assert.Error(t, err, db_not_initialized_error)
	})
}

func TestStatement(t *testing.T) {
	basePath := test.MountAbsolutPath(test.DATABASE_ENVIRONMENT_PATH)

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	Initialize()

	t.Run("Statement without query", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		err = NewStatement(ctx, "").Execute()
		assert.Error(t, err, query_is_empty_error)
	})

	t.Run("Statement", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		birth, _ := time.Parse("2006-01-02", "2021-11-22")
		user := User{123, "Usuário teste stmt", birth, Profile{100, "ADMIN"}}
		err = NewStatement(ctx, "INSERT INTO users VALUES ($1, $2, $3, $4)", user.Id, user.Name, user.Birthday, user.Profile.Id).Execute()
		assert.NoError(t, err)

		result, err := NewQuery[User](ctx, query_base+" WHERE u.id = $1", user.Id).One()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Id, result.Id)
		assert.Equal(t, user.Name, result.Name)
		assert.Equal(t, user.Birthday.Local(), result.Birthday.Local())
		assert.Equal(t, user.Profile, result.Profile)
	})
}
