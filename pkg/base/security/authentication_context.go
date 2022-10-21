package security

import (
	"context"

	"github.com/google/uuid"
)

type IAuthenticationContext interface {
	GetUserID() uuid.UUID
	GetTenantID() uuid.UUID
}

type AuthenticationContext struct {
	tenantID uuid.UUID
	userID   uuid.UUID
}

func NewAuthenticationContext(tenantID, userID uuid.UUID) *AuthenticationContext {
	return &AuthenticationContext{tenantID, userID}
}

func (a *AuthenticationContext) GetTenantID() uuid.UUID {
	return a.tenantID
}

func (a *AuthenticationContext) GetUserID() uuid.UUID {
	return a.userID
}

func (a *AuthenticationContext) SetInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, AuthenticationContext{}, a)
}

func GetAuthenticationContext(ctx context.Context) *AuthenticationContext {
	if result := ctx.Value(AuthenticationContext{}); result != nil {
		return result.(*AuthenticationContext)
	}
	return nil
}
