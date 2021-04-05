// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/4

package rabbitmq

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Consume interface {
	message.Subscriber
	Name() string
}

type Consumer struct {
	*RabbitMq
	channel     *amqp.Channel
	Marshal     Marshaler
	queueName   string
	ConsumeName string
}

func (c *Consumer) Name() string {
	return c.ConsumeName
}

func (c *Consumer) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	if c.IsClosed() {
		return nil, errors.New("AMQP is connection closed")
	}

	if !c.IsConnected() {
		return nil, errors.New("not connected to AMQP")
	}
	if c.Marshal == nil {
		c.Marshal = DefaultMarshal{}
	}
	c.queueName = topic
	if c.queueName == "" {
		c.queueName = "micro"
	}
	var err error
	if c.channel, err = c.Connection.Channel(); err != nil {
		return nil, errors.Wrap(err, "cannot open channel")
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	if err = c.channel.Qos(1, 0, true); err != nil {
		return nil, err
	}

	if _, err = c.channel.QueueDeclare(
		c.queueName,
		//控制消息是否持久化，true开启
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	); err != nil {
		return nil, err
	}
	output := make(chan *message.Message)
	go func() {
		for {
			deliveries, err := c.channel.Consume(
				c.queueName,
				//用来区分多个消费者
				"OBS",
				//是否自动应答(自动应答确认消息，这里设置为否，在下面手动应答确认)
				false,
				//是否具有排他性
				false,
				//如果设置为true，表示不能将同一个connection中发送的消息
				//传递给同一个connection的消费者
				false,
				//是否为阻塞
				false,
				nil,
			)
			if err != nil {
				c.Logger.Error(err.Error())
				goto label
			}
			for {
				select {
				case d := <-deliveries:
					msgStruct, err := c.Marshal.Unmarshal(d)
					if err != nil {
						c.Logger.Error(err.Error())
						// 当requeue为true时，将该消息排队，以在另一个通道上传递给使用者。
						//当requeue为false或服务器无法将该消息排队时，它将被丢弃。
						if err = d.Reject(false); err != nil {
							c.Logger.Error(err.Error())
							goto label
						}
						continue
					}
					// 手动确认收到本条消息, true表示回复当前信道所有未回复的ack，用于批量确认。
					//false表示回复当前条目
					if err = d.Ack(false); err != nil {
						c.Logger.Error(err.Error())
						goto label
					}
					output <- msgStruct
				case <-ctx.Done():
					close(output)
					return
				}
			}
		label:
		}
	}()
	return output, nil
}

func (c *Consumer) Close() error {
	c.Logger.Debug("consume close")
	return c.channel.Close()
}
