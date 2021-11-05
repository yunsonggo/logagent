package service

import (
	"github.com/Shopify/sarama"
	"github.com/beego/beego/v2/core/logs"
)

var KafkaProducer sarama.SyncProducer

func InitKafkaProducer(addr string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	cli , err := sarama.NewSyncProducer([]string{addr},config)
	if err != nil {
		logs.Error("init kafka producer failed, err:", err)
		return
	}
	KafkaProducer = cli
	logs.Debug("init kafka succ")
	return
}

func SendToKafka(data,topic string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	pid,offset,err := KafkaProducer.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err:%v data:%v topic:%v", err, data, topic)
		return
	}
	logs.Info("send succ, pid:%v offset:%v, topic:%v\n", pid, offset, topic)
	return
}