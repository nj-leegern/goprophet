package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
)

/*
	redis agent
*/

var (
	client         *redis.Client
	clusterClient  *redis.ClusterClient
	mutex          sync.Mutex
	ErrRedisAgent  = errors.New("create redis client failed")
	ErrRedisOption = errors.New("server address not null")
)

type RedisAgent struct {
	Client *redis.Client
}

type RedisClusterAgent struct {
	ClusterClient *redis.ClusterClient
}

/* 创建redis agent实例 */
func NewRedisAgent(ops ...Option) (*RedisAgent, error) {
	if client == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if client == nil {
			defaultOps := parseOptions(ops...)
			if len(defaultOps.serverAddrs) == 0 {
				return nil, ErrRedisOption
			}
			cc := redis.NewClient(&redis.Options{
				Addr:         defaultOps.serverAddrs[0],
				Password:     defaultOps.password,
				DB:           defaultOps.dbIndex,
				PoolSize:     int(defaultOps.maxPoolSize),
				MinIdleConns: int(defaultOps.minIdleConns),
				DialTimeout:  defaultOps.dialTimeout,
				IdleTimeout:  defaultOps.idleTimeout,
			})
			if cc == nil {
				return nil, ErrRedisAgent
			}
			err := cc.Ping().Err()
			if err != nil {
				return nil, err
			}
			client = cc
		}
	}
	return &RedisAgent{Client: client}, nil
}

/* 创建redis cluster agent实例 */
func NewRedisClusterAgent(ops ...Option) (*RedisClusterAgent, error) {
	if clusterClient == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if clusterClient == nil {
			defaultOps := parseOptions(ops...)
			cc := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:        defaultOps.serverAddrs,
				Password:     defaultOps.password,
				PoolSize:     int(defaultOps.maxPoolSize),
				MinIdleConns: int(defaultOps.minIdleConns),
				DialTimeout:  defaultOps.dialTimeout,
				IdleTimeout:  defaultOps.idleTimeout,
			})
			if cc == nil {
				return nil, ErrRedisAgent
			}
			_, err := cc.Ping().Result()
			if err != nil {
				return nil, err
			}
			clusterClient = cc
		}
	}
	return &RedisClusterAgent{ClusterClient: clusterClient}, nil
}

// 解析配置项
func parseOptions(ops ...Option) options {
	defaultOps := defaultRedisOptions()
	for _, apply := range ops {
		apply(&defaultOps)
	}
	return defaultOps
}
