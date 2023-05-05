package security

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/cloud"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type firebaseAuthService struct {
	client *auth.Client
}

func newFirebaseAuthenticationService() *firebaseAuthService {
	auth, err := cloud.GetFirebaseSession().Auth(context.Background())
	if err != nil {
		logging.Fatal(connection_error, err)
	}

	return &firebaseAuthService{auth}
}

func (s *firebaseAuthService) GetUser(ctx context.Context, id string) (*User, error) {
	user, err := s.client.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    user.UID,
		Email: user.Email,
		Name:  user.DisplayName,
	}, nil
}

func (s *firebaseAuthService) CreateUser(ctx context.Context, user *UserCreate) error {
	password, err := user.decodeBase64Password()
	if err != nil {
		return err
	}

	userToCreate := (&auth.UserToCreate{}).
		UID(user.ID).
		Email(user.Email).
		Password(password).
		DisplayName(user.Name)

	if _, err = s.client.CreateUser(ctx, userToCreate); err != nil {
		return err
	}

	userToSetClaims := &auth.UserToUpdate{}
	userToSetClaims.CustomClaims(map[string]interface{}{
		profileField:  user.Profile,
		tenantIdField: user.TenantID,
	})
	_, err = s.client.UpdateUser(ctx, user.ID, userToSetClaims)

	return err
}

func (s *firebaseAuthService) UpdateUser(ctx context.Context, id string, user *UserUpdate) error {
	userToUpdate := &auth.UserToUpdate{}
	if user.Email != "" {
		userToUpdate.Email(user.Email)
	}

	if user.Password != "" {
		password, err := user.decodeBase64Password()
		if err != nil {
			return err
		}
		userToUpdate.Password(password)
	}

	if user.Name != "" {
		userToUpdate.DisplayName(user.Name)
	}

	if user.Profile != "" {
		userToUpdate.CustomClaims(map[string]interface{}{profileField: user.Profile})
	}

	_, err := s.client.UpdateUser(ctx, id, userToUpdate)
	return err
}

func (s *firebaseAuthService) DeleteUser(ctx context.Context, id string) error {
	return s.client.DeleteUser(ctx, id)
}

func (s *firebaseAuthService) EnableUser(ctx context.Context, id string) error {
	userToUpdate := (&auth.UserToUpdate{}).Disabled(false)
	_, err := s.client.UpdateUser(ctx, id, userToUpdate)
	return err
}

func (s *firebaseAuthService) DisableUser(ctx context.Context, id string) error {
	userToUpdate := (&auth.UserToUpdate{}).Disabled(true)
	_, err := s.client.UpdateUser(ctx, id, userToUpdate)
	return err
}
