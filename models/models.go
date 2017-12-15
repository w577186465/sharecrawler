package models

import (
	"time"
)

type Industry struct {
	Id      int
	Name    string
	Dfcfval int
	Fw      string
}

type Allshare struct {
	Id            int
	Symbol        string `gorm:"type:char(8);unique"`
	Code          string `gorm:"type:char(6);unique"`
	Name          string `gorm:"type:varchar(100)"`
	Trade         string
	Pricechange   string
	Changepercent string
	Buy           string
	Sell          string
	Settlement    string
	Open          string
	High          string
	Low           string
	Volume        string
	Amount        string
	Ticktime      string
	Per           float64
	Per_d         float64
	Nta           string
	Pb            float64
	Mktcap        float64
	Nmc           float64
	Turnoverratio float64
	Dfcfval       string
	Favor         string
	Guba          string
	UpdatedAt     time.Time
}

type Report struct {
	Id        int
	Code      string
	Name      string
	Title     string
	Hash      string
	Author    string
	Rate      string
	Change    string
	CreatedAt time.Time
	Insname   string
}

type ReportData struct {
	Id      int
	Content string
}

type Hyreport struct {
	Id          int
	Pjchange    string
	Date        time.Time
	Insname     string
	Indid       int
	Pjtype      string
	Expect      string
	Title       string
	Indname     string
	Fluctuation float64
	Hash        string
	UpdatedAt   time.Time
}

type HyreportData struct {
	Id      int
	Content string
}
