package kafka

/*
	kafka生产者
*/
import (
	"github.com/Shopify/sarama"
	"github.com/nj-leegern/goprophet/mq"
	"time"
)

type KafkaProducer struct {
	producer       sarama.SyncProducer  // 同步
	asyncProducer  sarama.AsyncProducer // 异步
	callbackHandle func()               // 处理结果回调
}

/* 创建kafka生产者实例 */
func NewKafkaProducer(conf KafkaProducerConf) (*KafkaProducer, error) {
	if conf.Config == nil {
		conf.DefaultProducerConfig()
	}
	// 同步生产者
	producer, err := sarama.NewSyncProducer(conf.BrokerServers, conf.Config)
	if err != nil {
		return nil, err
	}
	// 异步生产者
	asyncProducer, er := sarama.NewAsyncProducer(conf.BrokerServers, conf.Config)
	if er != nil {
		return nil, er
	}
	return &KafkaProducer{producer: producer, asyncProducer: asyncProducer}, nil
}

/* 发送消息 */
func (p *KafkaProducer) SendSync(topic string, key string, msg []byte) (bool, error) {
	messages := make([]*sarama.ProducerMessage, 1, 1)

	m := sarama.ProducerMessage{}
	m.Topic = topic
	m.Key = sarama.ByteEncoder(key)
	m.Value = sarama.ByteEncoder(msg)
	m.Timestamp = time.Now()
	messages[0] = &m

	err := p.producer.SendMessages(messages)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* 异步发送消息 */
func (p *KafkaProducer) SendAsync(topic, key string, msg []byte, handleResult func(sendResult mq.SendResult, e error)) error {
	sign := make(chan int, 1)
	if p.callbackHandle == nil {
		p.callbackHandle = func() {
			sign <- 0
			for {
				select {
				case success := <-p.asyncProducer.Successes():
					key, _ := success.Key.Encode()
					val, _ := success.Value.Encode()
					result := mq.KafkaResult{
						Topic:     success.Topic,
						Key:       key,
						Value:     val,
						Offset:    success.Offset,
						Partition: success.Partition,
					}
					handleResult(mq.SendResult{KafkaResult: result}, nil)
				case pErr := <-p.asyncProducer.Errors():
					entity := pErr.Msg
					key, _ := entity.Key.Encode()
					val, _ := entity.Value.Encode()
					result := mq.KafkaResult{
						Topic:     entity.Topic,
						Key:       key,
						Value:     val,
						Offset:    entity.Offset,
						Partition: entity.Partition,
					}
					handleResult(mq.SendResult{KafkaResult: result}, pErr.Err)
				}
			}
		}
		go p.callbackHandle()
	}

	<-sign

	m := sarama.ProducerMessage{}
	m.Topic = topic
	m.Key = sarama.ByteEncoder(key)
	m.Value = sarama.ByteEncoder(msg)
	m.Timestamp = time.Now()

	p.asyncProducer.Input() <- &m

	return nil
}

/* 释放资源 */
func (p *KafkaProducer) Destroy() error {
	err := p.producer.Close()
	if err != nil {
		return err
	}
	err = p.asyncProducer.Close()
	if err != nil {
		return err
	}
	p.callbackHandle = nil
	return nil
}
