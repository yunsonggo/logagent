package models

import "time"

type LogInfo struct {
	Id	int	`orm:"auto;pk" form:"id"`
	AppId      int    `orm:"app_id" form:"app_id"`
	AppName    string `orm:"app_name" form:"app_name"`
	LogId      int    `orm:"log_id" form:"log_id"`
	Status     int    `orm:"status" form:"status"`
	CreateTime time.Time `orm:"create_time;auto_now_add;type(datetime)" form:"create_time"`
	LogPath    string `orm:"log_path" form:"log_path"`
	Topic      string `orm:"topic" form:"topic"`
}

func (li *LogInfo) TableName() string {
	return "log_info"
}