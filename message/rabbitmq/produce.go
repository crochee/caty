// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/4

package rabbitmq

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// 工作模式
type Producer struct {
	*RabbitMq
	channel       *amqp.Channel
	Marshal       Marshaler
	queueName     string
	Transactional bool
}

func (p *Producer) Publish(topic string, messages ...*message.Message) error {
	if p.IsClosed() {
		return errors.New("AMQP is connection closed")
	}
	if !p.IsConnected() {
		return errors.New("not connected to AMQP")
	}
	if p.Marshal == nil {
		p.Marshal = DefaultMarshal{}
	}
	p.queueName = topic
	if p.queueName == "" {
		p.queueName = "micro"
	}
	//申请队列,如果不存在会自动创建，存在跳过创建，保证队列存在，消息能发送到队列中
	var err error
	if p.channel == nil {
		if p.channel, err = p.Connection.Channel(); err != nil {
			return errors.Wrap(err, "cannot open channel")
		}
	}

	if p.Transactional { // 开启事务
		if err = p.channel.Tx(); err != nil {
			return errors.Wrap(err, "cannot start transaction")
		}

		defer func() {
			if err != nil {
				if rollbackErr := p.channel.TxRollback(); rollbackErr != nil {
					err = multierror.Append(err, rollbackErr)
				}
				return
			}
			if commitErr := p.channel.TxCommit(); commitErr != nil {
				err = multierror.Append(err, commitErr)
			}
			return
		}()
	}
	if _, err = p.channel.QueueDeclare(
		p.queueName,
		//控制队列是否为持久的，当mq重启的时候不会丢失队列
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
		return err
	}
	if !p.Transactional {
		p.Logger.Info("enabling publishing confirms.")
		if err := p.channel.Confirm(false); err != nil {
			return errors.Wrap(err, "channel could not be put into confirm mode")
		}
	}
	for _, msg := range messages {
		if err := p.publishMessage(p.channel, msg); err != nil {
			return err
		}
	}
	if !p.Transactional {
		confirmed := <-p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		if confirmed.Ack {
			p.Logger.Infof("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
		}
		return errors.Errorf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
	return nil
}

func (p *Producer) Close() error {
	p.Logger.Debug("produce close")
	return p.channel.Close()
}

func (p *Producer) publishMessage(channel *amqp.Channel, message *message.Message) error {
	amqpMsg, err := p.Marshal.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "cannot marshal message")
	}
	//发送消息到队列中
	if err = channel.Publish(
		"",
		p.queueName,
		//如果为true，根据exchange类型和routekey类型，如果无法找到符合条件的队列，name会把发送的信息返回给发送者
		false,
		//如果为true，当exchange发送到消息队列后发现队列上没有绑定的消费者,则会将消息返还给发送者
		false,
		//发送信息
		amqpMsg,
	); err != nil {
		return errors.Wrap(err, "cannot publish msg")
	}
	return nil
}
