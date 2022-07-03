package common

const TIME_PTN = "2006-01-02"

type CrawTask interface {
	// 进行爬取
	Craw(url string);
}

type FundWorth struct {
	Id int64 `json:"id"`
	FundCode string `json:"fundCode"`
	InfoDate string `json:"infoDate"`
	UnitWorth float64 `json:"unitWorth"`
	TotalWorth float64 `json:"totalWorth"`
	Source string `json:"source"`
}