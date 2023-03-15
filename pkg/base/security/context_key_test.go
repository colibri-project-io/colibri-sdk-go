package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	t.Run("Should return authentication context key", func(t *testing.T) {
		expected := "context key " + string(contextKeyAuthenticationContext)
		assert.Equal(t, expected, contextKeyAuthenticationContext.String())
	})
}
