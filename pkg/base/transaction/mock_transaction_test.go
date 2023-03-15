package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockTransaction(t *testing.T) {
	t.Run("execute no error", func(t *testing.T) {
		f := func(ctx context.Context) error { return nil }
		m := NewMockTransaction()

		assert.NoError(t, m.ExecTx(context.Background(), f))
	})

	t.Run("execute with error", func(t *testing.T) {
		expectedErr := errors.New("could not execute")
		f := func(ctx context.Context) error {
			return fmt.Errorf("the force is weak in you: %w", expectedErr)
		}
		m := NewMockTransaction()

		assert.ErrorIs(t, m.ExecTx(context.Background(), f), expectedErr)
	})
}
