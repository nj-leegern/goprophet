package mq

/*
	生产者抽象接口
*/

type Producer interface {
	// 批量发送消息
	BatchSend(key []byte, msgs [][]byte) (bool, error)
	// 释放资源
	Destroy() error
}
