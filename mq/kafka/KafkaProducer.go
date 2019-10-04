package kafka

/*
	kafka生产者
*/
import (
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

type KafkaProducer struct {
	producer  sarama.SyncProducer
	TopicName string
}

/* 创建kafka生产者实例 */
func NewKafkaProducer(conf KafkaProducerConf) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(strings.Split(conf.BrokerServers, ","), conf.Config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer, TopicName: conf.TopicName}, nil
}

/* 批量发送消息 */
func (p *KafkaProducer) BatchSend(key []byte, msgs [][]byte) (bool, error) {
	messages := make([]*sarama.ProducerMessage, len(msgs))
	for i, msg := range msgs {
		m := sarama.ProducerMessage{}
		m.Topic = p.TopicName
		m.Key = sarama.ByteEncoder(key)
		m.Value = sarama.ByteEncoder(msg)
		m.Timestamp = time.Now()
		messages[i] = &m
	}
	err := p.producer.SendMessages(messages)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* 释放资源 */
func (p *KafkaProducer) Destroy() error {
	return p.producer.Close()
}
