package mq

/*
	消费者抽象接口
*/

type Consumer interface {
	// 释放资源
	Destroy() error
}
