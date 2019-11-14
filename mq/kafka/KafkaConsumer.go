package kafka

import (
	"fmt"
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
)

/*
	kafka消费者
*/

type KafkaConsumer struct {
	consumer    *cluster.Consumer
	subscribers map[string]func(msg string) error
}

/* 创建kafka消费者实例 */
func NewKafkaConsumer(conf KafkaConsumerConf) (*KafkaConsumer, error) {
	if conf.Config == nil {
		conf.DefaultConsumerConfig()
	}
	topics := make([]string, 0)
	sbs := make(map[string]func(msg string) error, 0)
	for _, subscriber := range conf.Subscribers {
		if len(subscriber.TopicName) > 0 && subscriber.HandleMsg != nil {
			sbs[subscriber.TopicName] = subscriber.HandleMsg
			topics = append(topics, subscriber.TopicName)
		}
	}
	if len(topics) == 0 {
		return nil, fmt.Errorf("no subscribe topics")
	}
	// 创建消费实例
	consumer, err := cluster.NewConsumer(conf.BrokerServers, conf.GroupId, topics, conf.Config)
	if err != nil {
		return nil, err
	}

	kc := &KafkaConsumer{consumer: consumer, subscribers: sbs}
	// 启动消费
	go kc.start()

	return kc, nil
}

// 启动消费
func (c *KafkaConsumer) start() {
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume partitions
	for {
		select {
		case part, ok := <-c.consumer.Partitions():
			if !ok {
				return
			}
			// start a separate goroutine to consume messages
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					if invokeHandle, ok := c.subscribers[msg.Topic]; ok {
						// 消息回调
						er := invokeHandle(string(msg.Value))
						if er == nil {
							// 标记已处理
							c.consumer.MarkOffset(msg, "")
						}
					}
				}
			}(part)
		case <-signals:
			return
		}
	}
}

/* 释放资源 */
func (c *KafkaConsumer) Destroy() error {
	return c.consumer.Close()
}
