// Package handler
package handler

import (
	"context"

	"sharefood/internal/appctx"
	"sharefood/internal/consts"
	uContract "sharefood/internal/ucase/contract"
	"sharefood/pkg/awssqs"
)

// SQSConsumerHandler sqs consumer message processor handler
func SQSConsumerHandler(msgHandler uContract.MessageProcessor) awssqs.MessageProcessorFunc {
	return func(decoder *awssqs.MessageDecoder) error {
		return msgHandler.Serve(context.Background(), &appctx.ConsumerData{
			Body:        []byte(*decoder.Body),
			Key:         []byte(*decoder.MessageId),
			ServiceType: consts.ServiceTypeConsumer,
		})
	}
}
