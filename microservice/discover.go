package microservice

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

type DiscoverService struct {
	client *clientv3.Client
	serviceList map[string]string
	Lock sync.Mutex
}

func NewDiscoverService(addr []string) (*DiscoverService,error) {
	clientConf := clientv3.Config{
		Endpoints: addr,
		DialTimeout: time.Second * 5,
	}

	client,err := clientv3.New(clientConf)
	if err != nil {
		return nil,err
	}

	dis := DiscoverService{
		client: client,
		serviceList: make(map[string]string),
	}

	return &dis,nil
}

func (this *DiscoverService)GetService(prefix string) ([]string,error) {
	resp,err := this.client.Get(context.TODO(),prefix,clientv3.WithPrefix())
	if err != nil {
		return nil,err
	}
	addrs := this.extracAddr(resp,prefix)

	go this.watcher(prefix)

	return addrs,nil
}

func (this *DiscoverService)extracAddr(resp *clientv3.GetResponse,prefix string) []string {
	if resp == nil || resp.Kvs == nil {
		return nil
	}
	var addrs = make([]string,len(resp.Kvs))
	for _,v := range resp.Kvs {
		if v != nil {
			this.setServiceList(string(v.Key),string(v.Value))
			addrs = append(addrs,string(v.Value))
		}
	}

	return addrs
}

func (this *DiscoverService)watcher(prefix string)  {
	resp := this.client.Watch(context.TODO(),prefix,clientv3.WithPrefix())

	for v := range resp {
		for _,vv := range v.Events {
			switch vv.Type {
			case mvccpb.PUT:
				this.setServiceList(string(vv.Kv.Key),string(vv.Kv.Value))
			case mvccpb.DELETE:
				this.delServiceList(string(vv.Kv.Key))
			default:
				log.Println(vv.Type)
			}
		}
	}
}

func (this *DiscoverService)setServiceList(key,value string)  {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.serviceList[key] = value
	log.Println("set key:",key,",value:",value)
}


func (this *DiscoverService)delServiceList(key string)  {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	delete(this.serviceList,key)
	log.Println("del key:",key)
}