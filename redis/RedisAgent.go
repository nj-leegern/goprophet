package redis

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis"
	"strings"
	"sync"
)

/*
	redis agent
*/

const (
	TYPE_SINGLE  = "single"
	TYPE_CLUSTER = "cluster"
)

var (
	clients        sync.Map
	ErrRedisAgent  = errors.New("create redis client failed")
	ErrRedisOption = errors.New("server address not null")
)

type RedisAgent struct {
	Client *redis.Client
	mutex  sync.Mutex
}

type RedisClusterAgent struct {
	ClusterClient *redis.ClusterClient
	mutex         sync.Mutex
}

/* 创建redis agent实例 */
func NewRedisAgent(ops ...Option) (*RedisAgent, error) {
	// 创建单实例客户端
	agent, err := createAgent(TYPE_SINGLE, ops...)
	if err != nil {
		return nil, err
	}
	if redisAgent, ok := agent.(*RedisAgent); ok {
		return redisAgent, nil
	}
	return nil, ErrRedisAgent
}

/* 创建redis cluster agent实例 */
func NewRedisClusterAgent(ops ...Option) (*RedisClusterAgent, error) {
	// 创建集群客户端
	agent, err := createAgent(TYPE_CLUSTER, ops...)
	if err != nil {
		return nil, err
	}
	if clusterAgent, ok := agent.(*RedisClusterAgent); ok {
		return clusterAgent, nil
	}
	return nil, ErrRedisAgent
}

// 创建redis客户端
func createAgent(agentType string, ops ...Option) (interface{}, error) {
	defaultOps := parseOptions(ops...)
	if len(defaultOps.serverAddrs) == 0 {
		return nil, ErrRedisOption
	}
	params := make([]string, 0)
	for _, addr := range defaultOps.serverAddrs {
		params = append(params, addr)
	}
	params = append(params, agentType)
	agentKey := generateAgentKey(params)
	var agent interface{}
	if agentType == TYPE_CLUSTER {
		agent = &RedisClusterAgent{}
	} else {
		agent = &RedisAgent{}
	}
	actual, loaded := clients.LoadOrStore(agentKey, agent)
	// single node
	if agent, ok := actual.(*RedisAgent); ok {
		if loaded && agent.Client != nil {
			return agent, nil
		} else {
			agent.mutex.Lock()
			defer agent.mutex.Unlock()

			cc := redis.NewClient(&redis.Options{
				Addr:         defaultOps.serverAddrs[0],
				Password:     defaultOps.password,
				DB:           defaultOps.dbIndex,
				PoolSize:     int(defaultOps.maxPoolSize),
				MinIdleConns: int(defaultOps.minIdleConns),
				DialTimeout:  defaultOps.dialTimeout,
				IdleTimeout:  defaultOps.idleTimeout,
				MaxRetries:   2,
			})
			if cc == nil {
				return nil, ErrRedisAgent
			}
			err := cc.Ping().Err()
			if err != nil {
				return nil, err
			}
			agent.Client = cc
			return agent, nil
		}
	}
	// cluster node
	if agent, ok := actual.(*RedisClusterAgent); ok {
		if loaded && agent.ClusterClient != nil {
			return agent, nil
		} else {
			agent.mutex.Lock()
			defer agent.mutex.Unlock()

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
			agent.ClusterClient = cc
			return agent, nil
		}
	}

	return nil, ErrRedisAgent
}

// 解析配置项
func parseOptions(ops ...Option) options {
	defaultOps := defaultRedisOptions()
	for _, apply := range ops {
		apply(&defaultOps)
	}
	return defaultOps
}

// 生成agent实例KEY
func generateAgentKey(params []string) string {
	data := md5.Sum([]byte(strings.Join(params, "_")))
	return hex.EncodeToString(data[:])
}
