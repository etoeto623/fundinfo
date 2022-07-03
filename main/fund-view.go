package main

import (
	"errors"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math"
	"neolong.me/fundinfo/common"
	"neolong.me/fundinfo/dbs"
	"net/http"
	"strconv"
	"time"
)

func handleQuery(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	code := params.Get("code")
	fromDate := params.Get("from") // yyyy-MM-dd
	toDate := params.Get("to") // yyyy-MM-dd
	if len(code) <= 0{
		fmt.Fprint(w, "param is illegal, please ensure parameter [code] is specified")
		return
	}

	datas, err := queryFund(code, fromDate, toDate)
	if nil != err {
		fmt.Fprint(w, "data query error:", err.Error())
		return
	}

	drawEchart(w, datas, code, "")
}

// 画echart图表
func drawEchart(w http.ResponseWriter, datas []common.FundWorth, code string, msg string) {
	var xAxis []string
	var lineData []opts.LineData
	min := 10.0
	max := 0.0

	for i := range datas {
		item := datas[i]
		xAxis = append(xAxis, item.InfoDate)
		lineData = append(lineData, opts.LineData{Value: item.UnitWorth})
		if item.UnitWorth < min {
			min = item.UnitWorth
		}
		if item.UnitWorth > max {
			max = item.UnitWorth
		}
	}

	if len(msg) <= 0 {
		msg = code;
	}

	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "基金净值趋势("+code+")",
			Subtitle: msg,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: roundFloat(min - 0.2, 1),
			Max: roundFloat(max + 0.2, 1)}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
			Trigger: "item",
			AxisPointer: &opts.AxisPointer{Type: "cross", Snap: true}},
		))

	// Put data into instance
	line.SetXAxis(xAxis).
		AddSeries(code, lineData).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func queryFund(code, fromDate, toDate string) ([]common.FundWorth, error) {
	// get db connect
	dbConn, err := dbs.OpenDefaultDB()
	if nil != err {
		return nil, err
	}
	defer dbConn.Close()

	sqlParams := make([]interface{}, 1)
	sqlStr := "select id, fund_code FundCode, info_date InfoDate, unit_worth UnitWorth, " +
		"total_worth TotalWorth from fund_worth where fund_code=?"
	sqlParams[0] = code
	if len(fromDate) > 0 {
		sqlStr += " and info_date >= ?"
		sqlParams = append(sqlParams, fromDate)
	}
	if len(toDate) > 0 {
		sqlStr += " and info_date <= ?"
		sqlParams = append(sqlParams, toDate)
	}
	sqlStr += " order by info_date asc"

	rows, err := dbConn.Query(sqlStr, sqlParams...)
	if nil != err {
		return nil, err
	}

	var datas []common.FundWorth
	for {
		if rows.Next() {
			item := common.FundWorth{}
			rows.Scan(&item.Id, &item.FundCode, &item.InfoDate, &item.UnitWorth, &item.TotalWorth)
			datas = append(datas, item)
		}else{
			break
		}
	}
	return datas, nil
}

//func showGain(w http.ResponseWriter, req *http.Request) {
//	calcGain(w, req)
//}

// 计算收益
func calcGain(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	code := params.Get("code")
	start := params.Get("start") // yyyy-MM-dd
	end := params.Get("end")
	amount := params.Get("amount") // 每次的金额
	if len(code) <= 0 || len(start) <= 0 || len(amount) <= 0 {
		fmt.Fprint(w, "illegal param, code/start/amount is required")
		return
	}
	amountVal, err := strconv.ParseInt(amount, 10, 64)
	if nil != err {
		fmt.Fprint(w, "amount illegal:", err.Error())
		return
	}

	// 计算endDate
	endDate := time.Now()
	if len(end) > 0 {
		endDate, err = time.Parse(common.TIME_PTN, end)
		if nil != err {
			fmt.Fprint(w, "end date illegal:", err.Error())
			return
		}
	}
	startDate, _ := time.Parse(common.TIME_PTN, start)

	// 查询所有的净值记录
	datas, err := queryFund(code, start, "")
	if nil != err {
		fmt.Fprint(w, "data query error:", err.Error())
		return
	}
	if len(datas) <= 0 {
		fmt.Fprint(w, "no fund worth data found")
		return
	}
	dataMap := listToMap(datas)

	endWorth, err := queryFundUnitWorth(endDate, dataMap);
	if nil != err {
		fmt.Fprint(w, "no worth data for end date")
		return
	}

	totalInvest := int64(0) // 总投入
	units := 0.0
	// 开始计算收益
	for {
		var fund *common.FundWorth
		if fund, err = queryFundUnitWorth(startDate, dataMap); nil != err {
			fmt.Fprint(w, "query fund error:", err.Error())
			return
		}
		totalInvest += amountVal
		units += roundFloat(float64(amountVal) / fund.UnitWorth, 4)
		startDate = startDate.AddDate(0, 1, 0)
		if startDate.After(endDate) {
			break
		}
	}
	final_worth := roundFloat(units * endWorth.UnitWorth, 2)
	msg := fmt.Sprintf("总投资: %d  总份额: %f  总市值: %f 盈亏：%f", totalInvest, units,
		final_worth, final_worth - float64(totalInvest))
	drawEchart(w, datas, code, msg)
}
// 查询当前投资日期的净值数据
func queryFundUnitWorth(date time.Time, fundMap map[string]common.FundWorth) (*common.FundWorth, error){
	now := time.Now()
	for {
		dateStr := date.Format(common.TIME_PTN)
		if worth, ok := fundMap[dateStr]; ok {
			return &worth, nil
		}
		// 往后一天
		date = date.AddDate(0, 0, 1)
		if date.After(now) {
			return nil, errors.New("no fund unit worth found")
		}
	}
}

// list转map
func listToMap(datas []common.FundWorth) map[string]common.FundWorth {
	result := make(map[string]common.FundWorth)
	for i := range datas {
		result[datas[i].InfoDate] = datas[i]
	}
	return result
}

func main() {
	// http://localhost:8888/view/?code=000962&from=2022-04-11
	http.HandleFunc("/view/", handleQuery)
	http.HandleFunc("/gain/", calcGain)
	fmt.Println("server listening at 8888")
	http.ListenAndServe(":8888", nil)
}