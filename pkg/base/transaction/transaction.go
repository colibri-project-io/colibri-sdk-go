package transaction

import (
	"context"
)

const SqlTxContext = "SqlTxContext"

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}
