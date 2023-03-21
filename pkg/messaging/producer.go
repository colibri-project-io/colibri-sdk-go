package messaging

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/security"
	"github.com/google/uuid"
)

type Producer struct {
	topic string
}

func NewProducer(topicName string) *Producer {
	return &Producer{topicName}
}

func (p *Producer) Publish(ctx context.Context, action string, message any) {
	txn := monitoring.GetTransactionInContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			logging.Error("panic recovering publish topic %s: \n%s", p.topic, string(debug.Stack()))
			monitoring.NoticeError(txn, errors.New(string(debug.Stack())))
		}
	}()

	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, messaging_producer_transaction, map[string]any{
			"topic": p.topic,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	msg := &ProviderMessage{
		Id:      uuid.New(),
		Origin:  config.APP_NAME,
		Action:  action,
		Message: message,
	}

	authContext := security.GetAuthenticationContext(ctx)
	if authContext != nil {
		msg.TenantId = authContext.GetTenantID()
		msg.UserId = authContext.GetUserID()
	}

	if err := instance.producer(ctx, p, msg); err != nil {
		logging.Error("Could not send message with id %s to topic %s. Error: %v", msg.Id, p.topic, err)
		monitoring.NoticeError(txn, err)
	}
}
