package main

import (
	"beegoDemo/dial"
	_ "beegoDemo/routers"
	"beegoDemo/service"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 配置日志
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"logs/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	logs.EnableFuncCallDepth(true)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 启用数据库
	mysqlAddr, err := beego.AppConfig.String("mysql_addr")
	if err != nil {
		logs.Warn("mysql address lost,please check the app config", err)
		return
	}
	err = dial.DbServer(mysqlAddr)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 启用ETCD
	etcdAddrs, err := beego.AppConfig.String("etcd_listen")
	if err != nil {
		logs.Warn("etcd address lost,please check the app config", err)
		return
	}
	err = dial.EtcdServer(etcdAddrs)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 收集ETCD保存的配置
	confs, err := service.CollectAllEtcdKeyConf()
	if err != nil {
		logs.Warn(err)
		return
	}
	// tail根据confs 展开收集工作
	chanSize, err := beego.AppConfig.Int("tail_chansize")
	if err != nil {
		logs.Warn(err)
		return
	}
	err = service.InitTail(confs, chanSize)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 开启kafka生产端
	kafkaAddr, err := beego.AppConfig.String("kafka_addr")
	if err != nil {
		logs.Warn(err)
		return
	}
	err = service.InitKafkaProducer(kafkaAddr)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 启动发送服务
	go service.SendMsgServer()
	// 注册kafka消费端
	topic, err := beego.AppConfig.String("topic")
	if err != nil {
		logs.Warn(err)
		return
	}
	err = service.InitKfkCS(kafkaAddr, topic)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 注册es
	esAddr, err := beego.AppConfig.String("es_addr")
	if err != nil {
		logs.Warn(err)
		return
	}
	err = service.InitESServer(esAddr)
	if err != nil {
		logs.Warn(err)
		return
	}
	// 消费kafka到ES
	err = service.WatchMsg()
	if err != nil {
		logs.Warn(err)
		return
	}
	beego.Run("192.168.1.102:8090")
}
