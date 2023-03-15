package types

import (
	"fmt"
	"strings"
)

// Page is the page response contract
type Page[T any] struct {
	Content       []T    `json:"content"`
	TotalElements uint64 `json:"totalElements"`
}

// SortDirection is the field sort direction
type SortDirection string

const (
	ASC  SortDirection = "ASC"
	DESC SortDirection = "DESC"
)

// Sort is the contract to sort an field
type Sort struct {
	Direction SortDirection
	Field     string
}

// NewSort returns a new Sort
func NewSort(direction SortDirection, field string) Sort {
	return Sort{direction, field}
}

// PageRequest is the contract of request page
type PageRequest struct {
	Page  uint16
	Size  uint16
	Order []Sort
}

// NewPageRequest returns a new page request pointer
func NewPageRequest(page uint16, size uint16, order []Sort) *PageRequest {
	return &PageRequest{page, size, order}
}

// GetOrder returns string contains concated order list
func (p *PageRequest) GetOrder() string {
	orders := make([]string, 0, len(p.Order))

	for _, order := range p.Order {
		orders = append(orders, fmt.Sprintf("%s %s", order.Field, order.Direction))
	}

	return strings.Join(orders, ", ")
}
