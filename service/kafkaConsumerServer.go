package service

import (
	"github.com/Shopify/sarama"
	"github.com/beego/beego/v2/core/logs"
	"strings"
	"sync"
)

type KafkaConsumerServer struct {
	Consumer sarama.Consumer
	Addr     string
	Topic    string
	Wg       sync.WaitGroup
}

var KFKCS *KafkaConsumerServer

func InitKfkCS(addr, topic string) (err error) {
	KFKCS = &KafkaConsumerServer{}
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("init kafka consumer server failed, err:%v", err)
		return
	}
	KFKCS.Consumer = consumer
	KFKCS.Addr = addr
	KFKCS.Topic = topic
	logs.Info("init kafka consumer success")
	return
}


