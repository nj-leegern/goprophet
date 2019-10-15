package rocketmq

/*
	RocketMQ服务配置参数
*/

/* RocketMQ消费端配置 */
type RocketConsumerConf struct {
	NameServers   []string
	GroupName     string
	ConsumerOrder bool           // 顺序消费
	Subscribers   []Subscription // 订阅配置
}

/* 订阅者配置 */
type Subscription struct {
	TopicName string                    // 主题名称
	Tag       string                    // 标签
	HandleMsg func(msgs []string) error // 处理消息
}
