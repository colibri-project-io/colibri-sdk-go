package messaging

import (
	"context"
	"fmt"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/graceful-shutdown"
	"sync"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
)

type Consumer struct {
	sync.WaitGroup
	queue  string
	fn     func(ctx context.Context, message *ProviderMessage) error
	done   chan interface{}
	hasDLQ bool
}

type consumerObserver struct {
	c *Consumer
}

func (o consumerObserver) Close() {
	o.c.close()
}

func NewConsumerWithDLQ(queueName string, function func(ctx context.Context, message *ProviderMessage) error) {
	newConsumer(queueName, true, function)
}

func NewConsumerWithoutDLQ(queueName string, function func(ctx context.Context, message *ProviderMessage) error) {
	newConsumer(queueName, false, function)
}

func newConsumer(queueName string, hasDLQ bool, function func(ctx context.Context, message *ProviderMessage) error) {
	c := &Consumer{
		WaitGroup: sync.WaitGroup{},
		queue:     queueName,
		fn:        function,
		done:      make(chan interface{}),
		hasDLQ:    hasDLQ,
	}

	gracefulshutdown.Attach(consumerObserver{c: c})
	startListener(c)
}

func startListener(c *Consumer) {
	ch := createConsumer(c)

	go func() {
		for {
			msg := <-ch
			ctx := context.Background()
			security.NewAuthenticationContext(msg.TenantId, msg.UserId).SetInContext(ctx)
			if err := c.fn(ctx, msg); err != nil {
				c.sendMessageDLQ(ctx, msg, err)
			}
		}
	}()
}

func createConsumer(c *Consumer) chan *ProviderMessage {
	txn, ctx := monitoring.StartTransaction(fmt.Sprintf(messaging_consumer_transaction, c.queue))
	defer monitoring.EndTransaction(txn)

	ch, err := instance.consumer(ctx, c)
	if err != nil {
		logging.Error("An error occurred when trying to create a consumer to queue %s. Error: %v", c.queue, err)
		monitoring.NoticeError(txn, err)
		return nil
	}

	return ch
}

func (c *Consumer) sendMessageDLQ(ctx context.Context, msg *ProviderMessage, err error) {
	if !c.hasDLQ {
		return
	}
	dlqQueueName := fmt.Sprintf("%s_DLQ", c.queue)
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, messaging_dlq_transaction, map[string]interface{}{
			"queue": c.queue,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	logging.Error("An error occurred on processing message with id %s from queue %s. Sending to DLQ. Error: %v", msg.Id, c.queue, err)
	monitoring.NoticeError(txn, err)

	if err = instance.sendToDLQ(ctx, dlqQueueName, msg); err != nil {
		logging.Error("Error on send message with id %s from DQL %s. Error: %v", msg.Id, dlqQueueName, err)
		monitoring.NoticeError(txn, err)
	}
}

func (c *Consumer) close() {
	logging.Info("closing queue consumer %s", c.queue)
	close(c.done)
	c.Wait()
}

func (c *Consumer) isCanceled() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}
