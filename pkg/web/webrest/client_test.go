package webrest

import (
	"context"
	"fmt"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

type userStruct struct {
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func InitWiremock() *RestClient {
	wiremockContainer := test.UseWiremockContainer(test.MountAbsolutPath("../../../development-environment/wiremock/"))

	monitoring.Initialize()
	return NewRestClient("test-rest-client", fmt.Sprintf("http://localhost:%d/users-api/v1", wiremockContainer.Port()), 1)
}

func TestGet(t *testing.T) {
	restClient := InitWiremock()

	t.Run("Should return users", func(t *testing.T) {
		users, err := Get[[]userStruct](context.Background(), restClient, "", nil)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, *users, 5)
	})

	t.Run("Should return one user", func(t *testing.T) {
		user, err := Get[userStruct](context.Background(), restClient, "/1", nil)
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestPost(t *testing.T) {
	restClient := InitWiremock()

	t.Run("Should post a new user", func(t *testing.T) {
		newUser := userStruct{Name: "User 10", Email: "user_10@email.com"}
		savedUser, err := Post[userStruct](context.Background(), restClient, "/users", &newUser, nil)
		assert.NoError(t, err)
		assert.NotNil(t, savedUser)
		assert.Equal(t, uint(10), savedUser.ID)
	})
}

func TestPut(t *testing.T) {
	restClient := InitWiremock()

	t.Run("Should put a user", func(t *testing.T) {
		newUser := userStruct{ID: 10, Name: "User 10 edited", Email: "user_10@email.com"}
		savedUser, err := Put[userStruct](context.Background(), restClient, fmt.Sprintf("/users/%d", newUser.ID), &newUser, nil)
		assert.NoError(t, err)
		assert.NotNil(t, savedUser)
		assert.Equal(t, uint(10), savedUser.ID)
		assert.Equal(t, "User 10 edited", savedUser.Name)
	})
}

func TestDelete(t *testing.T) {
	restClient := InitWiremock()

	t.Run("Should delete a user with response body", func(t *testing.T) {
		deletedUser, err := Delete[userStruct](context.Background(), restClient, "/users/11", nil)
		assert.NoError(t, err)
		assert.NotNil(t, deletedUser)
		assert.Equal(t, uint(11), deletedUser.ID)
		assert.Equal(t, "User 11 deleted", deletedUser.Name)
	})

	t.Run("Should delete a user with empty response body", func(t *testing.T) {
		deletedUser, err := Delete[userStruct](context.Background(), restClient, "/users/12", nil)
		assert.NoError(t, err)
		assert.Nil(t, deletedUser)
	})
}
