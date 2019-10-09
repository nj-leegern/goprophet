package etcd

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

/*
	ETCD服务发现
*/

type EtcdDiscovery struct {
	client     *clientv3.Client
	services   map[string]string // 注册的服务列表
	lock       sync.Mutex
	prefixName string // 服务名前缀
}

/* 创建服务发现实例 */
func NewEtcdDiscovery(addrs []string, prefixName string) (*EtcdDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	disc := &EtcdDiscovery{
		client:     client,
		services:   make(map[string]string, 0),
		prefixName: prefixName,
	}
	// 监听注册服务状态
	err = disc.watcher(prefixName)
	if err != nil {
		return nil, err
	}
	return disc, nil
}

/* 获取注册的服务列表 */
func (e *EtcdDiscovery) GetServices() map[string]string {
	return e.services
}

/* 实时获取注册的服务 */
func (e *EtcdDiscovery) GetServicesNow() (map[string]string, error) {
	return e.get()
}

/* 释放资源 */
func (e *EtcdDiscovery) Destroy() (bool, error) {
	err := e.client.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 监听服务状态
func (e *EtcdDiscovery) watcher(prefixName string) error {
	// 获取注册的服务
	services, err := e.get()
	if err != nil {
		return err
	}

	e.services = services

	go func() {
		for {
			select {
			case resp := <-e.client.Watch(context.Background(), prefixName, clientv3.WithPrefix()):
				if nil != resp.Err() {
					return
				}
				for _, ev := range resp.Events {
					switch ev.Type {
					case mvccpb.PUT:
						e.putService(string(ev.Kv.Key), string(ev.Kv.Value))
					case mvccpb.DELETE:
						e.removeService(string(ev.Kv.Key))
					}
				}
			}
		}
	}()

	return nil
}

// 获取注册的服务
func (e *EtcdDiscovery) get() (map[string]string, error) {
	rsp, err := e.client.Get(context.Background(), e.prefixName, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	services := make(map[string]string, 0)
	if rsp != nil && nil != rsp.Kvs {
		for _, v := range rsp.Kvs {
			if v != nil {
				services[string(v.Key)] = string(v.Value)
			}
		}
	}
	return services, nil
}

// 缓存服务
func (e *EtcdDiscovery) putService(key, val string) bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.services[key] = val
	return true
}

// 清理服务
func (e *EtcdDiscovery) removeService(key string) bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.services, key)
	return true
}
