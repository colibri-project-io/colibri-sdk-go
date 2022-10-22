package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/cloud"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type sqsNotification struct {
	MessageId string `json:"MessageId"`
	Message   string `json:"Message"`
}

type awsMessaging struct {
	snsService *sns.SNS
	sqsService *sqs.SQS
}

func newAwsMessaging() *awsMessaging {
	var m awsMessaging

	m.snsService = sns.New(cloud.GetAwsSession())
	m.sqsService = sqs.New(cloud.GetAwsSession())

	if _, err := m.snsService.ListTopics(nil); err != nil {
		logging.Fatal(connection_error, err)
	}

	return &m
}

func (m *awsMessaging) producer(ctx context.Context, p *Producer, msg *ProviderMessage) error {
	_, err := m.snsService.PublishWithContext(ctx, &sns.PublishInput{
		Message:  aws.String(msg.String()),
		TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:us-east-1:000000000000:%s", p.topic)),
	})

	return err
}

func (m *awsMessaging) consumer(ctx context.Context, c *Consumer) (chan *ProviderMessage, error) {
	ch := make(chan *ProviderMessage, 1)
	queueUrl := m.getQueueUrl(ctx, c.queue)

	go func() {
		for {
			if c.isCanceled() {
				c.Done()
				return
			}
			msgs, err := m.readMessages(ctx, queueUrl)
			if err != nil {
				logging.Error("Could not read messages from queue %s. Error: %v", c.queue, err)
			}

			if len(msgs.Messages) > 0 {
				msg := msgs.Messages[0]

				var n sqsNotification
				if err = json.Unmarshal([]byte(*msg.Body), &n); err != nil {
					logging.Error(couldNotReadMsgBody, *msg.MessageId, c.queue, err)
				}

				var pm ProviderMessage
				if err = json.Unmarshal([]byte(n.Message), &pm); err != nil {
					logging.Error(couldNotReadMsgBody, *msg.MessageId, c.queue, err)
				} else {
					ch <- &pm
					m.removeMessageFromQueue(ctx, queueUrl, msg)
				}
			}
		}
	}()

	return ch, nil
}

func (m *awsMessaging) readMessages(ctx context.Context, queueResult *sqs.GetQueueUrlOutput) (*sqs.ReceiveMessageOutput, error) {
	var msgs, err = m.sqsService.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              queueResult.QueueUrl,
		MaxNumberOfMessages:   aws.Int64(1),
		WaitTimeSeconds:       aws.Int64(1),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	return msgs, err
}

func (m *awsMessaging) removeMessageFromQueue(ctx context.Context, queueResult *sqs.GetQueueUrlOutput, msg *sqs.Message) {
	if _, err := m.sqsService.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      queueResult.QueueUrl,
		ReceiptHandle: msg.ReceiptHandle,
	}); err != nil {
		logging.Error("Could not delete message with id %s from queue %s. Error: %v", *msg.MessageId, *queueResult.QueueUrl, err)
	}
}

func (m *awsMessaging) sendToDLQ(ctx context.Context, queue string, msg *ProviderMessage) error {
	queueUrl := m.getQueueUrl(ctx, queue)

	if _, err := m.sqsService.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(msg.String()),
		QueueUrl:    queueUrl.QueueUrl,
	}); err != nil {
		return err
	}

	return nil
}

func (m *awsMessaging) getQueueUrl(ctx context.Context, queue string) *sqs.GetQueueUrlOutput {
	queueResult, err := m.sqsService.GetQueueUrlWithContext(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(queue)})
	if err != nil {
		logging.Fatal("Could not connect to queue %s. Error: %v", queue, err)
	}

	if queueResult.QueueUrl == nil {
		logging.Fatal("Queue %s not found", queue)
	}

	return queueResult
}
