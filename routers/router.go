package routers

import (
	"github.com/astaxie/beego"
	"spider/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/task/allshare", &controllers.TaskController{}, "get:Allshare")
}
