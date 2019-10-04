package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

/*
	kafka服务配置参数
*/
/* 生产者配置 */
type KafkaProducerConf struct {
	// 服务地址和端口
	BrokerServers string `json:"brokerServers"`
	// 主题名
	TopicName string `json:"topicName"`
	// 配置参数
	Config *sarama.Config `json:"config"`
}

/* 消费者配置 */
type KafkaConsumerConf struct {
	// 服务地址和端口
	BrokerServers string `json:"brokerServers"`
	// 主题名
	TopicName string `json:"topicName"`
	// 配置参数
	Config *cluster.Config `json:"config"`
	// 分组ID
	GroupId string `json:"groupId"`
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

	c.Config = config
}

/* 消费者默认配置 */
func (c *KafkaConsumerConf) DefaultConsumerConfig() {
	config := cluster.NewConfig()
	config.Group.Mode = cluster.ConsumerModePartitions

	c.Config = config
}
