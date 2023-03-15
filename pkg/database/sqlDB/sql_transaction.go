package sqlDB

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/transaction"
)

type sqlTransaction struct {
	isolation sql.IsolationLevel
}

// NewTransaction creates new sqlTransaction with implements a transaction.Transaction
func NewTransaction(isolation ...sql.IsolationLevel) transaction.Transaction {
	isolationLevel := sql.LevelDefault
	if len(isolation) == 1 {
		isolationLevel = isolation[0]
	} else if len(isolation) > 1 {
		isolationLevel = isolation[0]
		logging.Warn("transaction isolation just use first parameter, others will be ignored")
	}

	return &sqlTransaction{isolation: isolationLevel}
}

// ExecTx execute a transactional sql
func (t *sqlTransaction) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, txCh, err := t.beginTx(ctx)
	if err != nil {
		return err
	}
	defer close(txCh)

	ctx = context.WithValue(ctx, transaction.SqlTxContext, tx)

	if err = fn(ctx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			fErr := fmt.Errorf("error when executing transaction rollback: %v: %w", err, rbErr)
			logging.Error("%v", fErr)
			txCh <- fErr
			return fErr
		}

		logging.Error("%v", err)
		txCh <- err
		return err
	}

	if err = tx.Commit(); err != nil {
		fErr := fmt.Errorf("could not commit transaction: %w", err)
		logging.Error("%v", fErr)
		txCh <- fErr
		return fErr
	}

	return nil
}

func (t *sqlTransaction) beginTx(ctx context.Context) (*sql.Tx, chan error, error) {
	tx, err := instance.BeginTx(ctx, &sql.TxOptions{Isolation: t.isolation})

	if err != nil {
		fErr := fmt.Errorf("could not start database transaction: %v", err)
		logging.Error("%v", fErr)
		return nil, nil, fErr
	}

	return tx, make(chan error, 1), nil
}
