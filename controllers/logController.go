package controllers

import (
	"beegoDemo/models"
	"beegoDemo/service"
	"beegoDemo/tools"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

type LogController struct {
	beego.Controller
}

var ls = service.NewLogServer()
var eks = service.NewEtcdKeyServer()
// 日志列表
func (l *LogController) LogList() {
	l.Layout = "layout/layout.html"
	list,err := ls.LogList()
	if err != nil {
		logs.Warn(err)
		l.Data["Error"] = err.Error()
		l.TplName = "error/error.html"
		return
	}
	l.Data["LogList"] = list
	l.TplName = "log/list.html"
	return
}

// 日志申请
func (l *LogController) LogApply () {
	l.Layout = "layout/layout.html"
	l.TplName = "log/apply.html"
}

// 提交申请
func (l *LogController) LogCreate () {
	info := models.LogInfo{}
	l.Layout = "layout/layout.html"
	err := l.ParseForm(&info)
	if err != nil || len(info.AppName) == 0 || len(info.LogPath) == 0 || len(info.Topic) == 0 {
		logs.Warn(err)
		l.Data["Error"] = err.Error()
		l.TplName = "error/error.html"
		return
	}
	_,err = ls.InsertOneInfo(info)
	if err != nil {
		logs.Warn(err)
		l.Data["Error"] = err.Error()
		l.TplName = "error/error.html"
		return
	}
	appInfo,err := as.FindIpWithName(info.AppName)
	if err != nil {
		logs.Warn(err)
		l.Data["Error"] = err.Error()
		l.TplName = "error/error.html"
		return
	}
	keyFormat := tools.CheckSuffix(appInfo.DevelopPath)
	// 判断部署是否已存在
	exist := as.FindAppInfoWithPath(keyFormat)
	ips := strings.Split(appInfo.Ip,",")
	for _, ip := range ips {
		etcdKey := keyFormat + ip + "/"
		fmt.Println(etcdKey)
		// 写入etcd
		err = ls.InsertInfoToEtcd(exist,etcdKey,info)
		if err != nil {
			logs.Warn("Set log conf to etcd failed, err:%v", err)
			continue
		}
		// 记录
		err = eks.InsertEtcdKey(etcdKey)
		if err != nil {
			logs.Warn("insert into etcdKey db failed, err:%v", err)
			continue
		}
	}
	l.Redirect("/log/list",302)
	return
}

// etcd key 列表
func (l *LogController) LogKeys () {
	l.Layout = "layout/layout.html"
	list,err := eks.EtcdKeyList()
	if err != nil {
		logs.Warn(err)
		l.Data["Error"] = err.Error()
		l.TplName = "error/error.html"
		return
	}
	l.Data["KeyList"] = list
	l.TplName = "log/key.html"
	return
}