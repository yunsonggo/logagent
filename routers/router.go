package routers

import (
	"beegoDemo/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{},"get:Index")
	beego.Router("/index", &controllers.MainController{},"get:Index")
	beego.Router("/list", &controllers.MainController{}, "get:Index")

	beego.Router("/app/apply", &controllers.MainController{}, "get:AppApply")
	beego.Router("/app/apply",&controllers.MainController{},"post:AppCreate")

	beego.Router("/log/list",&controllers.LogController{},"get:LogList")
	beego.Router("/log/apply",&controllers.LogController{},"get:LogApply")
	beego.Router("/log/apply",&controllers.LogController{},"post:LogCreate")
	beego.Router("/log/keys",&controllers.LogController{},"get:LogKeys")

}
