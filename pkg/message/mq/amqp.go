// Author: crochee
// Date: 2021/9/3

// Package mq
package mq

import (
	"fmt"
	"sync/atomic"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/streadway/amqp"
)

type Option struct {
	Transactional bool
	Config        *amqp.Config
	URI           string
	Marshal       MarshalAPI
	QueueName     func(topic string) string
	Exchange      string
	RoutingKey    string
}

type rabbitMq struct {
	*amqp.Connection
	connected uint32
}

// NewMq create a mq client
func NewMq(opts ...func(*Option)) (*rabbitMq, error) {
	option := Option{
		URI: "amqp://guest:guest@localhost:5672/",
	}
	for _, opt := range opts {
		opt(&option)
	}
	r := &rabbitMq{}
	var err error
	if option.Config == nil {
		r.Connection, err = amqp.Dial(option.URI)
	} else {
		r.Connection, err = amqp.DialConfig(option.URI, *option.Config)
	}
	if err != nil {
		return nil, err
	}
	atomic.AddUint32(&r.connected, 1)
	return r, nil
}

func (r *rabbitMq) IsConnected() bool {
	return atomic.LoadUint32(&r.connected) == 1
}

type MarshalAPI interface {
	Marshal(msg *message.Message) (amqp.Publishing, error)
	Unmarshal(amqpMsg *amqp.Delivery) (*message.Message, error)
}

const DefaultMessageUUIDHeaderKey = "_message_uuid"

type DefaultMarshal struct {
	PostprocessPublishing     func(amqp.Publishing) amqp.Publishing
	NotPersistentDeliveryMode bool
	MessageUUIDHeaderKey      string
}

func (d DefaultMarshal) Marshal(msg *message.Message) (amqp.Publishing, error) {
	headers := make(amqp.Table, len(msg.Metadata)+1) // metadata + plus uuid

	for key, value := range msg.Metadata {
		headers[key] = value
	}
	headers[d.computeMessageUUIDHeaderKey()] = msg.UUID

	publishing := amqp.Publishing{
		Body:    msg.Payload,
		Headers: headers,
	}
	if !d.NotPersistentDeliveryMode {
		publishing.DeliveryMode = amqp.Persistent
	}

	if d.PostprocessPublishing != nil {
		publishing = d.PostprocessPublishing(publishing)
	}

	return publishing, nil
}

func (d DefaultMarshal) Unmarshal(amqpMsg *amqp.Delivery) (*message.Message, error) {
	msgUUIDStr, err := d.unmarshalMessageUUID(amqpMsg)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(msgUUIDStr, amqpMsg.Body)
	msg.Metadata = make(message.Metadata, len(amqpMsg.Headers)-1) // headers - minus uuid

	for key, value := range amqpMsg.Headers {
		if key == d.computeMessageUUIDHeaderKey() {
			continue
		}

		var ok bool
		msg.Metadata[key], ok = value.(string)
		if !ok {
			return nil, fmt.Errorf("metadata %s is not a string, but %#v", key, value)
		}
	}
	return msg, nil
}

func (d DefaultMarshal) unmarshalMessageUUID(amqpMsg *amqp.Delivery) (string, error) {
	msgUUID, hasMsgUUID := amqpMsg.Headers[d.computeMessageUUIDHeaderKey()]
	if !hasMsgUUID {
		return "", nil
	}
	var msgUUIDStr string
	if msgUUIDStr, hasMsgUUID = msgUUID.(string); !hasMsgUUID {
		return "", fmt.Errorf("message UUID is not a string, but: %#v", msgUUID)
	}
	return msgUUIDStr, nil
}

func (d DefaultMarshal) computeMessageUUIDHeaderKey() string {
	if d.MessageUUIDHeaderKey != "" {
		return d.MessageUUIDHeaderKey
	}
	return DefaultMessageUUIDHeaderKey
}
