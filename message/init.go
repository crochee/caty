// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/3

package message

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"obs/logger"
	"obs/message/rabbitmq"
)

// todo self done
func Setup(ctx context.Context) error {
	//loggerx := watermill.NewStdLogger(true, true)
	//router, err := message.NewRouter(message.RouterConfig{}, loggerx)
	//if err != nil {
	//	return err
	//}
	conn, err := rabbitmq.NewMq(func(option *rabbitmq.Option) {
		option.AmqpURI = "amqp://admin:1234567@192.168.31.62:5672/"
	}, func(option *rabbitmq.Option) {
		option.Logger = logger.NewLogger()
	})
	if err != nil {
		return err
	}

	// SignalsHandler will gracefully shutdown Router when SIGTERM is received.
	// You can also close the router by just calling `r.Close()`.
	//router.AddPlugin(plugin.SignalsHandler)

	// Router level middleware are executed for every message sent to the router
	//router.AddMiddleware(
	//	// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
	//	middleware.CorrelationID,
	//
	//	// The handler function is retried if it returns an error.
	//	// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
	//	middleware.Retry{
	//		MaxRetries:      3,
	//		InitialInterval: time.Millisecond * 100,
	//		Logger:          loggerx,
	//	}.Middleware,
	//
	//	// Recoverer handles panics from handlers.
	//	// In this case, it passes them as errors to the Retry middleware.
	//	middleware.Recoverer,
	//)

	// For simplicity, we are using the gochannel Pub/Sub here,
	// You can replace it with any Pub/Sub implementation, it will work the same.
	//pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	//var publisher *amqp.Publisher
	//if publisher, err = amqp.NewPublisher(amqp.Config{
	//	Connection: amqp.ConnectionConfig{
	//		AmqpURI: "amqp://admin:1234567@192.168.31.62:5672/",
	//	},
	//	Marshaler: amqp.DefaultMarshaler{NotPersistentDeliveryMode: true},
	//	Exchange: amqp.ExchangeConfig{
	//		GenerateName: func(topic string) string {
	//			switch topic {
	//			case "micro_topic":
	//				return "micro"
	//			default:
	//				return ""
	//			}
	//		},
	//		Type: "fanout",
	//	},
	//	QueueBind: amqp.QueueBindConfig{
	//		GenerateRoutingKey: func(topic string) string {
	//			return ""
	//		},
	//	},
	//	Publish: amqp.PublishConfig{
	//		GenerateRoutingKey: func(topic string) string {
	//			return ""
	//		},
	//	},
	//	TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	//}, logger); err != nil {
	//	return err
	//}
	publisher := &rabbitmq.Producer{
		RabbitMq:      conn,
		Transactional: true,
	}
	//var subscriber *amqp.Subscriber
	//if subscriber, err = amqp.NewSubscriber(amqp.Config{
	//	Connection: amqp.ConnectionConfig{
	//		AmqpURI: "amqp://admin:1234567@192.168.31.62:5672/",
	//	},
	//	Marshaler: amqp.DefaultMarshaler{NotPersistentDeliveryMode: true},
	//	Exchange: amqp.ExchangeConfig{
	//		GenerateName: func(topic string) string {
	//			switch topic {
	//			case "micro_topic":
	//				return "micro"
	//			default:
	//				return ""
	//			}
	//		},
	//		Type: "fanout",
	//	},
	//	QueueBind: amqp.QueueBindConfig{
	//		GenerateRoutingKey: func(topic string) string {
	//			return ""
	//		},
	//	},
	//	Queue: amqp.QueueConfig{
	//		GenerateName: func(topic string) string {
	//			return topic
	//		},
	//	},
	//	Consume: amqp.ConsumeConfig{
	//		Qos: amqp.QosConfig{
	//			PrefetchCount: 1,
	//		},
	//	},
	//	TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	//}, logger); err != nil {
	//	return err
	//}
	subscriber := &rabbitmq.Consumer{
		RabbitMq:    conn,
		ConsumeName: "-0",
	}
	subscriber1 := &rabbitmq.Consumer{
		RabbitMq:    conn,
		ConsumeName: "-1",
	}
	// Producing some incoming messages in background
	go publishMessages(ctx, publisher)
	go printMessages(ctx, subscriber)
	go printMessages(ctx, subscriber1)
	go func() {
		<-ctx.Done()
		publisher.Close()
		subscriber.Close()
		subscriber1.Close()
		conn.Close()
	}()

	// AddHandler returns a handler which can be used to add handler level middleware
	//handler := router.AddHandler(
	//	"handler",        // handler name, must be unique
	//	"sub_read_topic", // topic from which we will read events
	//	subscriber,
	//	"pub_write_topic", // topic to which we will publish events
	//	publisher,
	//	HandlerFunc,
	//)

	// Handler level middleware is only executed for a specific handler
	// Such middleware can be added the same way the router level ones
	//handler.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
	//	return func(message *message.Message) ([]*message.Message, error) {
	//		log.Println("executing handler specific middleware for ", message.UUID)
	//		return h(message)
	//	}
	//})

	// just for debug, we are printing all messages received on `incoming_messages_topic`
	//router.AddNoPublisherHandler(
	//	"print_incoming_messages",
	//	"incoming_messages_topic",
	//	subscriber,
	//	printMessages,
	//)

	// just for debug, we are printing all events sent to `outgoing_messages_topic`
	//router.AddNoPublisherHandler(
	//	"print_outgoing_messages",
	//	"outgoing_messages_topic",
	//	subscriber,
	//	printMessages,
	//)

	// Now that all handlers are registered, we're running the Router.
	// Run is blocking while the router is running.
	//if err = router.Run(ctx); err != nil {
	//	return err
	//}
	return nil
}

func publishMessages(ctx context.Context, publisher message.Publisher) {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"+time.Now().String()))
			middleware.SetCorrelationID(watermill.NewUUID(), msg)

			log.Printf("sending message %s, correlation id: %s\n", msg.UUID, middleware.MessageCorrelationID(msg))

			if err := publisher.Publish("micro_topic", msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func printMessages(ctx context.Context, subscriber rabbitmq.Consume) {
	receiveChan, err := subscriber.Subscribe(ctx, "micro_topic")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case msg := <-receiveChan:
			fmt.Printf(
				"\n> Received %s message: %s\n> %s\n> metadata: %v\n\n",
				subscriber.Name(), msg.UUID, string(msg.Payload), msg.Metadata,
			)
		case <-ctx.Done():
			return
		}
	}
}

func HandlerFunc(msg *message.Message) ([]*message.Message, error) {
	log.Println("structHandler received message", msg.UUID)

	msg = message.NewMessage(watermill.NewUUID(), []byte("message produced by structHandler"))
	return message.Messages{msg}, nil
}
