// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/3

package rabbitmq

import (
	"sync/atomic"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"obs/logger"
)

type Option struct {
	Config  *amqp.Config
	Logger  logger.Builder
	AmqpURI string
}

type RabbitMq struct {
	Connection *amqp.Connection
	connected  uint32
	closed     uint32
	Option
}

func NewMq(opts ...func(*Option)) (*RabbitMq, error) {
	r := &RabbitMq{
		Option: Option{
			Logger:  logger.NoLogger{},
			AmqpURI: "amqp://guest:guest@localhost:5672/",
		},
	}
	for _, opt := range opts {
		opt(&r.Option)
	}

	var err error
	if r.Option.Config == nil {
		r.Connection, err = amqp.Dial(r.AmqpURI)
	} else {
		r.Connection, err = amqp.DialConfig(r.AmqpURI, *r.Option.Config)
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to AMQP")
	}
	atomic.AddUint32(&r.connected, 1)
	r.Logger.Info("connected to AMQP")
	return r, nil
}

func (r *RabbitMq) IsConnected() bool {
	return atomic.LoadUint32(&r.connected) == 1
}

func (r *RabbitMq) IsClosed() bool {
	return atomic.LoadUint32(&r.closed) == 1
}

func (r *RabbitMq) Close() error {
	r.Logger.Debug("AMQP close")
	atomic.AddUint32(&r.closed, 1)
	return r.Connection.Close()
}

// Marshaler marshals Watermill's message to amqp.Publishing and unmarshals amqp.Delivery to Watermill's message.
type Marshaler interface {
	Marshal(msg *message.Message) (amqp.Publishing, error)
	Unmarshal(amqpMsg amqp.Delivery) (*message.Message, error)
}

const DefaultMessageUUIDHeaderKey = "_micro_message_uuid"

type DefaultMarshal struct {
	// PostprocessPublishing can be used to make some extra processing with amqp.Publishing,
	// for example add CorrelationId and ContentType:
	//
	//  amqp.DefaultMarshaler{
	//		PostprocessPublishing: func(publishing stdAmqp.Publishing) stdAmqp.Publishing {
	//			publishing.CorrelationId = "correlation"
	//			publishing.ContentType = "application/json"
	//
	//			return publishing
	//		},
	//	}
	PostprocessPublishing func(amqp.Publishing) amqp.Publishing

	// When true, DeliveryMode will be not set to Persistent.
	//
	// DeliveryMode Transient means higher throughput, but messages will not be
	// restored on broker restart. The delivery mode of publishings is unrelated
	// to the durability of the queues they reside on. Transient messages will
	// not be restored to durable queues, persistent messages will be restored to
	// durable queues and lost on non-durable queues during server restart.
	NotPersistentDeliveryMode bool

	// Header used to store and read message UUID.
	//
	// If value is empty, DefaultMessageUUIDHeaderKey value is used.
	// If header doesn't exist, empty value is passed as message UUID.
	MessageUUIDHeaderKey string
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

func (d DefaultMarshal) Unmarshal(amqpMsg amqp.Delivery) (*message.Message, error) {
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
			return nil, errors.Errorf("metadata %s is not a string, but %#v", key, value)
		}
	}

	return msg, nil
}

func (d DefaultMarshal) unmarshalMessageUUID(amqpMsg amqp.Delivery) (string, error) {
	var msgUUIDStr string

	msgUUID, hasMsgUUID := amqpMsg.Headers[d.computeMessageUUIDHeaderKey()]
	if !hasMsgUUID {
		return "", nil
	}

	msgUUIDStr, hasMsgUUID = msgUUID.(string)
	if !hasMsgUUID {
		return "", errors.Errorf("message UUID is not a string, but: %#v", msgUUID)
	}

	return msgUUIDStr, nil
}

func (d DefaultMarshal) computeMessageUUIDHeaderKey() string {
	if d.MessageUUIDHeaderKey != "" {
		return d.MessageUUIDHeaderKey
	}

	return DefaultMessageUUIDHeaderKey
}
