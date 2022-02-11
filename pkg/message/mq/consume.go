package mq

import (
	"context"
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/crochee/lirity/e"
	"github.com/crochee/lirity/logger"
)

type consumer struct {
	mq        *Client
	wg        sync.WaitGroup
	marshal   MarshalAPI
	queueName func(string) string
}

// NewConsumer create message.Subscriber
func NewConsumer(mq *Client, opts ...func(*Option)) message.Subscriber {
	option := Option{
		Marshal: DefaultMarshal{},
		QueueName: func(topic string) string {
			return topic
		},
	}
	for _, opt := range opts {
		opt(&option)
	}
	return &consumer{
		mq:        mq,
		marshal:   option.Marshal,
		queueName: option.QueueName,
	}
}

func (c *consumer) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	if c.mq.IsClosed() {
		return nil, errors.New("AMQP is connection closed")
	}
	if !c.mq.IsConnected() {
		return nil, errors.New("not connected to AMQP")
	}
	channel, err := c.mq.Channel()
	if err != nil {
		return nil, e.Wrap(err, "cannot open channel")
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	if err = channel.Qos(1, 0, true); err != nil { // nolint:gocritic
		return nil, err
	}
	queueName := c.queueName(topic)
	if _, err = channel.QueueDeclare(
		queueName,
		// 控制消息是否持久化，true开启
		true,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	); err != nil {
		return nil, err
	}
	c.wg.Add(1)
	output := make(chan *message.Message)
	go func() {
		for {
			deliveries, err := channel.Consume(
				queueName,
				// 用来区分多个消费者
				"dcs",
				// 是否自动应答(自动应答确认消息，这里设置为否，在下面手动应答确认)
				false,
				// 是否具有排他性
				false,
				// 如果设置为true，表示不能将同一个connection中发送的消息
				// 传递给同一个connection的消费者
				false,
				// 是否为阻塞
				false,
				nil,
			)
			if err != nil {
				logger.From(ctx).Error(err.Error())
				goto label
			}
			for {
				select {
				case d := <-deliveries:
					msgStruct, err := c.marshal.Unmarshal(&d) // nolint:govet
					if err != nil {
						logger.From(ctx).Error(err.Error())
						// 当requeue为true时，将该消息排队，以在另一个通道上传递给使用者。
						// 当requeue为false或服务器无法将该消息排队时，它将被丢弃。
						if err = d.Reject(false); err != nil {
							logger.From(ctx).Error(err.Error())
							goto label
						}
						continue
					}
					// 手动确认收到本条消息, true表示回复当前信道所有未回复的ack，用于批量确认。
					// false表示回复当前条目
					if err = d.Ack(false); err != nil {
						logger.From(ctx).Error(err.Error())
						goto label
					}
					output <- msgStruct
				case <-ctx.Done():
					close(output)
					c.wg.Done()
					if err = channel.Close(); err != nil {
						logger.From(ctx).Error(err.Error())
					}
					return
				}
			}
		label:
		}
	}()
	return output, nil
}

func (c *consumer) Close() error {
	c.wg.Wait()
	return nil
}
