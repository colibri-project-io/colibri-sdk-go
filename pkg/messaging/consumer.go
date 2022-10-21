package messaging

import (
	"context"
	"fmt"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
)

type Consumer struct {
	queue string
	fn    func(ctx context.Context, message *ProviderMessage) error
}

func NewConsumer(queueName string, function func(ctx context.Context, message *ProviderMessage) error) (c Consumer) {
	c.queue = queueName
	c.fn = function
	return
}

func AddConsumer(consumer Consumer) {
	appConsumers = append(appConsumers, consumer)
}

func initializeConsumers() {
	for _, consumer := range appConsumers {
		consumer.executeConsumer()
	}
}

func (c Consumer) executeConsumer() {
	ctx, ch := createConsumer(c.queue)

	go func() {
		for {
			msg := <-ch
			security.NewAuthenticationContext(msg.TenantId, msg.UserId).SetInContext(ctx)
			if err := c.fn(ctx, msg); err != nil {
				c.sendMessageDLQ(ctx, msg, err)
			}
		}
	}()
}

func createConsumer(queueName string) (context.Context, chan *ProviderMessage) {
	txn, ctx := monitoring.StartTransaction(fmt.Sprintf(messaging_consumer_transaction, queueName))
	defer monitoring.EndTransaction(txn)

	ch, err := instance.consumer(ctx, queueName)
	if err != nil {
		logging.Error("An error occurred when trying to create a consumer to queue %s. Error: %v", queueName, err)
		monitoring.NoticeError(txn, err)
		return nil, nil
	}

	return ctx, ch
}

func (c Consumer) sendMessageDLQ(ctx context.Context, msg *ProviderMessage, err error) {
	dlqQueueName := fmt.Sprintf("%s_DLQ", c.queue)
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, messaging_dlq_transaction, map[string]interface{}{
			"queue": c.queue,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	logging.Error("An error ocurred on processing message with id %s from queue %s. Sending to DLQ. Error: %v", msg.Id, c.queue, err)
	monitoring.NoticeError(txn, err)

	if err = instance.sendToDLQ(ctx, dlqQueueName, msg); err != nil {
		logging.Error("Error on send message with id %s from DQL %s. Error: %v", msg.Id, dlqQueueName, err)
		monitoring.NoticeError(txn, err)
	}
}
