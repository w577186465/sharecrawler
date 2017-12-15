package modules

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/libs/database"
	"github.com/bitly/go-simplejson"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vsuper/spider/request"
	"spider/models"
	"strings"
)

var (
	url string
	num = 80 // 每页数量
)

func init() {
	var err error
	beego.SetLogger("file", `{"filename":"logs/gplist_err.log"}`)
	db, err = database.Open()
	if err != nil {
		beego.Error(err)
	}

}

func AllShare() {
	url = `http://money.finance.sina.com.cn/d/api/openapi_proxy.php/?__s=[["hq","%s","",0,%d,%d]]`
	get("hs_a")
}

func get(t string) {
	count := getcount(t) // 获取数据量

	// 获取分页数
	pagen := count / num
	if pagen*num < count {
		pagen++
	}

	var data []models.Allshare // 声明保存数据

	// 获取全部股票
	for i := 1; i <= pagen; i++ {
		getitems(i, t, &data)
	}
	savedata(&data)
}

// 获取数据量
func getcount(t string) int {
	json, ok := getdata(1, t, 1)
	if !ok {
		panic("获取分页失败")
	}
	count, err := json.Get("count").Int() // 数量
	if err != nil {
		beego.Error(err)
		panic("获取分页失败")
	}
	return count
}

func savedata(data *[]models.Allshare) {
	tx := db.Begin()
	var allshare models.Allshare // 模型
	for _, v := range *data {
		var get models.Allshare
		db.Select("code").Where("code = ?", v.Code).First(&get) // 获取是否存在
		var err error
		if get.Code == "" {
			err = tx.Create(&v).Error // 不存在新增
		} else {
			err = tx.Model(allshare).Where("code = ?", v.Code).Updates(v).Error // 存在更新
		}

		if err != nil {
			fmt.Println("error:", err)
			tx.Rollback()
			break
		}
	}

	tx.Commit()
}

// joson 打开
func getdata(page int, t string, num int) (*simplejson.Json, bool) {
	url := fmt.Sprintf(url, t, page, num)

	req := &request.Request{
		Url: url,
	}

	json, err := request.NewJson(req)
	if err != nil {
		beego.Error(err)
		return nil, false
	}
	j0 := json.GetIndex(0)
	return j0, true
}

func getitems(page int, t string, allshare *[]models.Allshare) {
	j, ok := getdata(page, t, num)
	if !ok {
		return
	}

	items := j.Get("items")

	index := 0
	for {
		item := items.GetIndex(index)
		if _, err := item.Array(); err != nil {
			break
		}

		var ritem models.Allshare
		ritem.Symbol = item.GetIndex(0).MustString()
		ritem.Code = item.GetIndex(1).MustString()
		ritem.Name = item.GetIndex(2).MustString()
		ritem.Trade = item.GetIndex(3).MustString()
		ritem.Pricechange = item.GetIndex(4).MustString()
		ritem.Changepercent = item.GetIndex(5).MustString()
		ritem.Buy = item.GetIndex(6).MustString()
		ritem.Sell = item.GetIndex(7).MustString()
		ritem.Settlement = item.GetIndex(8).MustString()
		ritem.Open = item.GetIndex(9).MustString()
		ritem.High = item.GetIndex(10).MustString()
		ritem.Low = item.GetIndex(11).MustString()
		ritem.Volume = item.GetIndex(12).MustString()
		ritem.Amount = item.GetIndex(13).MustString()
		ritem.Ticktime = item.GetIndex(14).MustString()
		ritem.Per = item.GetIndex(15).MustFloat64()
		ritem.Per_d = item.GetIndex(16).MustFloat64()
		ritem.Nta = item.GetIndex(17).MustString()
		ritem.Pb = item.GetIndex(18).MustFloat64()
		ritem.Mktcap = item.GetIndex(19).MustFloat64()
		ritem.Nmc = item.GetIndex(20).MustFloat64()
		ritem.Turnoverratio = item.GetIndex(21).MustFloat64()
		ritem.Favor = item.GetIndex(22).MustString()
		ritem.Guba = item.GetIndex(23).MustString()

		*allshare = append(*allshare, ritem)
		index++
	}
}

func Industry(all bool) {
	shares := getshares(all) // 获取股票

	// panic 处理
	if err := recover(); err != nil {
		fmt.Println(err)
	}

	var data []models.Allshare
	for _, v := range *shares {
		code := v.Code
		// 股票市场
		arr := []rune(code)
		var n string
		if string(arr[0]) == "6" {
			n = "1"
		} else {
			n = "2"
		}

		// 打开页面
		href := fmt.Sprintf(`http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?type=CT&cmd=E.%s%s&sty=DCRRBKCPALTB&st=z&sr=-1&p=&ps=5&js=[(x)]&token=7bc05d0d4c3c22ef9fca8c2a912d779c`, code, n)

		req := &request.Request{Url: href}
		res, err := request.NewJson(req)
		if err != nil {
			continue
		}

		// 获取行业id
		info, err := res.GetIndex(0).String()
		if err != nil {
			fmt.Println(err)
			return
		}
		infos := strings.Split(info, ",")
		dfcfval := strings.Replace(infos[1], "BK", "", 1)
		data = append(data, models.Allshare{Code: code, Dfcfval: dfcfval})
	}

	savehy(&data) // 保存
}

func getshares(all bool) *[]models.Allshare {
	var shares []models.Allshare
	if all {
		db.Select("code").Find(&shares)
	} else {
		db.Select("code").Where("dfcfval = ''").Find(&shares)
	}

	return &shares
}

// 保存
func savehy(data *[]models.Allshare) {
	tx := db.Begin()
	for _, v := range *data {
		if err := tx.Model(&v).Where("code = ?", v.Code).Update("Dfcfval", v.Dfcfval); err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
}
