package messaging

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/test"
	"github.com/stretchr/testify/assert"
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
	chSuccess := make(chan string)
	consumer := func(ctx context.Context, message *ProviderMessage) error {
		chSuccess <- fmt.Sprintf("processing message: %v", message)
		return nil
	}

	producer := NewProducer(testTopicName)
	NewConsumerWithDLQ(testQueueName, consumer)

	model := userMessageTest{"User Name", "user@email.com"}
	producer.Publish(context.Background(), "create", model)

	timeout := time.After(2 * time.Second)
	select {
	case msgProcessing := <-chSuccess:
		assert.NotEmpty(t, msgProcessing)
	case <-timeout:
		t.Fatal("Test didn't finish after 2s")
	}
}

func TestMessagingFail(t *testing.T) {
	chFail := make(chan string)
	chDLQ := make(chan string)
	consumedDLQMsg := "consumed dlq"
	consumer := func(ctx context.Context, message *ProviderMessage) error {
		err := fmt.Errorf("email not valid")
		chFail <- err.Error()
		return err
	}
	consumerDLQ := func(ctx context.Context, message *ProviderMessage) error {
		fmt.Println(message)
		chDLQ <- "consumed dlq"
		return nil
	}

	producer := NewProducer(testFailTopicName)
	NewConsumerWithDLQ(testFailQueueName, consumer)
	NewConsumerWithoutDLQ(testFailDLQQueueName, consumerDLQ)

	model := userMessageTest{"User Name", "user@email.com"}
	producer.Publish(context.Background(), "create", model)

	msgFail := <-chFail
	assert.NotEmpty(t, msgFail)

	timeout := time.After(2 * time.Second)
	select {
	case msgDLQ := <-chDLQ:
		assert.Equal(t, consumedDLQMsg, msgDLQ)
	case <-timeout:
		t.Fatal("Test didn't finish after 2s")
	}
}
