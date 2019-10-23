package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go"
	"github.com/apache/rocketmq-client-go/primitive"
	"github.com/apache/rocketmq-client-go/producer"
	"github.com/nj-leegern/goprophet/mq"
)

/*
	RocketMQ生产者
*/

type RocketProducer struct {
	p rocketmq.Producer
}

/* 创建RocketMQ生产者实例 */
func NewRocketProducer(nameServers []string, retries int) (*RocketProducer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(nameServers),
		producer.WithRetry(retries),
	)
	if err != nil {
		return nil, err
	}
	err = p.Start()
	if err != nil {
		return nil, err
	}
	rp := &RocketProducer{
		p: p,
	}
	return rp, nil
}

/* 发送消息 */
func (p *RocketProducer) SendSync(topic string, tag string, msg []byte) (bool, error) {
	m := primitive.NewMessage(topic, msg)
	if len(tag) > 0 {
		// 目前版本不支持TAG
		//m.WithTag(tag)  // expression
	}
	_, err := p.p.SendSync(context.Background(), m)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* 异步发送消息 */
func (p *RocketProducer) SendAsync(topic, tag string, msg []byte, handleResult func(sendResult mq.SendResult, e error)) error {
	m := primitive.NewMessage(topic, msg)
	if len(tag) > 0 {
		// 目前版本不支持TAG
		//m.WithTag(tag)  // expression
	}
	return p.p.SendAsync(context.Background(), func(i context.Context, result *primitive.SendResult, e error) {
		rs := mq.RocketMqResult{
			MsgID:         result.MsgID,
			QueueOffset:   result.QueueOffset,
			TransactionID: result.TransactionID,
			OffsetMsgID:   result.OffsetMsgID,
		}
		queue := result.MessageQueue
		if queue != nil {
			rs.Topic = queue.Topic
			rs.BrokerName = queue.BrokerName
			rs.QueueId = queue.QueueId
		}
		handleResult(mq.SendResult{RocketResult: rs}, e)
	}, m)
}

/* 释放资源 */
func (p *RocketProducer) Destroy() error {
	return p.p.Shutdown()
}