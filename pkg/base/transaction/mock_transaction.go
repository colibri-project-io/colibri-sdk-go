package transaction

import "context"

// mockSqlTransaction implements a transaction.Transaction
type mockSqlTransaction struct {
}

// NewMockTransaction returns a new mockSqlTransaction.
//
// No parameters.
// Returns a Transaction.
func NewMockTransaction() Transaction {
	return mockSqlTransaction{}
}

// Execute executes a mockSqlTransaction.
//
// ctx: The context for the transaction.
// fn: The function to be executed.
// Returns an error.
func (m mockSqlTransaction) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
