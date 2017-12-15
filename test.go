package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/libs/database"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"spider/models"
	"spider/spider-modules"
	"time"
)

type Testt struct {
	Insname     string
	Indid       int
	Pjtype      string
	Expect      string
	Title       string
	Indname     string
	Fluctuation float64
	Content     string
	Md5         string
	State       int
}

var db, _ = database.Open()

func init() {
	beego.SetLogger("file", `{"filename":"logs/err.log"}`)
	db, _ = database.Open()
}

func main() {
	hyreportTest()
}

func selectTest() {
	var code models.Allshare
	db.Select("code").Where("code = ?", "0000011").First(&code)
	if code.Code == "" {
		fmt.Println("have not")
	}
}

func test() {
	var shares []models.Allshare
	db.Select("code").Where("dfcfval = ''").Find(&shares)
	fmt.Println(shares)
}

// func gplistTest() {
// 	modules.AllShare()
// 	modules.Industry(false)
// 	db.Close()
// }

func hyreportTest() {
	var set modules.Hyreport
	set.Pagesize = 50
	set.Method = "update"
	set.Startpage = 1
	set.Thread = 10

	set.HyReport()
	return
	var hyreport models.Hyreport
	var count int
	db.Model(&hyreport).Count(&count)
	fmt.Println(count)
	return
	db.Order("date ASC").Select("date").First(&hyreport)

	t, _ := time.Parse("2006/1/2 15:04:05", "2017/9/27")
	fmt.Println(t.Before(hyreport.Date))
}
