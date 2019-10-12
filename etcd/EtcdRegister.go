package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"time"
)

/*
	ETCD服务注册
*/

type EtcdRegister struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	cancelFunc    func() // 关闭续租回调
}

/* 创建服务注册实例 */
func NewEtcdRegister(addrs []string, leaseTime int64) (*EtcdRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	reg := &EtcdRegister{client: client}
	// 设置租约
	if rst, err := reg.setLease(leaseTime); !rst {
		return nil, err
	}
	return reg, nil
}

/* 以租约方式注册服务 */
func (e *EtcdRegister) Register(key string, val string) (bool, error) {
	kv := clientv3.NewKV(e.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(e.leaseResp.ID))
	if err != nil {
		return false, err
	}
	return true, nil
}

/* 撤销租约 */
func (e *EtcdRegister) RevokeLease() (bool, error) {
	e.cancelFunc()
	time.Sleep(1 * time.Second)
	_, err := e.lease.Revoke(context.TODO(), e.leaseResp.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* 监听续租情况 */
func (e *EtcdRegister) AddLeaseListener(invokeHandler func(resp clientv3.LeaseKeepAliveResponse)) {
	go func() {
		for {
			select {
			case rsp := <-e.keepAliveChan:
				invokeHandler(*rsp)
				// 续约失败退出
				if rsp == nil {
					goto END
				}
			}
			time.Sleep(1 * time.Second)
		}
	END:
	}()
}

/* 释放资源 */
func (e *EtcdRegister) Destroy() (bool, error) {
	err := e.client.Close()
	if err != nil {
		return false, err
	}
	return true, err
}

// 设置租约
func (e *EtcdRegister) setLease(leaseTime int64) (bool, error) {
	lease := clientv3.NewLease(e.client)
	// 设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), leaseTime)
	if err != nil {
		return false, err
	}
	ctx, cancelFunc := context.WithCancel(context.TODO())
	// 自动续租
	keepAliveChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return false, err
	}
	e.lease = lease
	e.leaseResp = leaseResp
	e.cancelFunc = cancelFunc
	e.keepAliveChan = keepAliveChan
	return true, nil
}
