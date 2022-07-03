package unit_worth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"neolong.me/fundinfo/common"
	"neolong.me/fundinfo/dbs"
	"net/http"
	"os"
	"time"
)

const API = "https://stock.finance.sina.com.cn/fundInfo/api/openapi.php/CaihuiFundInfoService.getNav?symbol=%s&datefrom=%s&dateto=%s&page=%d"

type UnitWorthCrawler struct {
}

type UWResp struct {
	Result struct {
		Status struct {
			Code int `json:"code"`
		} `json:"status"`
		Data struct {
			Data []FundDetail `json:"data"`
			TotalNum string `json:"total_num"`
		} `json:"data"`
	} `json:"result"`
}
type FundDetail struct {
	Fbrq string `json:"fbrq"`
	Jjjz string `json:"jjjz"`
	Ljjz string `json:"ljjz"`
}

func (crawler UnitWorthCrawler) Craw(url string) *UWResp {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if nil != err {
		fmt.Println("craw error:", err.Error())
		return nil
	}

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("resp read error:", err.Error())
		return nil
	}
	respJson := UWResp{}
	err = json.Unmarshal(dataBytes, &respJson)
	if nil != err {
		fmt.Println("resp format error:", err.Error())
		return nil
	}

	return &respJson
}

// 爬取某一支基金的历史净值数据
func CrawFundAllWorth(fundCode string) {
	pageNo := 1
	dbConn := openDB()
	defer dbConn.Close()

	// 查询数据库中最新的数据
	rows, err := dbConn.Query("select max(info_date) md from fund_worth where fund_code=?", fundCode)
	if nil != err && err!=sql.ErrNoRows{
		fmt.Println("db query error:", err.Error())
		return
	}
	var dateFrom string
	var dateTo string
	if nil != rows && rows.Next() {
		rows.Scan(&dateFrom)
	}
	if len(dateFrom) > 0 {
		dateTo = time.Now().Format(common.TIME_PTN)
	}

	for  {
		url := fmt.Sprintf(API, fundCode, dateFrom, dateTo, pageNo)
		fmt.Println("---- prepare to craw", fundCode, "in page", pageNo)
		resp := UnitWorthCrawler{}.Craw(url)
		if isFinished(resp) {
			fmt.Println("++++", fundCode, "craw finished, total page is", pageNo-1)
			return
		}

		items := resp.Result.Data.Data
		for i := range items {
			saveFund(&items[i], dbConn, fundCode)
		}
		pageNo++
	}
}

func openDB() *sql.DB {
	dbConn, err := dbs.OpenDefaultDB()
	if nil != err {
		fmt.Println("dbs connect fail", err.Error())
		os.Exit(1)
	}
	return dbConn
}

func saveFund(item *FundDetail, db *sql.DB, fundCode string) {
	sqlStr := "insert into fund_worth(fund_code, info_date, unit_worth, total_worth, source) value (?,?,?,?,?)"
	_, err := db.Exec(sqlStr, fundCode, item.Fbrq, item.Jjjz, item.Ljjz, "sina")
	if nil != err {
		fmt.Println("fund detail save fail:", err.Error(), "【",fundCode,item.Fbrq,item.Jjjz,item.Ljjz,"】")
	}
}

func isFinished(resp *UWResp) bool {
	return nil == resp || resp.Result.Status.Code != 0 ||
		len(resp.Result.Data.Data) <= 0
}