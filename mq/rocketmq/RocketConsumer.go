package rocketmq

import (
	"context"
	"github.com/swift9/rocketmq-client-go/v2"
	"github.com/swift9/rocketmq-client-go/v2/consumer"
	"github.com/swift9/rocketmq-client-go/v2/primitive"
	"strings"
)

/*
	RocketMQ消费者
*/

type RocketConsumer struct {
	// topic -> consumer
	consumers map[string]rocketmq.PushConsumer
	// topic -> handleMsg
	handlers map[string]Subscription
}

/* 创建RocketMQ消费端实例 */
func NewRocketConsumer(conf RocketConsumerConf) (*RocketConsumer, error) {
	client := RocketConsumer{}
	consumers := make(map[string]rocketmq.PushConsumer, 0)
	handlers := make(map[string]Subscription, 0)
	subscribers := conf.Subscribers

	for _, subscriber := range subscribers {
		handlers[subscriber.TopicName] = subscriber
	}

	for _, subscriber := range subscribers {
		groupName := strings.Join([]string{conf.GroupName, subscriber.TopicName}, "_")
		// 实例化consumer
		c, err := rocketmq.NewPushConsumer(
			consumer.WithNameServer(conf.NameServers),
			consumer.WithGroupName(groupName),
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
			for _, msg := range msgs {
				// 消息回调处理
				err = handlers[msg.Topic].HandleMsg([]string{string(msg.Body)})
				if err != nil {
					return consumer.ConsumeRetryLater, err
				}
			}
			return consumer.ConsumeSuccess, nil
		})
		if err != nil {
			return nil, err
		}
		// 启动
		if err := c.Start(); err != nil {
			return nil, err
		}
		consumers[subscriber.TopicName] = c
	}
	client.consumers = consumers
	client.handlers = handlers
	return &client, nil
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
