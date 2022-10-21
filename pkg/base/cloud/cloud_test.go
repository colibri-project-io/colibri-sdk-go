package cloud

import (
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	t.Run("Should nil if not initialize", func(t *testing.T) {
		assert.Nil(t, instance)
	})

	t.Run("Should initialize AWS with local enviroment", func(t *testing.T) {
		config.ENVIRONMENT = config.ENV_DEVELOPMENT
		config.CLOUD = config.CLOUD_AWS

		Initialize()

		assert.NotNil(t, instance)
		assert.NotNil(t, instance.aws)
		assert.NotNil(t, GetAwsSession())
	})

	t.Run("Should initialize AWS with cloud enviroment", func(t *testing.T) {
		config.ENVIRONMENT = config.ENV_PRODUCTION
		config.CLOUD = config.CLOUD_AWS

		Initialize()

		assert.NotNil(t, instance)
		assert.NotNil(t, instance.aws)
		assert.NotNil(t, GetAwsSession())
	})
}
