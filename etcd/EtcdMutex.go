package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/pkg/errors"
	"time"
)

/*
	ETCD分布式锁
*/

type EtcdMutex struct {
	session *concurrency.Session
	mutex   *concurrency.Mutex
}

/* 创建ETCD客户端实例 */
func NewEtcdClient(addrs []string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
}

/* 创建锁实例 */
func NewEtcdMutex(client *clientv3.Client, key string) (*EtcdMutex, error) {
	s, err := concurrency.NewSession(client)
	if err != nil {
		return nil, err
	}
	m := concurrency.NewMutex(s, key)
	if m == nil {
		return nil, errors.New("create mutex error")
	}
	return &EtcdMutex{session: s, mutex: m}, nil
}

/* 获取锁 */
func (e *EtcdMutex) TryLock(timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := e.mutex.Lock(ctx); err != nil {
		return false, err
	}
	return true, nil
}

/* 释放锁 */
func (e *EtcdMutex) Unlock() (err error) {
	if err = e.mutex.Unlock(context.TODO()); err != nil {
		return err
	}
	if err = e.session.Close(); err != nil {
		return err
	}
	return
}
