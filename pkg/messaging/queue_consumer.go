package messaging

import "context"

type QueueConsumer interface {
	Consume(ctx context.Context, providerMessage *ProviderMessage) error
	QueueName() string
}
