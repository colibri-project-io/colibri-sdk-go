package transaction

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockTransaction(t *testing.T) {
	t.Run("execute no error", func(t *testing.T) {
		f := func(ctx context.Context) error { return nil }
		m := NewMockTransaction()

		assert.NoError(t, m.Execute(context.Background(), f))
	})

	t.Run("execute with error", func(t *testing.T) {
		expectedErr := errors.New("could not execute")
		f := func(ctx context.Context) error {
			return fmt.Errorf("the force is weak in you: %w", expectedErr)
		}
		m := NewMockTransaction()

		assert.ErrorIs(t, m.Execute(context.Background(), f), expectedErr)
	})
}
