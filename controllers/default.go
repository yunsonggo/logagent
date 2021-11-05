package controllers

import (
	"beegoDemo/models"
	"beegoDemo/service"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

var as = service.NewAppInfoServer()

// 首页 项目列表数据渲染
func (c *MainController) Index() {
	list,err := as.AppInfoList()
	if err != nil {
		logs.Warn(err)
		c.Data["Error"] = err.Error()
		c.Layout = "layout/layout.html"
		c.TplName = "error/error.html"
		return
	}
	if len(list) == 0 {
		info := models.AppInfo{
			Id:          0,
			AppName:     "none",
			AppType:     "none",
			DevelopPath: "none",
			Ip:          "none",
		}
		list = append(list,info)
	}
	c.Data["infoList"] = list
	c.Layout = "layout/layout.html"
	c.TplName = "index.html"
}

func (c *MainController) AppApply() {
	c.Layout = "layout/layout.html"
	c.TplName = "app/apply.html"
}

//
func (c *MainController) AppCreate() {
	info := models.AppInfo{}
	c.Layout = "layout/layout.html"
	err := c.ParseForm(&info)
	if err != nil || len(info.AppName) == 0 || len(info.AppType) == 0 || len(info.DevelopPath) == 0{
		logs.Warn(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error/error.html"
		return
	}
	_,_ = as.InsertAppInfo(info)

	c.Redirect("/",302)
	return
}