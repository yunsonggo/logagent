package service

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	elasticV7 "github.com/olivere/elastic/v7"
	"time"
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var ESClient *elasticV7.Client
var esAddr string

func InitESServer(addr string) (err error) {
	esAddr = addr
	cli, err := elasticV7.NewClient(elasticV7.SetSniff(false), elasticV7.SetURL(addr))
	if err != nil {
		logs.Warn("connect es server error", err)
		return
	}
	ESClient = cli
	logs.Info("es client init success")
	return
}

func SendMsgToES(topic string, data []byte) (err error) {
	msg := &LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)
	info,code,err := ESClient.Ping(esAddr).Do(context.Background())
	time.Sleep(time.Millisecond * 2)
	if err != nil {
		logs.Error("ping es server error:",err)
		return
	}
	logs.Info("es info :%v,code :%d",info,code)
	exists,existsErr := ESClient.IndexExists(topic).Do(context.Background())
	time.Sleep(time.Millisecond * 2)
	if existsErr != nil {
		logs.Error("check es index:%s error:%v",topic,err)
		return
	}
	if !exists {
		createIndex,createErr := ESClient.CreateIndex(topic).BodyJson(msg).Do(context.Background())
		time.Sleep(time.Millisecond * 2)
		if createErr != nil {
			fmt.Printf("%v",err)
			logs.Error("create es index error")
			return
		}
		if !createIndex.Acknowledged {
			logs.Error("create es index failed")
			return
		}
	}
	_, err = ESClient.Index().Index(topic).BodyJson(msg).Do(context.Background())
	if err != nil {
		logs.Error(err)
		logs.Debug("send msg to es faild,err :%v", err)
		return
	}
	logs.Info("send msg to es success")
	return
}
