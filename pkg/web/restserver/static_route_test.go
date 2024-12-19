package restserver

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
)

func TestStartStatisRoutes(t *testing.T) {
	test.InitializeBaseTest()

	AddStaticRoute(StaticRoute{
		Path: "../../../development-environment/rest",
		URI:  "/",
	})

	l := listener(t)
	config.PORT = l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	baseURL := fmt.Sprintf("http://localhost:%d/", config.PORT)

	go ListenAndServe()

	t.Run("should return the json file", func(t *testing.T) {
		fmt.Println(baseURL)
		response, err := http.Get(baseURL)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		body, err := io.ReadAll(response.Body)

		assert.NoError(t, err)
		assert.Contains(t, string(body), "minha-cozinha")
	})
}

func listener(t *testing.T) (l net.Listener) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Error("could not get available port")
	}

	return l
}
