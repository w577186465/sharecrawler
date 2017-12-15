package modules

import (
	"crypto/md5"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/vsuper/spider/request"
	"spider/models"
	"strings"
	"time"
)

type Report struct {
	Pagesize  int
	Method    string
	Startpage int
	Thread    int
}

var (
	rpset Report
)

func (set Report) Getreport() {
	rpset = set
	// 获取信息
	json, err := reqreport(rpset.Pagesize, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	pagenum, _ := json.Get("pages").Int() // 分页数

	// 起始页
	startpage := rpset.Startpage
	if startpage == 0 {
		startpage = 1
	}

	ch := make(chan int, 10)
	defer close(ch)
	goon := true // 是否继续
	for i := startpage; i <= pagenum; i++ {
		if goon == false {
			fmt.Println("err")
			return
		}
		ch <- i
		page := i
		go func() {
			fmt.Println(page)
			goon = parsereport(page)
			<-ch
		}()

	}
}

func reqreport(pagesize, page int) (*simplejson.Json, error) {
	url := fmt.Sprintf(`http://datainterface.eastmoney.com//EM_DataCenter/js.aspx?type=SR&sty=GGSR&js={"data":[(x)],"pages":(pc),"update":"(ud)","count":(count)}&ps=%d&p=%d&mkt=0&stat=0&cmd=2&code=&rt=50154142`, pagesize, page)
	return request.NewJson(&request.Request{Url: url})
}

func parsereport(page int) bool {
	json, err := reqreport(rpset.Pagesize, page)
	if err != nil {
		return true
	}
	data := json.Get("data")
	arr, err := data.Array()
	if err != nil {
		return true
	}

	var reports []models.Report
	var reportdatas []models.ReportData
	for _, v := range arr {
		var report models.Report
		var reportdata models.ReportData
		item := v.(map[string]interface{})

		// 时间
		datetime := item["datetime"].(string)
		t, _ := time.ParseInLocation("2006-01-02T15:04:05", datetime, time.Local)
		day := t.Format("20060102")
		infoCode := item["infoCode"].(string)

		fullcode := item["secuFullCode"].(string)
		report.Code = fullcode[0:6]
		report.Name = item["secuName"].(string)
		report.Title = item["title"].(string)
		report.Author = item["author"].(string)
		report.Rate = item["rate"].(string)
		report.Change = item["change"].(string)
		report.Insname = item["insName"].(string)
		report.CreatedAt = t
		report.Hash = fmt.Sprintf("%x", md5.Sum([]byte(infoCode+day+report.Rate+report.Change)))

		// 判断报告是否存在
		var has int
		db.Model(&models.Report{}).Where("hash = ?", report.Hash).Count(&has)
		if has > 0 {
			if rpset.Method == "update" {
				reportsave(&reports, &reportdatas)
				fmt.Println(page)
				fmt.Println(report)
				fmt.Println("更新完成")
				return false
			} else {
				fmt.Println("该报告已存在")
				continue
			}
		}

		// 抓取content
		reportdata.Content = reportcontent(day, infoCode)

		reports = append(reports, report)
		reportdatas = append(reportdatas, reportdata)
	}
	reportsave(&reports, &reportdatas)
	return true
}

func reportcontent(day, infoCode string) string {
	var content string
	url := fmt.Sprintf("http://data.eastmoney.com/report/%s/%s.html", day, infoCode)
	document, err := request.NewDocument(&request.Request{Url: url, Char: "GBK"})
	if err != nil {
		return content
	}
	content, _ = document.Find(".newsContent").Html()
	return strings.TrimSpace(content)
}

func reportsave(reports *[]models.Report, reportdatas *[]models.ReportData) {
	fmt.Println(len(*reports))

	tx := db.Begin()
	for _, v := range *reports {
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
