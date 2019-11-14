package es

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/elastic/go-elasticsearch/v5"
	"net/http"
	"strings"
	"sync"
)

/*
	elastic search agent
*/

var (
	clients    sync.Map
	ErrEsAgent = fmt.Errorf("create elastic search client failed")
)

type EsAgent struct {
	EsClient *elasticsearch.Client
	mutex    sync.Mutex
}

/* 创建es agent实例 */
func NewEsAgent(ops ...EsOptions) (*EsAgent, error) {
	defaultOps := parseOptions(ops...)
	if len(defaultOps.addresses) == 0 {
		return nil, ErrEsAgent
	}
	params := make([]string, 0)
	for _, addr := range defaultOps.addresses {
		params = append(params, addr)
	}
	if len(defaultOps.username) > 0 {
		params = append(params, defaultOps.username)
	}
	agentKey := generateAgentKey(params)
	actual, loaded := clients.LoadOrStore(agentKey, &EsAgent{})
	if agent, ok := actual.(*EsAgent); ok {
		if loaded && agent.EsClient != nil {
			return agent, nil
		} else {
			agent.mutex.Lock()
			defer agent.mutex.Unlock()

			cfg := elasticsearch.Config{
				Addresses: defaultOps.addresses,
				Username:  defaultOps.username,
				Password:  defaultOps.password,
				Transport: &http.Transport{
					MaxIdleConnsPerHost:   int(defaultOps.maxIdleConns),
					MaxConnsPerHost:       int(defaultOps.maxConns),
					IdleConnTimeout:       defaultOps.idleTimeout,
					ResponseHeaderTimeout: defaultOps.responseHeaderTimeout,
					DialContext:           defaultOps.dialContext,
					TLSClientConfig:       &defaultOps.tlsClientConfig,
				},
			}
			es, err := elasticsearch.NewClient(cfg)
			if err != nil {
				return nil, err
			}
			agent.EsClient = es
			return agent, nil
		}
	}
	return nil, ErrEsAgent
}

// 解析配置项
func parseOptions(ops ...EsOptions) options {
	defaultOps := defaultEsOptions()
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
