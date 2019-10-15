package kafka

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
)

/*
	kafka消费者
*/

type KafkaConsumer struct {
	consumer *cluster.Consumer
}

/* 创建kafka消费者实例 */
func NewKafkaConsumer(conf KafkaConsumerConf) (*KafkaConsumer, error) {
	if conf.Config == nil {
		conf.DefaultConsumerConfig()
	}
	consumer, err := cluster.NewConsumer(conf.BrokerServers, conf.GroupId, conf.TopicNames, conf.Config)
	if err != nil {
		return nil, err
	}
	kc := &KafkaConsumer{consumer: consumer}
	// 启动消费
	go kc.start(conf.HandleMsg)

	return kc, nil
}

// 启动消费
func (c *KafkaConsumer) start(invokeHandle func(msg string) error) {
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
					er := invokeHandle(string(msg.Value))
					if er == nil {
						c.consumer.MarkOffset(msg, "")
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
