package service

import (
	"github.com/Shopify/sarama"
	"github.com/beego/beego/v2/core/logs"
)
// 消费kafka 消息 到ES
func WatchMsg() (err error) {
	partitions,err := KFKCS.Consumer.Partitions(KFKCS.Topic)
	if err != nil {
		logs.Error("kafka consumer server's partitions lost,err:",err)
		return
	}
	for partition := range partitions {
		pc,err := KFKCS.Consumer.ConsumePartition(KFKCS.Topic,int32(partition),sarama.OffsetNewest)
		if err != nil {
			logs.Error("kafka consumer partition %d ,err : %s\n",partition,err)
			continue
		}
		defer pc.AsyncClose()
		go func(pc sarama.PartitionConsumer) {
			KFKCS.Wg.Add(1)
			for msg := range pc.Messages() {
				err = SendMsgToES(KFKCS.Topic,msg.Value)
				if err != nil {
					logs.Warn("send to es failed, err:%v", err)
					continue
				}
				logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
			KFKCS.Wg.Done()
		}(pc)
	}
	KFKCS.Wg.Wait()
	return
}
