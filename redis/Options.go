package redis

import (
	"runtime"
	"time"
)

/*
	redis配置项
*/

type options struct {
	serverAddrs  []string      // 服务器地址:端口
	password     string        // 密码
	dbIndex      int           // 数据库编号
	maxPoolSize  uint          // 最大连接数
	minIdleConns uint          // 最小空闲数
	dialTimeout  time.Duration // 请求超时时间
	idleTimeout  time.Duration // 连接空闲时间
}

type Option func(*options)

/* 服务器地址端口 */
func WithServerAddrs(serverAddrs []string) Option {
	return func(ops *options) {
		if len(serverAddrs) > 0 {
			ops.serverAddrs = serverAddrs
		}
	}
}

/* 密码 */
func WithPassword(password string) Option {
	return func(ops *options) {
		ops.password = password
	}
}

/* 数据库编号 */
func WithDBIndex(index int) Option {
	return func(ops *options) {
		if index >= 0 {
			ops.dbIndex = index
		}
	}
}

/* 最大连接数 */
func WithMaxPoolSize(maxPoolSize uint) Option {
	return func(ops *options) {
		ops.maxPoolSize = maxPoolSize
	}
}

/* 最小空闲数 */
func WithMinIdleConns(minIdleConns uint) Option {
	return func(ops *options) {
		ops.minIdleConns = minIdleConns
	}
}

/* 请求超时时间 */
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(ops *options) {
		if dialTimeout > 0 {
			ops.dialTimeout = dialTimeout
		}
	}
}

/* 连接空闲时间 */
func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(ops *options) {
		if idleTimeout > 0 {
			ops.idleTimeout = idleTimeout
		}
	}
}

// 默认配置项
func defaultRedisOptions() options {
	return options{
		dbIndex:      0,
		maxPoolSize:  uint(runtime.NumCPU()) * 10,
		minIdleConns: 10,
		dialTimeout:  5 * time.Second,
		idleTimeout:  3 * time.Minute,
	}
}
