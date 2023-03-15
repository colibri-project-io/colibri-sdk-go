package restserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseServer(t *testing.T) {
	srv = &GorillaServer{srv: &http.Server{}}

	restObserver{}.Close()
	assert.Nil(t, srv)
}
