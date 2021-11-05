package service

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

func SendMsgServer() {
	fmt.Printf("send msg server forever running....\n")
	logs.Warn("send msg server forever running....")
	for {
		msg := GetOneLine()
		err := SendToKafka(msg.Msg,msg.Topic)
		if err != nil {
			logs.Error("send to kafka failed, err:%v", err)
			time.Sleep(time.Second * 2)
			continue
		}
		logs.Info("send a message ago")
	}
}

