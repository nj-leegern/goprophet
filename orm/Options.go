package orm

import (
	"runtime"
	"time"
)

/*
	数据源配置
*/

const (
	DIALECT_MYSQL   = "mysql"
	DIALECT_SQLITE3 = "sqlite3"
	DIALECT_PG      = "postgres"
)

type options struct {
	hostname        string        // 服务器主机名:端口
	database        string        // 数据库名
	username        string        // 用户名
	password        string        // 密码
	maxIdleConns    int           // 最大空闲数
	maxOpenConns    int           // 最大连接数
	connMaxLifetime time.Duration // 连接生命周期
}

type Option func(ops *options)

// 数据源默认配置
func defaultDBOptions() options {
	return options{
		maxIdleConns:    10,
		maxOpenConns:    runtime.NumCPU() * 10,
		connMaxLifetime: 10 * time.Minute,
	}
}

/* 服务器主机名 */
func WithHostname(hostname string) Option {
	return func(ops *options) {
		if len(hostname) > 0 {
			ops.hostname = hostname
		}
	}
}

/* 数据库名 */
func WithDatabase(database string) Option {
	return func(ops *options) {
		if len(database) > 0 {
			ops.database = database
		}
	}
}

/* 用户名 */
func WithUsername(username string) Option {
	return func(ops *options) {
		if len(username) > 0 {
			ops.username = username
		}
	}
}

/* 密码 */
func WithPassword(password string) Option {
	return func(ops *options) {
		if len(password) > 0 {
			ops.password = password
		}
	}
}

/* 最大空闲数 */
func WithMaxIdleConns(maxIdleConns uint) Option {
	return func(ops *options) {
		if maxIdleConns > 0 {
			ops.maxIdleConns = int(maxIdleConns)
		}
	}
}

/* 最大连接数 */
func WithMaxOpenConns(maxOpenConns uint) Option {
	return func(ops *options) {
		if maxOpenConns > 0 {
			ops.maxOpenConns = int(maxOpenConns)
		}
	}
}

/* 连接生命周期 */
func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(ops *options) {
		if connMaxLifetime > 0 {
			ops.connMaxLifetime = connMaxLifetime
		}
	}
}
