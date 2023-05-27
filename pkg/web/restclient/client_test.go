package restclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/net"
)

type userResponseTestStruct struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponseErrorTestStruct struct {
	Message string `json:"message"`
}

type userRequestTestStruct struct {
	ID    uint   `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func InitWiremock() *RestClient {
	wiremockContainer := test.UseWiremockContainer(test.MountAbsolutPath(test.WIREMOCK_ENVIRONMENT_PATH))

	test.InitializeBaseTest()
	return NewRestClient("test-client", fmt.Sprintf("http://localhost:%d/users-api/v1", wiremockContainer.Port()), 1)
}

func TestGet(t *testing.T) {
	restClient := InitWiremock()

	t.Run("User ok", func(t *testing.T) {
		t.Parallel()
		resp := Get[userResponseTestStruct](context.Background(), restClient, "/1", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body(), "user is nil")
	})

	t.Run("User with timeout", func(t *testing.T) {
		t.Parallel()
		resp := Get[userResponseTestStruct](context.Background(), restClient, "/2", nil)
		assert.Errorf(t, resp.Error(), "expected timeout error: %v\n", resp.Error())
		assert.Truef(t, net.IsTimeout(resp.Error()), "expected timeout error: %v\n", resp.Error())
		assert.Nilf(t, resp.Body(), "user is not null\n")
	})

	t.Run("List users ok", func(t *testing.T) {
		t.Parallel()
		resp := Get[[]userResponseTestStruct](context.Background(), restClient, "", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.True(t, len(*resp.Body()) == 5)
	})
}

func TestPost(t *testing.T) {
	restClient := InitWiremock()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		newUser := userRequestTestStruct{Name: "User 10", Email: "user_10@email.com"}
		resp := Post[userResponseTestStruct](context.Background(), restClient, "/users", &newUser, nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, uint(10), resp.Body().ID)
	})
}

func TestPut(t *testing.T) {
	restClient := InitWiremock()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		newUser := userRequestTestStruct{ID: 10, Name: "User 10 edited", Email: "user_10@email.com"}
		resp := Put[userResponseTestStruct](context.Background(), restClient, fmt.Sprintf("/users/%d", newUser.ID), &newUser, nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, uint(10), resp.Body().ID)
		assert.Equal(t, "User 10 edited", resp.Body().Name)
	})
}

func TestPatch(t *testing.T) {
	restClient := InitWiremock()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		newUser := userRequestTestStruct{Name: "User 10 edited"}
		resp := Patch[userResponseTestStruct](context.Background(), restClient, "/users/10", &newUser, nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, uint(10), resp.Body().ID)
		assert.Equal(t, "User 10 edited", resp.Body().Name)
	})
}

func TestDelete(t *testing.T) {
	restClient := InitWiremock()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		resp := Delete[userResponseTestStruct](context.Background(), restClient, "/users/11", nil)
		assert.NotNil(t, resp)
		assert.NoError(t, resp.Error())
		assert.NotNil(t, resp.Body())
		assert.Equal(t, uint(11), resp.Body().ID)
		assert.Equal(t, "User 11 deleted", resp.Body().Name)
	})

	t.Run("OK no content", func(t *testing.T) {
		t.Parallel()
		resp := Delete[userResponseTestStruct](context.Background(), restClient, "/users/12", nil)
		assert.NoError(t, resp.Error())
		assert.Nil(t, resp.Body())
	})
}

func TestPostNotEmptyResponseBodyError(t *testing.T) {
	restClient := InitWiremock()

	t.Run("Should return an empty body", func(t *testing.T) {
		t.Parallel()
		newUser := userRequestTestStruct{Name: "Empty Body Response", Email: "post_user_empty_body@error.com"}

		resp, respErr := PostWithErrorData[any, userResponseErrorTestStruct](context.Background(), restClient, "/users", &newUser, nil)

		assert.NotNil(t, resp)
		assert.Nil(t, resp.Body())
		assert.ErrorContains(t, resp.Error(), "500")
		assert.Nil(t, respErr)
	})

	t.Run("Should return a body on statusCode 500", func(t *testing.T) {
		t.Parallel()
		newUser := userRequestTestStruct{Name: "Body Response", Email: "post_user_with_body@error.com"}

		resp, respErr := PostWithErrorData[userResponseTestStruct, userResponseErrorTestStruct](context.Background(), restClient, "/users", &newUser, nil)

		assert.NotNil(t, resp)
		assert.Nil(t, resp.Body())
		assert.ErrorContains(t, resp.Error(), "500")
		assert.Equal(t, respErr.Message, "Error message post user")
	})

	t.Run("Should return an decode error on statusCode 500", func(t *testing.T) {
		t.Parallel()
		newUser := userRequestTestStruct{Name: "Body Response", Email: "post_user_decode_error@error.com"}

		resp, respErr := PostWithErrorData[userResponseTestStruct, userResponseErrorTestStruct](context.Background(), restClient, "/users", &newUser, nil)

		assert.NotNil(t, resp)
		assert.Nil(t, resp.Body())
		assert.ErrorContains(t, resp.Error(), "500")
		assert.ErrorContains(t, resp.Error(), "could not decode response")
		assert.Nil(t, respErr)
	})
}

func TestPostWithBodyString(t *testing.T) {
	restClient := InitWiremock()

	t.Run("post x-www-form-urlencoded", func(t *testing.T) {
		data := url.Values{}
		data.Set("user", "darth_vader")
		data.Set("pass", "force")
		resp := PostBodyString[userResponseTestStruct](context.Background(), restClient, "/login", data.Encode(), map[string]string{"Content-type": "application/x-www-form-urlencoded"})

		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Body())
		assert.Equal(t, uint(100), resp.Body().ID)
		assert.Equal(t, "user_100@email.com", resp.Body().Email)
		assert.Equal(t, "User 100", resp.Body().Name)
	})
}

func TestCreateSegmentTestEnv(t *testing.T) {
	monitoring.Initialize()
	seg := createSegment(context.Background(), nil, http.MethodPost, "/api/users")
	assert.Nil(t, seg)
}
