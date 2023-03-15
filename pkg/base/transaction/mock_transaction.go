package transaction

import "context"

// mockSqlTransaction struct with implements Transaction for tests purpose
type mockSqlTransaction struct {
}

// NewMockTransaction create a mock transaction for tests purpose
func NewMockTransaction() Transaction {
	return mockSqlTransaction{}
}

func (m mockSqlTransaction) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
