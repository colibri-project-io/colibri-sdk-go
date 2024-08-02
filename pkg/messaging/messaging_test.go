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

type queueConsumerTest struct {
	fn    func(ctx context.Context, n *ProviderMessage) error
	qName string
}

func (q *queueConsumerTest) Consume(ctx context.Context, pm *ProviderMessage) error {
	return q.fn(ctx, pm)
}

func (q *queueConsumerTest) QueueName() string {
	return q.qName
}

func TestMessaging_AWS(t *testing.T) {
	test.InitializeTestLocalstack()

	Initialize()

	executeMessagingTest(t)
}

func TestMessaging_GCP(t *testing.T) {
	test.InitializeGcpEmulator()

	Initialize()

	executeMessagingTest(t)
}

func executeMessagingTest(t *testing.T) {
	t.Run("Should return nil when process message with success", func(t *testing.T) {
		chSuccess := make(chan string)
		qc := queueConsumerTest{
			fn: func(ctx context.Context, message *ProviderMessage) error {
				chSuccess <- fmt.Sprintf("processing message: %v", message)
				return nil
			},
			qName: testQueueName,
		}

		producer := NewProducer(testTopicName)
		NewConsumer(&qc)

		model := userMessageTest{"User Name", "user@email.com"}
		producer.Publish(context.Background(), "create", model)

		timeout := time.After(2 * time.Second)
		select {
		case msgProcessing := <-chSuccess:
			assert.NotEmpty(t, msgProcessing)
		case <-timeout:
			t.Fatal("Test didn't finish after 2s")
		}
	})

	t.Run("Should return error when process message with error and send message to dlq", func(t *testing.T) {
		chFail := make(chan string)
		qc := queueConsumerTest{
			fn: func(ctx context.Context, message *ProviderMessage) error {
				err := fmt.Errorf("email not valid")
				chFail <- err.Error()
				return err
			},
			qName: testFailQueueName,
		}

		producer := NewProducer(testFailTopicName)
		NewConsumer(&qc)

		model := userMessageTest{"User Name", "user@email.com"}
		producer.Publish(context.Background(), "create", model)

		timeout := time.After(2 * time.Second)
		select {
		case msgDLQ := <-chFail:
			assert.Equal(t, "email not valid", msgDLQ)
		case <-timeout:
			t.Fatal("Test didn't finish after 2s")
		}
	})
}
