package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
)

type consumer struct {
	sync.WaitGroup
	queue string
	fn    func(ctx context.Context, message *ProviderMessage) error
	done  chan interface{}
}

type consumerObserver struct {
	c *consumer
}

func (o consumerObserver) Close() {
	o.c.close()
}

func NewConsumer(qc QueueConsumer) {
	c := &consumer{
		WaitGroup: sync.WaitGroup{},
		queue:     qc.QueueName(),
		fn:        qc.Consume,
		done:      make(chan interface{}),
	}

	observer.Attach(consumerObserver{c: c})
	startListener(c)
}

func startListener(c *consumer) {
	ch := createConsumer(c)

	go func() {
		for {
			msg := <-ch
			ctx := context.Background()
			security.NewAuthenticationContext(msg.TenantId, msg.UserId).SetInContext(ctx)
			if err := c.fn(ctx, msg); err != nil {
				logging.Error("could not process message %s: %v", msg.Id, err)
			}
		}
	}()
}

func createConsumer(c *consumer) chan *ProviderMessage {
	txn, ctx := monitoring.StartTransaction(context.Background(), fmt.Sprintf(messaging_consumer_transaction, c.queue))
	defer monitoring.EndTransaction(txn)

	ch, err := instance.consumer(ctx, c)
	if err != nil {
		logging.Error("An error occurred when trying to create a consumer to queue %s: %v", c.queue, err)
		monitoring.NoticeError(txn, err)
		return nil
	}

	return ch
}

func (c *consumer) close() {
	logging.Info("Closing queue consumer %s", c.queue)
	close(c.done)
	c.Wait()
}

func (c *consumer) isCanceled() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}
