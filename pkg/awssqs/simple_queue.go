// Package awssqs
// @author Daud Valentino
package awssqs

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"sharefood/pkg/logger"
)

const (
	componentName = "sqs"
)

// SimpleQueue interface
type SimpleQueue interface {
	PublishMessage(ctx context.Context, queueURL, message string) error
	DeleteMessage(ctx context.Context, queueURL string, message *MessageDecoder) error
	StartPolling(ctx context.Context, qCtx *QueueContext)
	FetchMessage(ctx context.Context, qCtx *QueueContext) (*[]MessageDecoder, error)
}

type simpleQueue struct {
	queue *sqs.SQS
}

// NewSimpleQueue return aws sqs simple queue
func NewSimpleQueue(sess *session.Session, cfg ...*aws.Config) SimpleQueue {
	return &simpleQueue{
		queue: sqs.New(sess, cfg...),
	}
}

// PublishMessage publish message to queue
func (q *simpleQueue) PublishMessage(ctx context.Context, queueURL, message string) error {

	_, err := q.queue.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(message),
	})

	return err
}

// DeleteMessage deleted message in queue
func (q *simpleQueue) DeleteMessage(ctx context.Context, queueURL string, message *MessageDecoder) error {

	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),  // Required
		ReceiptHandle: message.ReceiptHandle, // Required
	}
	_, err := q.queue.DeleteMessage(params)

	return err

}

// FetchMessage get message from queue
func (q *simpleQueue) FetchMessage(ctx context.Context, qCtx *QueueContext) (*[]MessageDecoder, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(qCtx.QueueURL), // Required
		MaxNumberOfMessages: aws.Int64(qCtx.NumberOfMessage),
		AttributeNames: []*string{
			aws.String("All"), // Required
		},
		WaitTimeSeconds: aws.Int64(maxWaitTime(qCtx.WaitTimeSecond)),
	}

	resp, err := q.queue.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}

	msgDecoders := []MessageDecoder{}

	for i := range resp.Messages {
		msgDecoders = append(msgDecoders, MessageDecoder{
			Body:          resp.Messages[i].Body,
			ReceiptHandle: resp.Messages[i].ReceiptHandle,
			MessageId:     resp.Messages[i].MessageId,
		})
	}

	return &msgDecoders, nil

}

// StartPolling start polling the message
// when message handler firing return error nil will automatically delete message
func (q *simpleQueue) StartPolling(ctx context.Context, qCtx *QueueContext) {
	for {
		select {
		case <-ctx.Done():
			logger.Warn("worker: Stopping polling because a context kill signal was sent", logger.Any("component", componentName))

			return
		default:
			logger.Debug(fmt.Sprintf("[sqs] consume message %s", queueName(qCtx.QueueURL)), logger.Any("component", componentName), logger.Any("sqs_setting", qCtx))

			params := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(qCtx.QueueURL), // Required
				MaxNumberOfMessages: aws.Int64(qCtx.NumberOfMessage),
				AttributeNames: []*string{
					aws.String("All"), // Required
				},
				WaitTimeSeconds: aws.Int64(maxWaitTime(qCtx.WaitTimeSecond)),
			}

			resp, err := q.queue.ReceiveMessage(params)
			if err != nil {
				logger.Error(logger.MessageFormat(`get message got:%v`, err),
					logger.Any("component", componentName),
					logger.Any("sqs_setting", qCtx),
				)
				continue
			}

			if len(resp.Messages) > 0 {
				q.runner(qCtx, resp.Messages)
			}
		}
	}
}

// runner messages handler
func (q *simpleQueue) runner(qCtx *QueueContext, messages []*sqs.Message) {
	numMessages := len(messages)
	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			// launch goroutine
			defer wg.Done()

			msgDecoder := &MessageDecoder{
				Body:          m.Body,
				MessageId:     m.MessageId,
				ReceiptHandle: m.ReceiptHandle,
			}

			err := qCtx.Handler(msgDecoder)

			if err == nil {

				errDel := q.DeleteMessage(context.Background(), qCtx.QueueURL, msgDecoder)
				if errDel != nil {
					logger.Error(logger.MessageFormat(`delete message got:%v`, errDel),
						logger.Any("component", componentName),
						logger.Any("sqs_setting", qCtx),
					)
				}
			}

		}(messages[i])
	}

	wg.Wait()
}

func maxWaitTime(waitTime int64) int64 {
	if waitTime < 1 {
		return 0
	}

	if waitTime > 20 {
		return 20
	}

	return waitTime
}

func queueName(queueUrl string) string {

	q := strings.Split(queueUrl, "/")

	return q[len(q)-1]
}
