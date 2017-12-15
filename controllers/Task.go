package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type TaskController struct {
	beego.Controller
}

func (c *TaskController) Get() {
	c.TplName = "task.tpl"
}

func (c *TaskController) Allshare() {

}
