package rocketmq

import (
	"git.oschina.net/cloudzone/cloudmq-go-client/cloudmq"
)

type DidaRocketConsumer struct {
	// nameSrv -> consumer
	clusterConsumer map[string]cloudmq.Consumer
	// topic -> handleMsg
	consumerHandler map[string]Subscription
}

/* 创建消费端实例 */
func NewDidaRocketConsumer(conf RocketConsumerConf) (*DidaRocketConsumer, error) {
	clusterConsumer := make(map[string]cloudmq.Consumer, len(conf.NameServers))
	for _, nameServer := range conf.NameServers {
		consumer, err := cloudmq.NewDefaultConsumer(conf.GroupName, &cloudmq.Config{
			Nameserver:   nameServer,
			InstanceName: "DEFAULT",
		})
		if err != nil {
			return nil, err
		}
		clusterConsumer[nameServer] = consumer
	}

	consumerHandler := make(map[string]Subscription, len(conf.Subscribers))
	for _, subscriber := range conf.Subscribers {
		consumerHandler[subscriber.TopicName] = subscriber
	}

	for _, consumer := range clusterConsumer {
		for _, subscriber := range conf.Subscribers {
			// 订阅消息
			consumer.Subscribe(subscriber.TopicName, "*")
		}
		consumer.RegisterMessageListener(func(msgs []*cloudmq.MessageExt) (int, error) {
			for _, msg := range msgs {
				if handler, ok := consumerHandler[msg.Topic]; ok {
					// 消息处理回调
					err := handler.HandleMsg([]string{string(msg.Body)})
					if err != nil {
						return cloudmq.Action.ReconsumeLater, err
					}
				}
			}
			return cloudmq.Action.CommitMessage, nil
		})
	}

	if len(clusterConsumer) > 0 {
		for _, consumer := range clusterConsumer {
			err := consumer.Start()
			if err != nil {
				return nil, err
			}
		}
	}
	return &DidaRocketConsumer{clusterConsumer: clusterConsumer}, nil
}

/* 释放资源 */
func (c *DidaRocketConsumer) Destroy() error {
	for _, consumer := range c.clusterConsumer {
		consumer.Shutdown()
	}
	return nil
}
