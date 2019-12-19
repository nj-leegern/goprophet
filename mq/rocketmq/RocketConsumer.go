package rocketmq

import (
	"context"
	"github.com/nj-leegern/rocketmq-client-go"
	"github.com/nj-leegern/rocketmq-client-go/consumer"
	"github.com/nj-leegern/rocketmq-client-go/primitive"
	"strings"
)

/*
	RocketMQ消费者
*/

type RocketConsumer struct {
	// topic -> consumer
	consumers map[string]rocketmq.PushConsumer
}

/* 创建RocketMQ消费端实例 */
func NewRocketConsumer(conf RocketConsumerConf) (*RocketConsumer, error) {
	consumers := make(map[string]rocketmq.PushConsumer, 0)
	subscribers := conf.Subscribers
	for _, subscriber := range subscribers {
		// 实例化consumer
		c, err := rocketmq.NewPushConsumer(
			consumer.WithNameServer(conf.NameServers),
			consumer.WithGroupName(strings.Join([]string{conf.GroupName, subscriber.TopicName}, "_")),
			consumer.WithConsumerModel(consumer.Clustering),
			consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
			consumer.WithConsumerOrder(conf.ConsumerOrder),
		)
		if err != nil {
			return nil, err
		}
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
		consumers[subscriber.TopicName] = c
	}
	// 启动
	for _, consumer := range consumers {
		if err := consumer.Start(); err != nil {
			return nil, err
		}
	}
	return &RocketConsumer{consumers: consumers}, nil
}

/* 释放资源 */
func (c *RocketConsumer) Destroy() error {
	for _, consumer := range c.consumers {
		if err := consumer.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}
