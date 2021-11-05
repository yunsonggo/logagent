package models

import "time"

type AppInfo struct {
	Id       int    `orm:"auto;pk" form:"id"`
	AppName     string `orm:"app_name" form:"app_name"`
	AppType     string `orm:"app_type" form:"app_type"`
	CreateTime  time.Time `orm:"create_time;auto_now_add;type(datetime)" form:"create_time"`
	DevelopPath string `orm:"develop_path" form:"develop_path"`
	Ip         string `orm:"ip" form:"ip"`
}

func (ai *AppInfo) TableName() string {
	return "app_info"
}