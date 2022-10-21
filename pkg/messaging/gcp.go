package messaging

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type gcpMessaging struct {
	client *pubsub.Client
}

func newGcpMessaging() *gcpMessaging {
	client, err := pubsub.NewClient(context.Background(), "config.PROJECT_ID")
	if err != nil {
		logging.Fatal(connection_error, err)
	}

	return &gcpMessaging{client}
}

func (m *gcpMessaging) producer(ctx context.Context, p *Producer, msg *ProviderMessage) error {
	topic := m.client.Topic(p.topic)
	result := topic.Publish(ctx, &pubsub.Message{Data: []byte(msg.String())})
	_, err := result.Get(ctx)
	return err
}

func (m *gcpMessaging) consumer(ctx context.Context, queue string) (chan *ProviderMessage, error) {
	ch := make(chan *ProviderMessage, 1)
	sub := m.client.Subscription(queue)
	go func() {
		sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			var pm ProviderMessage
			if err := json.Unmarshal(([]byte(msg.Data)), &pm); err != nil {
				logging.Error(couldNotReadMsgBody, msg.ID, queue, err)
			} else {
				ch <- &pm
				msg.Ack()
			}
		})
	}()

	return ch, nil
}

func (m *gcpMessaging) sendToDLQ(ctx context.Context, queue string, msg *ProviderMessage) error {
	// Pensar melhor

	return nil
}
