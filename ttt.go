package main

import (
	"spider/spider-modules"
)

func main() {
	report := modules.Report{
		Pagesize:  50,
		Method:    "all",
		Startpage: 1,
		Thread:    10,
	}
	report.Getreport()
}
