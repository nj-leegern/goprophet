package rocketmq

import (
	"context"
	"github.com/nj-leegern/rocketmq-client-go"
	"github.com/nj-leegern/rocketmq-client-go/consumer"
	"github.com/nj-leegern/rocketmq-client-go/primitive"
)

/*
	RocketMQ消费者
*/

type RocketConsumer struct {
	c rocketmq.PushConsumer
}

/* 创建RocketMQ消费端实例 */
func NewRocketConsumer(conf RocketConsumerConf) (*RocketConsumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(conf.NameServers),
		consumer.WithGroupName(conf.GroupName),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerOrder(conf.ConsumerOrder),
	)
	if err != nil {
		return nil, err
	}
	subscribers := conf.Subscribers
	for _, subscriber := range subscribers {
		selector := consumer.MessageSelector{}
		// 订阅标签
		if len(subscriber.Tag) > 0 {
			selector.Type = consumer.TAG
			selector.Expression = subscriber.Tag
		}
		// 订阅消息
		err = c.Subscribe(subscriber.TopicName, selector, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			messages := make([]string, 0)
			for _, msg := range msgs {
				messages = append(messages, string(msg.Body))
			}
			// 消息处理回调
			err = subscriber.HandleMsg(messages)
			if err != nil {
				return consumer.ConsumeRetryLater, err
			}
			return consumer.ConsumeSuccess, nil
		})
		if err != nil {
			return nil, err
		}
	}
	err = c.Start()
	if err != nil {
		return nil, err
	}
	return &RocketConsumer{c: c}, nil
}

/* 释放资源 */
func (c *RocketConsumer) Destroy() error {
	return c.c.Shutdown()
}
