package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	Initialize()
	assert.NotNil(t, instance)
	assert.NotNil(t, instance.validator)
	assert.NotNil(t, instance.formDecoder)
}
