package transaction

import (
	"context"
)

// Transaction defines the interface for a transaction
type Transaction interface {
	Execute(context.Context, func(ctx context.Context) error) error
}
