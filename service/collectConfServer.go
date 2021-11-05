package service

import (
	"beegoDemo/dial"
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"go.etcd.io/etcd/api/v3/mvccpb"

	"time"
)

var eks = NewEtcdKeyServer()

// 根据数据库现有etcdKey收集配置列表
func CollectAllEtcdKeyConf() (collectConf []CollectConf,err error) {
	etcdAddrs,err := beego.AppConfig.String("etcd_listen")
	if err != nil {
		logs.Error("collect conf etcdAddrs err:%v", err)
	}
	etcdKeys,err := eks.EtcdKeyList()
	for _, key := range etcdKeys {
		cList,err := CollectOneEtcdKeyConf(key.KeyName)
		if err != nil {
			logs.Error("collectOne conf failed,key:%s, err:%v",key.KeyName, err)
			continue
		}
		collectConf = append(collectConf, cList...)
		time.Sleep(time.Second)
		// 每完成一个去监听一个
		go WatchEtcdKey(key.KeyName,etcdAddrs)
	}
	return
}

// 收集单个etcdKey的配置列表
func CollectOneEtcdKeyConf(etcdKey string) (collectConf []CollectConf,err error) {
	ctx,cancel := context.WithTimeout(context.Background(),time.Second * 5)
	resp,err := dial.EC.Get(ctx,etcdKey)
	if err != nil {
		logs.Error(err)
		cancel()
		return nil, err
	}
	cancel()
	for _,v := range resp.Kvs {
		if string(v.Key) == etcdKey {
			err = json.Unmarshal(v.Value,&collectConf)
			if err != nil {
				logs.Error("unmarshal failed, err:%v", err)
				continue
			}
		}
	}
	return
}

func WatchEtcdKey(key ,etcdAddrs string) {
	cli,err := dial.EtcdWatchServer(etcdAddrs)
	if err != nil {
		logs.Error("watcher client connect etcd failed, err:", err)
		return
	}
	for {
		rch := cli.Watch(context.Background(),key)
		var collectConf []CollectConf
		var getConfSucc = true
		for wresp := range rch {
			for _,ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value,&collectConf)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", collectConf)
				// 通知tail 更新 配置
				_ = UpdateConfig(collectConf)
			}
		}
	}
}