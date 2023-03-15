package sqlDB

import (
	"context"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
	"github.com/stretchr/testify/assert"
)

func TestPageQueryWithoutInitialize(t *testing.T) {
	basePath := test.MountAbsolutPath(test.DATABASE_ENVIRONMENT_PATH)

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	instance = nil

	t.Run("Page", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		orders := []types.Sort{
			{Direction: types.DESC, Field: "name"},
			{Direction: types.ASC, Field: "birthday"},
		}
		page := types.NewPageRequest(1, 1, orders)
		_, err = NewPageQuery[User](context.Background(), page, query_base).Execute()
		assert.Error(t, err, db_not_initialized_error)
	})
}

func TestPageQuery(t *testing.T) {
	ctx := context.Background()
	orders := []types.Sort{
		{Direction: types.DESC, Field: "u.name"},
		{Direction: types.ASC, Field: "u.birthday"},
	}
	page := types.NewPageRequest(1, 1, orders)

	basePath := test.MountAbsolutPath(test.DATABASE_ENVIRONMENT_PATH)
	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()
	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	Initialize()

	t.Run("Page without page info", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		_, err = NewPageQuery[User](ctx, nil, query_base).Execute()
		assert.Error(t, err, page_is_empty_error)
	})

	t.Run("Page without query", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		_, err = NewPageQuery[User](ctx, page, "").Execute()
		assert.Error(t, err, query_is_empty_error)
	})

	t.Run("Page", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		result, err := NewPageQuery[User](ctx, page, query_base).Execute()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "OTHER USER", result.Content[0].Name)
		assert.Equal(t, uint64(2), result.TotalElements)
	})
}
