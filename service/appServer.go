package service

import (
	"beegoDemo/dial"
	"beegoDemo/models"
)

type AppInfoServer interface {
	// 获取所有信息
	AppInfoList() (list []models.AppInfo,err error)
	// 插入一条信息
	InsertAppInfo(info models.AppInfo) (id int64,err error)
	// 根据名称获取信息
	FindIpWithName(name string) (info models.AppInfo,err error)
	// 根据部署路径查询
	FindAppInfoWithPath(path string) (exist bool)
}

type appInfoServer struct {}

func NewAppInfoServer() AppInfoServer {
	return &appInfoServer{}
}
// 获取所有信息
func (as *appInfoServer) AppInfoList() (list []models.AppInfo,err error) {
	qs := dial.O.QueryTable("app_info")
	 _,err = qs.All(&list)
	return
}
// 插入一条信息
func (as *appInfoServer) InsertAppInfo(info models.AppInfo) (id int64,err error) {
	var i models.AppInfo
	i.AppName = info.AppName
	i.AppType = info.AppType
	i.DevelopPath = info.DevelopPath
	i.Ip = info.Ip
	id,err = dial.O.Insert(&i)
	return
}
// 根据名称获取信息
func (as *appInfoServer) FindIpWithName(name string) (info models.AppInfo,err error) {
	info.AppName = name
	err = dial.O.Read(&info,"app_name")
	return
}
// 根据部署路径查询
func (as *appInfoServer) FindAppInfoWithPath(path string) (exist bool) {
	qs := dial.O.QueryTable("app_info")
	exist = qs.Filter("develop_path",path).Exist()
	return
}