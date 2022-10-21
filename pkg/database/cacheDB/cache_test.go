package cacheDB

import (
	"context"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

type userCached struct {
	Id   int
	Name string
}

func TestCacheNotInitializedValidation(t *testing.T) {
	t.Run("Should return error when cacheDB is not initialized", func(t *testing.T) {
		instance = nil
		cache := NewCache[userCached]("cache-test", time.Hour)

		_, err := cache.Many(context.Background())
		assert.NotNil(t, err)
	})
}

func TestCache(t *testing.T) {
	expected := []userCached{
		{Id: 1, Name: "User 1"},
		{Id: 2, Name: "User 2"},
		{Id: 3, Name: "User 3"},
	}
	test.InitializeCacheDBTest()
	Initialize()

	t.Run("Should return a new point of cache", func(t *testing.T) {
		result := NewCache[userCached]("cache-test", time.Hour)
		assert.NotNil(t, result)
		assert.NotNil(t, result.name)
		assert.NotNil(t, result.ttl)
	})

	t.Run("Should return error when cache name is empty value", func(t *testing.T) {
		cache := Cache[userCached]{}

		_, err := cache.Many(context.Background())
		assert.NotNil(t, err)
	})

	t.Run("Should set many data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)

		err := cache.Set(context.Background(), expected)
		assert.Nil(t, err)

		result, err := cache.Many(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, expected, result)
	})

	t.Run("Should set one data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)

		err := cache.Set(context.Background(), expected[0])
		assert.Nil(t, err)

		result, err := cache.One(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected[0], *result)
	})

	t.Run("Should del data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)

		err := cache.Set(context.Background(), expected)
		assert.Nil(t, err)

		result, err := cache.Many(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, expected, result)

		err = cache.Del(context.Background())
		assert.Nil(t, err)

		result, err = cache.Many(context.Background())
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
