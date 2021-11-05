package service

import (
	"beegoDemo/dial"
	"beegoDemo/models"
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"time"
)

type LogServer interface {
	// 所有日志配置
	LogList() (list []models.LogInfo,err error)
	// 插入一条
	InsertOneInfo(info models.LogInfo) (id int64,err error)
	// 写入etcd配置
	InsertInfoToEtcd(exist bool,etcdKey string,info models.LogInfo) (err error)
}

type logServer struct {}

func NewLogServer() LogServer {
	return &logServer{}
}

// 所有日志配置
func (ls *logServer) LogList() (list []models.LogInfo,err error) {
	qs := dial.O.QueryTable("log_info")
	_,err = qs.All(&list)
	return
}
// 插入一条
func (ls *logServer) InsertOneInfo (info models.LogInfo) (id int64,err error) {
	var i models.LogInfo
	i.AppName = info.AppName
	i.LogPath = info.LogPath
	i.Topic = info.Topic
	id,err = dial.O.Insert(&i)
	return
}

// 写入etcd配置
func (ls *logServer) InsertInfoToEtcd(exist bool,etcdKey string,info models.LogInfo) (err error) {
	var logConfArr []CollectConf
	// 如果etcdkey 已经存在
	if exist {
		resp,err := dial.FindEtcdResp(etcdKey)
		if err != nil {
			logs.Warn("etcd key :%s,find resp err:%s",etcdKey,err)
		} else {
			for _, v := range resp.Kvs {
				if string(v.Key) == etcdKey {
					err = json.Unmarshal(v.Value,&logConfArr)
					if err != nil {
						logs.Error("Unmarshal resp.Kvs.v.value err:",err)
						continue
					}
					break
				}
			}
		}
	}
	conf := CollectConf{
		LogPath: info.LogPath,
		Topic: info.Topic,
	}
	logConfArr = append(logConfArr,conf)
	data,err := json.Marshal(logConfArr)
	if err != nil {
		logs.Warn("json marshal err :",err)
		return
	}
	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	_,err = dial.EC.Put(ctx,etcdKey,string(data))
	cancel()
	if err != nil {
		logs.Warn("put etcd err :",err)
		return
	}
	// 写入成功,该etcdKey如果本来就存在,更新配置后etcd 的 watch 会通知tail
	// 此处不用处理
	// 如果该key不存在,则etcd 的 watch 没有监听
	//......此处需要处理新入的key-value 加入监听队列并通知tail
	if !exist {
		// 推送tail
		_ = UpdateConfig(logConfArr)
		// 监听
		etcdAddrs,err := beego.AppConfig.String("etcd_listen")
		if err != nil {
			logs.Warn("watch etcd get etcdAddrs err :",err)
		}
		logs.Info("log server etcdAddrs:%v",etcdAddrs)
		go WatchEtcdKey(etcdKey,etcdAddrs)
	}
	logs.Debug("put etcd succ, data:%v", string(data))
	return
}