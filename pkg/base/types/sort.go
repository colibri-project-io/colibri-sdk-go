package types

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
