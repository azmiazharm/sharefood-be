// Package awssqs
// @author Daud Valentino
package awssqs

import (
	"encoding/json"
	"fmt"
	"reflect"
)


type MessageProcessorFunc func(*MessageDecoder) error

type MessageDecoder struct {
	// The message's contents (not URL-encoded).
	Body *string `type:"string"`

	// A unique identifier for the message. A MessageIdis considered unique across
	// all AWS accounts for an extended period of time.
	MessageId *string `type:"string"`

	// An identifier associated with the act of receiving the message. A new receipt
	// handle is returned every time you receive a message. When deleting a message,
	// you provide the last received receipt handle to delete the message.
	ReceiptHandle *string `type:"string"`
}

// DecodeJSON  json string message body decode
func (m *MessageDecoder) DecodeJSON(out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("%s", "output destination cannot addressable")
	}

	err := json.Unmarshal([]byte(*m.Body), out)

	return err
}
