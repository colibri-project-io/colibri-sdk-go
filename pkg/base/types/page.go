package types

import (
	"fmt"
	"strings"
)

type Page[T interface{}] struct {
	Content       []T    `json:"content"`
	TotalElements uint64 `json:"totalElements"`
}

type SortDirection string

const (
	ASC  SortDirection = "ASC"
	DESC SortDirection = "DESC"
)

type Sort struct {
	Direction SortDirection
	Field     string
}

func NewSort(direction SortDirection, field string) Sort {
	return Sort{direction, field}
}

type PageRequest struct {
	Page  uint16
	Size  uint16
	Order []Sort
}

func NewPageRequest(page uint16, size uint16, order []Sort) *PageRequest {
	return &PageRequest{page, size, order}
}

func (p *PageRequest) GetOrder() string {
	orders := make([]string, 0, len(p.Order))

	for _, order := range p.Order {
		orders = append(orders, fmt.Sprintf("%s %s", order.Field, order.Direction))
	}

	return strings.Join(orders, ", ")
}
