package modules

import (
	"github.com/astaxie/beego/libs/database"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func init() {
	var err error

	// 连接数据库
	db, err = database.Open()
	if err != nil {

	}

}
