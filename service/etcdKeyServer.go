package service

import (
	"beegoDemo/dial"
	"beegoDemo/models"
	"github.com/beego/beego/v2/core/logs"
)

type EtcdKeyServer interface {
	// 查询一条记录
	CheckEtcdKey(key string) (exist bool)
	// 插入一条记录
	InsertEtcdKey(key string) (err error)
	// 获取所有key
	EtcdKeyList() (list []models.EtcdKey,err error)
}

type etcdKeyServer struct {}

func NewEtcdKeyServer() EtcdKeyServer {
	return &etcdKeyServer{}
}

// 插入一条记录
func (eks *etcdKeyServer) InsertEtcdKey(key string) (err error) {
	exist := eks.CheckEtcdKey(key)
	if exist {
		logs.Info("insert etcdKey: %s was exist ",key)
		return
	}
	var etcdKey models.EtcdKey
	etcdKey.KeyName = key
	_,err = dial.O.Insert(&etcdKey)
	return
}

// 获取所有key
func (eks *etcdKeyServer) EtcdKeyList() (list []models.EtcdKey,err error) {
	qs := dial.O.QueryTable("etcd_key")
	_,err = qs.All(&list)
	return
}

// 判断一条记录是否存在
func (eks *etcdKeyServer) CheckEtcdKey(key string) (exist bool) {
	qs := dial.O.QueryTable("etcd_key")
	exist = qs.Filter("key_name",key).Exist()
	return
}
