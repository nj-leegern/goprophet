package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"time"
)

/*
	kafka服务配置参数
*/
/* 生产者配置 */
type KafkaProducerConf struct {
	// 服务地址和端口
	BrokerServers []string
	// 配置参数
	Config *sarama.Config
}

/* 消费者配置 */
type KafkaConsumerConf struct {
	// 服务地址和端口
	BrokerServers []string
	// 主题名
	TopicNames []string
	// 配置参数
	Config *cluster.Config
	// 分组ID
	GroupId string
	// 处理消息
	HandleMsg func(msg string) error
}

/* 生产者默认配置 */
func (c *KafkaProducerConf) DefaultProducerConfig() {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true
	// 返回失败通知
	config.Producer.Return.Errors = true
	// 超时时间默认10sec
	config.Producer.Timeout = 8 * time.Second

	c.Config = config
}

/* 消费者默认配置 */
func (c *KafkaConsumerConf) DefaultConsumerConfig() {
	config := cluster.NewConfig()
	// 分区消费模式
	config.Group.Mode = cluster.ConsumerModePartitions

	c.Config = config
}
