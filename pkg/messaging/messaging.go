package messaging

import (
	"context"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
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
	consumer(ctx context.Context, c *consumer) (chan *ProviderMessage, error)
}

var instance messaging

type messagingObserver struct {
	closed bool
}

func (o *messagingObserver) Close() {
	logging.Info("waiting to safely close messaging module")
	if observer.WaitRunningTimeout() {
		logging.Warn("WaitGroup timed out, forcing close the messaging module")
	}

	o.closed = true
}

func Initialize() {
	switch config.CLOUD {
	case config.CLOUD_AWS:
		instance = newAwsMessaging()
	case config.CLOUD_GCP, config.CLOUD_FIREBASE:
		instance = newGcpMessaging()
	}

	logging.Info("Message broker connected")
	observer.Attach(&messagingObserver{})
}
