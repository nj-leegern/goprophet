package es

import (
	"context"
	"crypto/tls"
	"net"
	"runtime"
	"time"
)

/*
	ES配置项
*/

type options struct {
	addresses             []string      // 服务器地址:端口
	username              string        // 用户名
	password              string        // 密码
	maxConns              int16         // 最大连接数
	maxIdleConns          int16         // 最大空闲连接数
	idleTimeout           time.Duration // 连接空闲时间
	responseHeaderTimeout time.Duration // 请求响应超时时间
	dialContext           func(ctx context.Context, network, addr string) (net.Conn, error)
	tlsClientConfig       tls.Config
}

type EsOptions func(ops *options)

/* 服务器地址端口 */
func WithAddresses(addrs []string) EsOptions {
	return func(ops *options) {
		if len(addrs) > 0 {
			ops.addresses = addrs
		}
	}
}

/* 用户名 */
func WithUsername(username string) EsOptions {
	return func(ops *options) {
		if len(username) > 0 {
			ops.username = username
		}
	}
}

/* 密码 */
func WithPassword(password string) EsOptions {
	return func(ops *options) {
		if len(password) > 0 {
			ops.password = password
		}
	}
}

/* 最大连接数 */
func WithMaxConns(maxConns int16) EsOptions {
	return func(ops *options) {
		ops.maxConns = maxConns
	}
}

/* 最大空闲连接数 */
func WithMaxIdleConns(maxIdleConns int16) EsOptions {
	return func(ops *options) {
		ops.maxIdleConns = maxIdleConns
	}
}

/* 连接超时时间 */
func WithIdleTimeout(idleTimeout time.Duration) EsOptions {
	return func(ops *options) {
		ops.idleTimeout = idleTimeout
	}
}

/* 请求响应空闲时间 */
func WithResponseHeaderTimeout(responseHeaderTimeout time.Duration) EsOptions {
	return func(ops *options) {
		ops.responseHeaderTimeout = responseHeaderTimeout
	}
}

// ES默认配置项
func defaultEsOptions() options {
	return options{
		maxConns:              int16(runtime.NumCPU() * 10),
		maxIdleConns:          10,
		idleTimeout:           5 * time.Minute,
		responseHeaderTimeout: 3 * time.Second,
		dialContext:           (&net.Dialer{Timeout: 2 * time.Second}).DialContext,
		tlsClientConfig:       tls.Config{MinVersion: tls.VersionTLS11},
	}
}
