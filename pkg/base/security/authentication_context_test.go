package security

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticationContext(t *testing.T) {
	var defaultUserId uuid.UUID = uuid.MustParse("5e859dae-c879-11eb-b8bc-0242ac130003")
	var defaultTenantId uuid.UUID = uuid.MustParse("5e859dae-c879-11eb-b8bc-0242ac130004")

	t.Run("Should return tenant and user", func(t *testing.T) {
		result := NewAuthenticationContext(defaultTenantId, defaultUserId)
		assert.NotNil(t, result)
		assert.Equal(t, defaultTenantId, result.GetTenantID())
		assert.Equal(t, defaultUserId, result.GetUserID())
	})

	t.Run("Should set in context", func(t *testing.T) {
		result := NewAuthenticationContext(defaultTenantId, defaultUserId).SetInContext(context.Background())
		assert.NotNil(t, result)
		assert.NotNil(t, result.Value(contextKeyAuthenticationContext))
	})

	t.Run("Should return nil when context is nil", func(t *testing.T) {
		result := GetAuthenticationContext(context.Background())
		assert.Nil(t, result)
	})

	t.Run("Should get in context", func(t *testing.T) {
		context := NewAuthenticationContext(defaultTenantId, defaultUserId).SetInContext(context.Background())
		assert.NotNil(t, context)

		result := GetAuthenticationContext(context)
		assert.Equal(t, defaultTenantId, result.GetTenantID())
		assert.Equal(t, defaultUserId, result.GetUserID())
	})
}
