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

func (m *gcpMessaging) consumer(ctx context.Context, c *consumer) (chan *ProviderMessage, error) {
	ch := make(chan *ProviderMessage, 1)
	sub := m.client.Subscription(c.queue)
	go func() {
		err := sub.Receive(ctx, func(innerCtx context.Context, msg *pubsub.Message) {
			if c.isCanceled() {
				c.Done()
				return
			}
			var pm ProviderMessage
			if err := json.Unmarshal(msg.Data, &pm); err != nil {
				logging.Error(couldNotReadMsgBody, msg.ID, c.queue, err)
			} else {
				ch <- &pm
				msg.Ack()
			}
		})
		if err != nil {
			logging.Error("Error on receive message from queue %s: %v", c.queue, err)
		}
	}()

	return ch, nil
}
