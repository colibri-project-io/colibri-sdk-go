package messaging

import (
	"context"
	"fmt"
	"testing"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

type userMessageTest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

const (
	testTopicName     = "COLIBRI_PROJECT_USER_CREATE"
	testQueueName     = "COLIBRI_PROJECT_USER_CREATE_APP_CONSUMER"
	testFailTopicName = "COLIBRI_PROJECT_FAIL_USER_CREATE"
	testFailQueueName = "COLIBRI_PROJECT_FAIL_USER_CREATE_APP_CONSUMER"
)

func TestMain(m *testing.M) {
	test.InitializeTestLocalstack()

	Initialize()

	m.Run()
}

func TestMessagingSuccess(t *testing.T) {
	chSuccess := make(chan string)
	consumer := func(ctx context.Context, message *ProviderMessage) error {
		chSuccess <- fmt.Sprintf("processing message: %v", message)
		return nil
	}

	producer := NewProducer(testTopicName)
	AddConsumer(NewConsumer(testQueueName, consumer))
	initializeConsumers()

	model := userMessageTest{"User Name", "user@email.com"}
	producer.Publish(context.Background(), "create", model)

	msgProcessing := <-chSuccess
	assert.NotEmpty(t, msgProcessing)
}

func TestMessagingFail(t *testing.T) {
	chFail := make(chan string)
	consumer := func(ctx context.Context, message *ProviderMessage) error {
		err := fmt.Errorf("email not valid")
		chFail <- err.Error()
		return err
	}

	producer := NewProducer(testFailTopicName)
	AddConsumer(NewConsumer(testFailQueueName, consumer))
	initializeConsumers()

	model := userMessageTest{"User Name", "user@email.com"}
	producer.Publish(context.Background(), "create", model)

	msgFail := <-chFail
	assert.NotEmpty(t, msgFail)
}
