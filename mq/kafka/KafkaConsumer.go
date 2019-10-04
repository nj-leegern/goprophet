package kafka

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"strings"
)

/*
	kafka消费者
*/

type KafkaConsumer struct {
	consumer *cluster.Consumer
}

/* 创建kafka消费者实例 */
func NewKafkaConsumer(conf KafkaConsumerConf) (*KafkaConsumer, error) {
	consumer, err := cluster.NewConsumer(strings.Split(conf.BrokerServers, ","), conf.GroupId, []string{conf.TopicName}, conf.Config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: consumer}, nil
}

/* 消费消息 */
func (c *KafkaConsumer) ConsumeMsg(invokeMethod func(msg interface{})) {
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
					invokeMethod(msg)
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
