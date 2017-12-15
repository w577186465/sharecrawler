package modules

import (
	"crypto/md5"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/vsuper/spider/request"
	"spider/models"
	"strconv"
	"strings"
	"time"
)

type Hyreport struct {
	Pagesize  int
	Method    string
	Startpage int
	Thread    int
}

var (
	hyrpset  Hyreport
	pagesize = 50
	savenum  = 1000
	method   string
)

func (set Hyreport) HyReport() {
	hyrpset = set
	method = "all"

	json, err := getPage(hyrpset.Pagesize, 1)
	if err != nil {
		return
	}
	pages, _ := json.Get("pages").Int()

	ch := make(chan int, hyrpset.Thread)
	defer close(ch)
	goon := true
	for i := 1; i <= pages; i++ {
		if !goon {
			return
		}
		ch <- i
		page := i
		go func() {
			fmt.Println(page)
			goon = parsehyreport(page)
			<-ch
		}()
	}
}

func parsehyreport(page int) bool {
	fail := 0
	var reports []models.Hyreport
	var reportdatas []models.HyreportData

	data, err := getPage(page, pagesize)
	if err != nil {
		fail += pagesize
		fmt.Println(err)
		return true
	}

	arr, err := data.Get("data").Array() // 获取data
	if err != nil {
		fmt.Println(err)
		return true
	}

	// ch := make(chan int, 10)
	for _, v := range arr {
		item := v.(string)
		// ch <- 1

		var report models.Hyreport
		var reportdata models.HyreportData

		arr := strings.Split(item, ",")
		t, _ := time.ParseInLocation("2006/1/2 15:04:05", arr[1], time.Local) // 将时间转换为时间类型
		day := t.Format("20060102")                                           // 生成详情页地址时间

		// go func() {

		// 类型转换
		indid, _ := strconv.Atoi(arr[6])
		fluctuation, _ := strconv.ParseFloat(arr[11], 64)

		report.Pjchange = arr[0] // 评级变动
		report.Date = t
		report.Insname = arr[4] // 机构名称
		report.Indid = indid    // 行业id
		report.Pjtype = arr[7]  // 评级类型
		report.Expect = arr[8]  // 看好
		report.Title = arr[9]
		report.Indname = arr[10]         // 行业名称
		report.Fluctuation = fluctuation // 涨跌幅

		hash := fmt.Sprintf("%x", md5.Sum([]byte(arr[2]+report.Indname+report.Pjchange+report.Pjtype+report.Expect+day))) // 生成hash
		report.Hash = hash

		// 判断是否存在
		var has int
		db.Model(&models.Hyreport{}).Where("hash = ?", hash).Count(&has)

		// 如果存在 更新操作停止抓取 all操作跳过
		if has > 0 {
			if hyrpset.Method == "update" {
				fmt.Println("更新完成")
				if num > 0 {
					hyreportsave(&reports, &reportdatas)
				}
				return false
			} else if hyrpset.Method == "all" {
				fmt.Println("已存在")
				// <-ch
				continue
			}

		}

		// content
		contenturl := fmt.Sprintf("http://data.eastmoney.com/report/%s/hy,%s.html", day, arr[2]) // 详情页地址
		content := getcontent(contenturl)                                                        // 获取研报详情
		reportdata.Content = content
		reports = append(reports, report)
		reportdatas = append(reportdatas, reportdata)

		// 抓取数量达到savenum 保存到数据库并清空抓取数据

		// <-ch
		// }()
	}

	hyreportsave(&reports, &reportdatas)
	return true
}

// func GetOld() {
// 	var hyreport models.Hyreport
// 	db.Order("date ASC").Select("date, hash").First(&hyreport) // 获取最后更新

// 	// 获取第一页
// 	db.Model(&hyreport).Count(&count) //获取数据量
// 	first := count / pagesize
// 	if count%pagesize != 0 {
// 		first++
// 	}

// 	getdata(first)
// }

func getPage(page, pagesize int) (*simplejson.Json, error) {
	starturl := fmt.Sprintf(`http://datainterface.eastmoney.com//EM_DataCenter/js.aspx?type=SR&sty=HYSR&mkt=0&stat=0&cmd=4&code=&sc=&ps=%d&p=%d&js={"data":[(x)],"pages":(pc),"update":"(ud)","count":(count)}`, pagesize, page)
	return request.NewJson(&request.Request{Url: starturl})
}

// func get() {
// 	count := getcount() // 获取信息数量

// 	// 获取分页数
// 	pages := int(count / 50)
// 	if count/50 > int(count/50) {
// 		pages += 1
// 	}

// 	ch := make(chan int, 10)

// 	for i := 1; i <= pages; i++ {
// 		ch <- 1
// 		fmt.Printf("正在获取第%d页", i)
// 		go func() {
// 			d, _ := getdata(i)
// 			savedata := getsavedata(d)
// 			save(savedata)
// 			if over {
// 				fmt.Println("抓取完成")
// 				return
// 			}
// 			<-ch
// 		}()
// 	}
// }

func hyreportsave(reports *[]models.Hyreport, reportdatas *[]models.HyreportData) {
	d := *reports
	num := len(d)
	if num == 0 {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tx := db.Begin()
	for _, v := range d {
		if err := tx.Create(&v).Error; err != nil {
			fmt.Println("保存失败")
			fmt.Println(err)
			tx.Rollback()
			return
		}
	}

	for _, v := range *reportdatas {
		if err := tx.Create(&v).Error; err != nil {
			fmt.Println("保存失败")
			fmt.Println(err)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	fmt.Println("保存成功")
}

func getcontent(url string) string {
	var content string
	req := request.Request{
		Url:  url,
		Char: "gbk",
	}
	document, err := request.NewDocument(&req)
	if err != nil {
		fmt.Println(err)
		return content
	}
	content, err = document.Find(".newsContent").Html()
	if err != nil {
		fmt.Println(err)
		return content
	}
	return strings.TrimSpace(content)
}
