package sqlDB

import (
	"context"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/test"
	"github.com/stretchr/testify/assert"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/cacheDB"
)

func TestQueryWithoutInitialize(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/database/")

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	instance = nil

	t.Run("One without params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		_, err = NewQuery[User](context.Background(), query_base+" LIMIT 1").One()
		assert.Error(t, err, db_not_initialized_error)
	})

	t.Run("One with params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		_, err = NewQuery[User](context.Background(), query_base+" WHERE u.id = $1", 1).One()
		assert.Error(t, err, db_not_initialized_error)
	})

	t.Run("Many without params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		_, err = NewQuery[User](context.Background(), query_base).Many()
		assert.Error(t, err, db_not_initialized_error)
	})

	t.Run("Many with params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		_, err = NewQuery[User](context.Background(), query_base+" WHERE u.name = $1", "ADMIN USER").Many()
		assert.Error(t, err, db_not_initialized_error)
	})
}

func TestQuery(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/database/")

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	Initialize()

	t.Run("One without query", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		_, err = NewQuery[User](ctx, "").One()
		assert.Error(t, err, query_is_empty_error)

	})

	t.Run("One without params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		result, err := NewQuery[User](ctx, query_base+" LIMIT 1").One()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ADMIN USER", result.Name)

	})

	t.Run("One with params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		result, err := NewQuery[User](ctx, query_base+" WHERE u.name = $1", "ADMIN USER").One()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ADMIN USER", result.Name)

	})

	t.Run("Many without query", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		_, err = NewQuery[User](ctx, "").Many()
		assert.Error(t, err, query_is_empty_error)
	})

	t.Run("Many without params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		result, err := NewQuery[User](ctx, query_base).Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "ADMIN USER", result[0].Name)

	})

	t.Run("Many with params", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		result, err := NewQuery[User](ctx, query_base+" WHERE u.name = $1", "ADMIN USER").Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "ADMIN USER", result[0].Name)
	})
}

func TestQueryWithoutCacheDBInitialize(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/database/")

	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	Initialize()

	t.Run("One without params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsOneWithoutParams", time.Hour)
		result, err := NewCachedQuery(ctx, cache, query_base+" LIMIT 1").One()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ADMIN USER", result.Name)

		_, err = cache.One(ctx)
		assert.Error(t, err, "Cache not initialized")
	})

	t.Run("One with params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsOneWithoutParams", time.Hour)
		result, err := NewCachedQuery(ctx, cache, query_base+" WHERE u.name = $1", "ADMIN USER").One()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ADMIN USER", result.Name)

		_, err = cache.One(ctx)
		assert.Error(t, err, "Cache not initialized")
	})

	t.Run("Many without params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsManyWithoutParams", time.Hour)
		result, err := NewCachedQuery(ctx, cache, query_base).Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "ADMIN USER", result[0].Name)

		_, err = cache.Many(ctx)
		assert.Error(t, err, "Cache not initialized")
	})

	t.Run("Many with params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsManyWithoutParams", time.Hour)
		result, err := NewCachedQuery(ctx, cache, query_base+" WHERE u.name = $1", "ADMIN USER").Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "ADMIN USER", result[0].Name)

		_, err = cache.Many(ctx)
		assert.Error(t, err, "Cache not initialized")
	})
}

func TestCachedQuery(t *testing.T) {
	basePath := test.MountAbsolutPath("../../../development-environment/database/")

	test.InitializeCacheDBTest()
	test.InitializeSqlDBTest()
	pc := test.UsePostgresContainer()

	if err := pc.Dataset(basePath, "schema.sql"); err != nil {
		logging.Fatal(err.Error())
	}

	cacheDB.Initialize()
	Initialize()

	t.Run("One without query with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsOneWithoutQuery", time.Hour)
		cacheInitialData, err := cache.One(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, cacheInitialData)

		_, err = NewCachedQuery(ctx, cache, "").One()
		assert.Error(t, err, query_is_empty_error)
	})

	t.Run("One without params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsOneWithoutParams", time.Hour)
		cacheInitialData, err := cache.One(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, cacheInitialData)

		result, err := NewCachedQuery(ctx, cache, query_base+" LIMIT 1").One()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ADMIN USER", result.Name)

		cacheFinalData, err := cache.One(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, cacheFinalData)
		assert.Equal(t, "ADMIN USER", cacheFinalData.Name)

		err = cache.Del(ctx)
		assert.Nil(t, err)
	})

	t.Run("One with params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsOneWithoutParams", time.Hour)
		cacheInitialData, err := cache.One(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, cacheInitialData)

		result, err := NewCachedQuery(ctx, cache, query_base+" WHERE u.name = $1", "ADMIN USER").One()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ADMIN USER", result.Name)

		cacheFinalData, err := cache.One(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, cacheFinalData)
		assert.Equal(t, "ADMIN USER", cacheFinalData.Name)

		err = cache.Del(ctx)
		assert.Nil(t, err)
	})

	t.Run("Many without query with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsManyWithoutQuery", time.Hour)
		cacheInitialData, err := cache.One(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, cacheInitialData)

		_, err = NewCachedQuery(ctx, cache, "").Many()
		assert.Error(t, err, query_is_empty_error)
	})

	t.Run("Many without params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsManyWithoutParams", time.Hour)
		cacheInitialData, err := cache.One(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, cacheInitialData)

		result, err := NewCachedQuery(ctx, cache, query_base).Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "ADMIN USER", result[0].Name)

		cacheFinalData, err := cache.Many(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, cacheFinalData)
		assert.Len(t, result, 2)
		assert.Equal(t, "ADMIN USER", result[0].Name)

		err = cache.Del(ctx)
		assert.Nil(t, err)
	})

	t.Run("Many with params with cache", func(t *testing.T) {
		datasets := []string{"clear-database.sql", "add-users.sql"}
		err := pc.Dataset(basePath, datasets...)
		assert.NoError(t, err)

		ctx := context.Background()
		cache := cacheDB.NewCache[User]("DbTestsManyWithoutParams", time.Hour)
		cacheInitialData, err := cache.One(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, cacheInitialData)

		result, err := NewCachedQuery(ctx, cache, query_base+" WHERE u.name = $1", "ADMIN USER").Many()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "ADMIN USER", result[0].Name)

		cacheFinalData, err := cache.Many(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, cacheFinalData)
		assert.Len(t, result, 1)
		assert.Equal(t, "ADMIN USER", result[0].Name)

		err = cache.Del(ctx)
		assert.Nil(t, err)
	})
}
