package mq

import (
	"errors"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/crochee/lib/e"
	"github.com/hashicorp/go-multierror"
	"github.com/streadway/amqp"
)

type producer struct {
	transactional bool
	mq            *rabbitMq
	exchange      string
	routingKey    string
	wg            sync.WaitGroup
	marshal       MarshalAPI
	queueName     func(string) string
}

// NewProducer create message.Publisher
func NewProducer(mq *rabbitMq, opts ...func(*Option)) message.Publisher {
	option := Option{
		Marshal: DefaultMarshal{},
		QueueName: func(topic string) string {
			return topic
		},
	}
	for _, opt := range opts {
		opt(&option)
	}
	return &producer{
		transactional: option.Transactional,
		mq:            mq,
		queueName:     option.QueueName,
		exchange:      option.Exchange,
		routingKey:    option.RoutingKey,
		marshal:       option.Marshal,
	}
}

func (p *producer) Publish(topic string, messages ...*message.Message) error {
	if p.mq.IsClosed() {
		return errors.New("AMQP is connection closed")
	}
	if !p.mq.IsConnected() {
		return errors.New("not connected to AMQP")
	}
	p.wg.Add(1)
	defer p.wg.Done()
	// 申请队列,如果不存在会自动创建，存在跳过创建，保证队列存在，消息能发送到队列中
	channel, err := p.mq.Channel()
	if err != nil {
		return e.Wrap(err, "cannot open channel")
	}
	defer func() {
		if channelCloseErr := channel.Close(); channelCloseErr != nil {
			err = multierror.Append(err, channelCloseErr)
		}
	}()
	if p.transactional { // 开启事务
		if err = channel.Tx(); err != nil {
			return e.Wrap(err, "cannot start transaction")
		}
		defer func() {
			if err != nil {
				if rollbackErr := channel.TxRollback(); rollbackErr != nil {
					err = multierror.Append(err, rollbackErr)
				}
				return
			}
			if commitErr := channel.TxCommit(); commitErr != nil {
				err = multierror.Append(err, commitErr)
			}
		}()
	}
	if _, err = channel.QueueDeclare(
		p.queueName(topic),
		// 控制队列是否为持久的，当mq重启的时候不会丢失队列
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
		return err
	}
	if !p.transactional {
		if err = channel.Confirm(false); err != nil {
			return e.Wrap(err, "channel could not be put into confirm mode")
		}
	}
	for _, msg := range messages {
		if err = p.publishMessage(channel, msg); err != nil {
			return err
		}
	}
	if !p.transactional {
		confirmed := <-channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		if !confirmed.Ack {
			return e.Errorf("failed delivery of delivery tag: %d ack:%v", confirmed.DeliveryTag, confirmed.Ack)
		}
	}
	return nil
}

func (p *producer) Close() error {
	p.wg.Wait()
	return nil
}

func (p *producer) publishMessage(channel *amqp.Channel, msg *message.Message) error {
	amqpMsg, err := p.marshal.Marshal(msg)
	if err != nil {
		return e.Wrap(err, "cannot marshal message")
	}
	// 发送消息到队列中
	if err = channel.Publish(
		p.exchange,
		p.routingKey,
		// 如果为true，根据exchange类型和routekey类型，如果无法找到符合条件的队列，name会把发送的信息返回给发送者
		false,
		// 如果为true，当exchange发送到消息队列后发现队列上没有绑定的消费者,则会将消息返还给发送者
		false,
		// 发送信息
		amqpMsg,
	); err != nil {
		return e.Wrap(err, "cannot publish msg")
	}
	return nil
}
