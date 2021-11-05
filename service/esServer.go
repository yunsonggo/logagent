package service

import (
	"context"
	"time"

	"github.com/beego/beego/v2/core/logs"
	elasticV7 "github.com/olivere/elastic/v7"
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var ESClient *elasticV7.Client

func InitESServer(addr string) (err error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = ESClient.Index().Index(topic).Type(topic).BodyJson(msg).Do(ctx)
	if err != nil {
		logs.Error(err)
		cancel()
		return
	}
	cancel()
	return
}
