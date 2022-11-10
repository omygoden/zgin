package microservice

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type RegisterService struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancelHandle  func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
}

func NewRegisterService(addr []string, timeNum int64) (*RegisterService, error) {
	clientConf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Second * 5,
	}

	client, err := clientv3.New(clientConf)
	if err != nil {
		return nil, err
	}
	reg := RegisterService{
		client: client,
	}
	err = reg.setLease(timeNum)

	if err != nil {
		return nil, err
	}
	go reg.ListenLeaseRespChan()

	return &reg, nil
}

//设置租约
func (this *RegisterService) setLease(timeNum int64) error {
	lease := clientv3.NewLease(this.client)

	//设置续租时间
	leaseResp, err := lease.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	//设置取消事件
	ctx, c := context.WithCancel(context.TODO())
	this.cancelHandle = c

	//设置心跳
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}
	this.keepAliveChan = leaseRespChan
	this.lease = lease
	this.leaseResp = leaseResp

	return nil
}

//监听
func (this *RegisterService) ListenLeaseRespChan() {
	for {
		//log.Println(<-this.keepAliveChan)
		//time.Sleep(time.Second)
		select {
		case v := <-this.keepAliveChan:
			log.Println("chan:", v)
			if v == nil {
				log.Println("已关闭续租功能")
				return
			} else {
				log.Println("续租成功")
			}
		}
	}
}

//通过续租，注册服务
func (this *RegisterService) PutService(key, val string) error {
	kv := clientv3.NewKV(this.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID))
	return err
}

func (this *RegisterService) GetService(key string) error {
	kv := clientv3.NewKV(this.client)
	resp,err := kv.Get(context.TODO(),key)
	log.Println("get:",err)
	log.Println("getvalue:",string(resp.Kvs[0].Value))
	return err
}

//撤销租约
func (this *RegisterService) RevokeLease() error {
	this.cancelHandle()
	time.Sleep(time.Second * 2)
	_, err := this.lease.Revoke(context.TODO(), this.leaseResp.ID)

	return err
}