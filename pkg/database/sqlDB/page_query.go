package sqlDB

import (
	"context"
	"errors"
	"fmt"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/types"
)

// PageQuery struct
type PageQuery[T any] struct {
	ctx   context.Context
	page  *types.PageRequest
	query string
	args  []interface{}
}

// NewPageQuery create a new pointer to PageQuery struct
func NewPageQuery[T any](ctx context.Context, page *types.PageRequest, query string, params ...interface{}) *PageQuery[T] {
	return &PageQuery[T]{ctx, page, query, params}
}

// Execute returns a pointer of page type with slice of T data
func (q *PageQuery[T]) Execute() (*types.Page[T], error) {
	if err := q.validate(); err != nil {
		return nil, err
	}

	var result types.Page[T]
	var err error
	result.TotalElements, err = q.pageTotal()
	if err != nil {
		return nil, err
	}

	result.Content, err = q.pageData()
	return &result, err
}

func (q *PageQuery[T]) pageTotal() (uint64, error) {
	query := fmt.Sprintf("SELECT COUNT(tb.*) FROM (%s) tb", q.query)

	var result uint64
	err := instance.QueryRowContext(q.ctx, query, q.args...).Scan(&result)
	return result, err
}

func (q *PageQuery[T]) pageData() ([]T, error) {
	query := fmt.Sprintf("%s ORDER BY %s LIMIT %d OFFSET %d", q.query, q.page.GetOrder(), q.page.Size, ((q.page.Page - 1) * q.page.Size))

	rows, err := instance.QueryContext(q.ctx, query, q.args...)
	if err != nil {
		return nil, err
	}
	defer closer(rows)

	return getDataList[T](rows)
}

func (q *PageQuery[T]) validate() error {
	if instance == nil {
		return errors.New(db_not_initialized_error)
	}

	if q.page == nil {
		return errors.New(page_is_empty_error)
	}

	if q.query == "" {
		return errors.New(query_is_empty_error)
	}

	return nil
}
