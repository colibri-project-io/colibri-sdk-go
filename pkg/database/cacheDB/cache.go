package cacheDB

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
)

// Cache struct
type Cache[T interface{}] struct {
	name string
	ttl  time.Duration
}

// NewCache create a new pointer to Cache struct.
func NewCache[T interface{}](name string, ttl time.Duration) *Cache[T] {
	return &Cache[T]{name, ttl}
}

// Many returns a slice of T value
func (c *Cache[T]) Many(ctx context.Context) ([]T, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	result, err := instance.Get(ctx, c.getNamePrefixed()).Bytes()
	if err != nil {
		return nil, err
	}

	list := make([]T, 0)
	if err = json.Unmarshal(result, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// One return a pointer of T value
func (c *Cache[T]) One(ctx context.Context) (*T, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	result, err := instance.Get(ctx, c.getNamePrefixed()).Bytes()
	if err != nil {
		return nil, err
	}

	model := new(T)
	if err = json.Unmarshal(result, &model); err != nil {
		return nil, err
	}

	return model, nil
}

// Set save data in cacheDB
func (c *Cache[T]) Set(ctx context.Context, data interface{}) error {
	if err := c.validate(); err != nil {
		return err
	}

	jsonData, _ := json.Marshal(data)
	return instance.Set(ctx, c.getNamePrefixed(), jsonData, c.ttl).Err()
}

// Del delete data in cachedDB
func (c *Cache[T]) Del(ctx context.Context) error {
	if err := c.validate(); err != nil {
		return err
	}

	return instance.Del(ctx, c.getNamePrefixed()).Err()
}

func (c *Cache[T]) validate() error {
	if instance == nil {
		return errors.New("Cache not initialized")
	}

	if c.name == "" {
		return errors.New("Cache without name")
	}

	return nil
}

func (c *Cache[T]) getNamePrefixed() string {
	return fmt.Sprintf("%s::%s", config.APP_NAME, c.name)
}
