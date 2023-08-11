package security

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

const (
	connection_error = "An error occurred when trying to connect to the authentication service. Error: %s"
)

type authenticationService interface {
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *UserCreate) error
	UpdateUser(ctx context.Context, id string, user *UserUpdate) error
	DeleteUser(ctx context.Context, id string) error
	EnableUser(ctx context.Context, id string) error
	DisableUser(ctx context.Context, id string) error
}

var instance authenticationService

func InitializeAuthenticationService() {
	switch config.CLOUD {
	case config.CLOUD_FIREBASE:
		instance = newFirebaseAuthenticationService()
	}

	logging.Info("Authentication service connected")
}

func GetUser(ctx context.Context, id string) (*User, error) {
	return instance.GetUser(ctx, id)
}

func CreateUser(ctx context.Context, user *UserCreate) error {
	return instance.CreateUser(ctx, user)
}

func UpdateUser(ctx context.Context, id string, user *UserUpdate) error {
	return instance.UpdateUser(ctx, id, user)
}

func DeleteUser(ctx context.Context, id string) error {
	return instance.DeleteUser(ctx, id)
}

func EnableUser(ctx context.Context, id string) error {
	return instance.EnableUser(ctx, id)
}

func DisableUser(ctx context.Context, id string) error {
	return instance.DisableUser(ctx, id)
}
