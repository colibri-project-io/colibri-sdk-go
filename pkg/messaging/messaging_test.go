package messaging

import (
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type userMessageTest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

const (
	testTopicName        = "COLIBRI_PROJECT_USER_CREATE"
	testQueueName        = "COLIBRI_PROJECT_USER_CREATE_APP_CONSUMER"
	testFailTopicName    = "COLIBRI_PROJECT_FAIL_USER_CREATE"
	testFailQueueName    = "COLIBRI_PROJECT_FAIL_USER_CREATE_APP_CONSUMER"
	testFailDLQQueueName = "COLIBRI_PROJECT_FAIL_USER_CREATE_APP_CONSUMER_DLQ"
)

func TestMain(m *testing.M) {
	test.InitializeBaseTest()
	test.InitializeTestLocalstack()

	Initialize()

	m.Run()
}

func TestMessagingSuccess(t *testing.T) {
	producer := NewProducer(testTopicName)
	expectedMessage := "processed consumer"
	ctx, cancel, ch := NewTestProducer(testQueueName, expectedMessage, false, 2*time.Second)
	defer cancel()

	model := userMessageTest{"User Name", "user@email.com"}
	assert.NoError(t, producer.Publish(ctx, "create", model))

	msg := <-ch
	assert.Equal(t, expectedMessage, msg)
}

func TestMessagingFail(t *testing.T) {
	expectedMessage := "email not valid"
	ctx, cancel, ch := NewTestProducer(testFailQueueName, expectedMessage, true, 2*time.Second)
	defer cancel()

	model := userMessageTest{"User Name", "user@email.com"}

	producer := NewProducer(testFailTopicName)
	assert.NoError(t, producer.Publish(ctx, "create", model))

	msg := <-ch
	assert.Equal(t, expectedMessage, msg)
}
