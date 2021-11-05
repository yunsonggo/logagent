package dial

import (
	"beegoDemo/models"
	"context"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	clientV3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

var (
	O orm.Ormer
	EC *clientV3.Client
)

func DbServer(mysqlAddr string) (err error) {
	_ = orm.RegisterDriver("mysql",orm.DRMySQL)
	_ = orm.RegisterDataBase("default","mysql",mysqlAddr)
	err = commandDbTables()
	o := orm.NewOrm()
	O = o
	return
}

func commandDbTables() error {
	orm.RegisterModel(new(models.AppInfo),new(models.LogInfo),new(models.EtcdKey))
	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	orm.RunCommand()
	if err != nil {
		logs.Warn(err)
	}
	return err
}

func EtcdServer(etcdAddrs string) (err error) {
	add := strings.Split(etcdAddrs,",")
	cli,err := clientV3.New(clientV3.Config{
		Endpoints: add,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Warn("connect etcd server fialed:",err)
		return
	}
	EC = cli
	return
}

// 根据key 查询value
func FindEtcdResp(key string) (*clientV3.GetResponse,error) {
	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	resp,err := EC.Get(ctx,key)
	cancel()
	return resp,err
}

func EtcdWatchServer(etcdAddrs string) (*clientV3.Client,error) {
	add := strings.Split(etcdAddrs,",")
	wClient,err := clientV3.New(clientV3.Config{
		Endpoints: add,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Warn("connect etcd server fialed:",err)
		return nil,err
	}
	return wClient,nil
}