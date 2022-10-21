package messaging

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

const (
	messaging_producer_transaction = "Producer"
	messaging_consumer_transaction = "Consumer-%s"
	messaging_dlq_transaction      = "Send-DLQ"
	connection_error               = "An error occurred when trying to connect to the message broker. Error: %s"
	couldNotReadMsgBody            = "Could not read message body with id %s from queue %s. Error: %v"
)

type messaging interface {
	producer(ctx context.Context, p *Producer, msg *ProviderMessage) error
	consumer(ctx context.Context, queue string) (chan *ProviderMessage, error)
	sendToDLQ(ctx context.Context, queue string, msg *ProviderMessage) error
}

var (
	instance     messaging
	appConsumers []Consumer
)

func Initialize() {
	switch config.CLOUD {
	case config.CLOUD_AWS:
		instance = newAwsMessaging()
	case config.CLOUD_GCP, config.CLOUD_FIREBASE:
		instance = newGcpMessaging()
	}

	logging.Info("Message broker connected")
	initializeConsumers()
}
