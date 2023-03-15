package sqlDB

import (
	"context"
	"database/sql"
	"errors"
)

// Statement struct
type Statement struct {
	ctx   context.Context
	query string
	args  []interface{}
}

// NewStatement create a new pointer to Statement struct
func NewStatement(ctx context.Context, query string, params ...interface{}) *Statement {
	return &Statement{ctx, query, params}
}

// Execute apply statement in database
func (s *Statement) Execute() error {
	if err := s.validate(); err != nil {
		return err
	}

	stmt, err := s.createStatement()
	if err != nil {
		return err
	}
	defer closer(stmt)

	if _, err = stmt.Exec(s.args...); err != nil {
		return err
	}

	return nil
}

func (s *Statement) validate() error {
	if instance == nil {
		return errors.New(db_not_initialized_error)
	}

	if s.query == "" {
		return errors.New(query_is_empty_error)
	}

	return nil
}

func (s *Statement) createStatement() (*sql.Stmt, error) {
	if tx := s.ctx.Value(SqlTxContext); tx != nil {
		return tx.(*sql.Tx).PrepareContext(s.ctx, s.query)
	}

	return instance.PrepareContext(s.ctx, s.query)
}
