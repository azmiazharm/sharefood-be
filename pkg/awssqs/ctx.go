// Package awssqs
// @author Daud Valentino
package awssqs

type QueueContext struct {
	QueueURL        string `json:"queue_url"`
	NumberOfMessage int64  `json:"number_of_message"`
	WaitTimeSecond  int64  `json:"wait_time_second"`
	Handler MessageProcessorFunc `json:"-"`
}

func NewQueueContext() *QueueContext {
	return &QueueContext{}
}
