package rocketmq

import (
	"git.oschina.net/cloudzone/cloudmq-go-client/cloudmq"
	"strings"
)

type DidaRocketConsumer struct {
	// nameSrv_topic -> consumer
	clusterConsumer map[string]cloudmq.Consumer
}

/* 创建消费端实例 */
func NewDidaRocketConsumer(conf RocketConsumerConf) (*DidaRocketConsumer, error) {
	clusterConsumer := make(map[string]cloudmq.Consumer, len(conf.NameServers))
	for _, subscriber := range conf.Subscribers {
		for _, nameServer := range conf.NameServers {
			consumer, err := cloudmq.NewDefaultConsumer(conf.GroupName, &cloudmq.Config{
				Nameserver:   nameServer,
				InstanceName: "DEFAULT",
			})
			if err != nil {
				return nil, err
			}
			// 订阅消息
			consumer.Subscribe(subscriber.TopicName, "*")
			consumer.RegisterMessageListener(func(msgs []*cloudmq.MessageExt) (int, error) {
				messages := make([]string, 0)
				for _, msg := range msgs {
					messages = append(messages, string(msg.Body))
				}
				// 消息处理回调
				err := subscriber.HandleMsg(messages)
				if err != nil {
					return cloudmq.Action.ReconsumeLater, err
				}
				return cloudmq.Action.CommitMessage, nil
			})
			key := strings.Join([]string{nameServer, subscriber.TopicName}, "_")
			clusterConsumer[key] = consumer
		}
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
