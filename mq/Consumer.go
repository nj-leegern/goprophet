package mq

/*
	消费者抽象接口
*/

type Consumer interface {
	// 消费消息
	ConsumeMsg(invokeMethod func(msg interface{}))
	// 释放资源
	Destroy() error
}
