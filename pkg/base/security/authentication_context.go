package security

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// IAuthenticationContext is the external user contract
type IAuthenticationContext interface {
	GetUserID() uuid.UUID
	GetTenantID() uuid.UUID
}

// AuthenticationContext is the external user struct
type AuthenticationContext struct {
	tenantID uuid.UUID
	userID   uuid.UUID
}

// NewAuthenticationContext returns a pointer of authentication context
func NewAuthenticationContext(tenantID, userID uuid.UUID) *AuthenticationContext {
	return &AuthenticationContext{tenantID, userID}
}

// GetTenantID returns the tenant id
func (a *AuthenticationContext) GetTenantID() uuid.UUID {
	return a.tenantID
}

// GetUserID returns the user id
func (a *AuthenticationContext) GetUserID() uuid.UUID {
	return a.userID
}

// SetInContext returns a context with authentication inside
func (a *AuthenticationContext) SetInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyAuthenticationContext, a)
}

// String returns the details of auth context
func (a *AuthenticationContext) String() string {
	return fmt.Sprintf("tenantId: %s | userId: %s", a.tenantID, a.userID)
}

// GetAuthenticationContext return the authentication context inside the context
func GetAuthenticationContext(ctx context.Context) *AuthenticationContext {
	if result := ctx.Value(contextKeyAuthenticationContext); result != nil {
		return result.(*AuthenticationContext)
	}
	return nil
}

// Valid returns a boolean if the context is valid
func (a *AuthenticationContext) Valid() bool {
	return (a.tenantID != uuid.Nil && a.userID != uuid.Nil)
}
