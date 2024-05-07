package cacheDB

import (
	"context"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
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
	ctx := context.Background()
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

	t.Run("Should return error when cache name is empty value on call many", func(t *testing.T) {
		cache := Cache[userCached]{}

		_, err := cache.Many(ctx)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when cache name is empty value on call one", func(t *testing.T) {
		cache := Cache[userCached]{}

		_, err := cache.One(ctx)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when cache name is empty value on call set", func(t *testing.T) {
		cache := Cache[userCached]{}

		err := cache.Set(ctx, cache)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when cache name is empty value on call del", func(t *testing.T) {
		cache := Cache[userCached]{}

		err := cache.Del(ctx)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when occurred error in json unmarshal on set data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)
		invalid := map[string]interface{}{
			"invalid": make(chan int),
		}

		err := cache.Set(ctx, invalid)
		assert.NotNil(t, err)
	})

	t.Run("Should set many data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)

		err := cache.Set(ctx, expected)
		assert.Nil(t, err)

		result, err := cache.Many(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, expected, result)
	})

	t.Run("Should set one data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)

		err := cache.Set(ctx, expected[0])
		assert.Nil(t, err)

		result, err := cache.One(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected[0], *result)
	})

	t.Run("Should del data in cache", func(t *testing.T) {
		cache := NewCache[userCached]("cache-test", time.Hour)

		err := cache.Set(ctx, expected)
		assert.Nil(t, err)

		result, err := cache.Many(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
		assert.Equal(t, expected, result)

		err = cache.Del(ctx)
		assert.Nil(t, err)

		result, err = cache.Many(ctx)
		assert.Nil(t, err)
		assert.Nil(t, result)
	})
}
