package sqlDB

import (
	"context"
	"database/sql"
	"errors"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/transaction"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/database/cacheDB"
)

// Query struct
type Query[T any] struct {
	ctx   context.Context
	cache *cacheDB.Cache[T]
	query string
	args  []any
}

// NewQuery create a new pointer to Query struct
func NewQuery[T any](ctx context.Context, query string, params ...any) *Query[T] {
	return &Query[T]{ctx, nil, query, params}
}

// NewCachedQuery create a new pointer to Query struct with cache
func NewCachedQuery[T any](ctx context.Context, cache *cacheDB.Cache[T], query string, params ...any) (q *Query[T]) {
	return &Query[T]{ctx, cache, query, params}
}

// Many returns a slice of T value
func (q *Query[T]) Many() ([]T, error) {
	if err := q.validate(); err != nil {
		return nil, err
	}

	if q.cache == nil {
		return q.fetchMany()
	}

	result, err := q.cache.Many(q.ctx)
	if err != nil {
		return q.fetchMany()
	}
	return result, nil
}

func (q *Query[T]) fetchMany() ([]T, error) {
	rows, err := q.queryContext()
	if err != nil {
		return nil, err
	}
	defer closer(rows)

	list, err := getDataList[T](rows)
	if err != nil {
		return nil, err
	}

	if q.cache != nil {
		q.cache.Set(q.ctx, list)
	}

	return list, nil
}

// One return a pointer of T value
func (q *Query[T]) One() (*T, error) {
	if err := q.validate(); err != nil {
		return nil, err
	}

	if q.cache == nil {
		return q.fetchOne()
	}

	result, err := q.cache.One(q.ctx)
	if err != nil {
		return q.fetchOne()
	}
	return result, nil
}

func (q *Query[T]) fetchOne() (*T, error) {
	model := new(T)
	if err := q.queryRowContext().Scan(reflectCols(model)...); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	if q.cache != nil {
		q.cache.Set(q.ctx, model)
	}

	return model, nil
}

func (q *Query[T]) validate() error {
	if instance == nil {
		return errors.New(db_not_initialized_error)
	}

	if q.query == "" {
		return errors.New(query_is_empty_error)
	}

	return nil
}

func (q *Query[T]) queryContext() (*sql.Rows, error) {
	if tx := q.ctx.Value(transaction.SqlTxContext); tx != nil {
		return tx.(*sql.Tx).QueryContext(q.ctx, q.query, q.args...)
	}

	return instance.QueryContext(q.ctx, q.query, q.args...)
}

func (q *Query[T]) queryRowContext() *sql.Row {
	if tx := q.ctx.Value(transaction.SqlTxContext); tx != nil {
		return tx.(*sql.Tx).QueryRowContext(q.ctx, q.query, q.args...)
	}

	return instance.QueryRowContext(q.ctx, q.query, q.args...)
}
