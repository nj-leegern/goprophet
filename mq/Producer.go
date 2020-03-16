package mq

/*
	生产者抽象接口
*/

type Producer interface {
	// 发送消息
	SendSync(topic, sign string, msg []byte) (SendResult, error)
	// 异步发送消息
	SendAsync(topic, sign string, msg []byte, handleResult func(sendResult SendResult, e error)) error
	// 释放资源
	Destroy() error
}
