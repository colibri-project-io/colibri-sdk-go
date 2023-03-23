package messaging

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrTestConsumerFail = errors.New("could not process consumer")

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

// NewTestProducer create a consumer to read queue when send messages on test environment, returning:
//   - context with timeout
//   - CancelFunc for context
//   - chan string to receive processed message in queue consumer
func NewTestProducer(queueName string, expectedMessage string, fail bool, timeout time.Duration) (context.Context, context.CancelFunc, chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	ch := make(chan string)
	qc := queueConsumerTest{
		fn: func(ctx context.Context, message *ProviderMessage) error {
			ch <- expectedMessage
			if fail {
				return ErrTestConsumerFail
			}
			return nil
		},
		qName: queueName,
	}
	NewConsumer(&qc)

	go func() {
		select {
		case <-time.After(timeout):
			panic(fmt.Sprintf("test timeout exceeded %s", timeout.String()))
		}
	}()

	return ctx, cancel, ch
}
