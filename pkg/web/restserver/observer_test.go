package restserver

import (
	"github.com/gofiber/fiber/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseServer(t *testing.T) {
	srv = &FiberWebServer{srv: &fiber.App{}}

	restObserver{}.Close()
	assert.Nil(t, srv)
}
