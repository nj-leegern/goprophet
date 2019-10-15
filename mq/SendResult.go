package mq

/* 发送返回结果 */
type SendResult struct {
	KafkaResult  KafkaResult
	RocketResult RocketMqResult
}

/* kafka返回结果 */
type KafkaResult struct {
	Topic     string
	Key       []byte
	Value     []byte
	Offset    int64
	Partition int32
}

/* RocketMQ返回结果 */
type RocketMqResult struct {
	MsgID         string
	QueueOffset   int64
	TransactionID string
	OffsetMsgID   string
	Topic         string
	BrokerName    string
	QueueId       int
}
