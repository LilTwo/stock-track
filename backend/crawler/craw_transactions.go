package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"regexp"
	"stock-tracker.com/types"
	"strconv"
	"strings"
	"sync"
	"time"
)

func writeToDb(db *gorm.DB, ch chan []types.RawTradeInfo, wg *sync.WaitGroup, resetCh chan bool) {
	db.AutoMigrate(&types.TradeInfo{})
	count := 0

	for {
		select {
		case fetchResult := <- ch:
			for _, rawTradeInfo := range fetchResult {
				tradeInfo := types.TradeInfo{RawTradeInfo: rawTradeInfo}
				tradeInfo.CodeDate = getUniqueCodeDate(tradeInfo.Date, tradeInfo.StockCode)
				db.Create(&tradeInfo)
				count = 0
			}
			if len(fetchResult) > 0 {
				fmt.Println(fetchResult[0].StockCode, "Complete")
			}
		case reset := <- resetCh:
			count = 0
			fmt.Println("reset counter", reset)
		default:
			count += 1
			fmt.Println("waiting", count)
			time.Sleep(1*time.Second)
			if count > 100 {
				wg.Done()
				return
			}
		}
	}


}

func parseRawDate(rawDate string) time.Time {
	datePattern := regexp.MustCompile("^([0-9]{3}).*([0-9]{2}).*([0-9]{2})")
	matchDate := datePattern.FindStringSubmatch(rawDate)
	year, month, day := matchDate[1], matchDate[2], matchDate[3]
	yearNum, err := strconv.Atoi(year)
	if err != nil {
		panic(err)
	}
	monthNum, err := strconv.Atoi(month)
	if err != nil {
		panic(err)
	}
	dayNum, err := strconv.Atoi(day)
	if err != nil {
		panic(err)
	}
	date := time.Date(yearNum + 1911, time.Month(monthNum), dayNum, 0, 0, 0, 0, time.Local)
	return date
}

func getUniqueCodeDate(date time.Time, code int) string {
	// return 2019March05-1102, for identifying a unique TradeInfo
	return fmt.Sprintf("%d%s%d-%d", date.Year(), date.Month().String(), date.Day(), code)
}

func convertToTradeInfo(data []string) *types.RawTradeInfo {
	var tradeInfo types.RawTradeInfo
	rawDate := data[0] // 109/02/02
	date := parseRawDate(rawDate)
	tradeInfo.Date = date
	tradeInfo.NumStock, _ = strconv.Atoi(strings.ReplaceAll(data[1], ",", ""))
	tradeInfo.TurnOver, _ = strconv.ParseUint(strings.ReplaceAll(data[2], ",", ""), 10, 64)
	tradeInfo.MaxPrice, _ = strconv.ParseFloat(data[3], 64)
	tradeInfo.MinPrice, _ = strconv.ParseFloat(data[4], 64)
	tradeInfo.OpeningPrice, _ = strconv.ParseFloat(data[5], 64)
	tradeInfo.ClosingPrice, _ = strconv.ParseFloat(data[6], 64)
	tradeInfo.NumTransaction, _ = strconv.Atoi(data[8])
	return &tradeInfo
}

func getJson(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResult map[string]interface{}
	err = json.Unmarshal(content, &jsonResult)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}

func fetchTradeInfo(code int, year int, month int, ch chan []types.RawTradeInfo) {
	// this api always returns the data of the entire month
	url := fmt.Sprintf("http://www.tpex.org.tw/web/stock/aftertrading/daily_trading_info/st43_result.php?stkno=%d&d=%d/%d", code, year - 1911, month)
	jsonResult, err := getJson(url)
	if err != nil {
		fmt.Println("fetch failed", err)
		return
	}
	monthlyData := jsonResult["aaData"].([]interface{})
	var result []types.RawTradeInfo
	for _, dailyData :=  range monthlyData {
		var data []string
		for _, rawData := range dailyData.([]interface{}) {
			data = append(data, rawData.(string))
		}
		tradeInfo := *convertToTradeInfo(data)
		tradeInfo.StockCode = code
		result = append(result, tradeInfo)
	}
	ch <- result
}

func fetchTradeInfoAtMarket(code int, year int, month int, ch chan []types.RawTradeInfo) {
	// this api always returns the data of the entire month
	url := fmt.Sprintf("https://www.twse.com.tw/exchangeReport/STOCK_DAY/?stockNo=%d&date=%02d%02d01", code, year, month)
	jsonResult, err := getJson(url)
	if err != nil {
		fmt.Println("fetch failed", err)
		return
	}
	monthlyData := jsonResult["data"].([]interface{})
	var result []types.RawTradeInfo
	for _, dailyData :=  range monthlyData {
		var data []string
		for _, rawData := range dailyData.([]interface{}) {
			data = append(data, rawData.(string))
		}
		tradeInfo := *convertToTradeInfo(data)
		tradeInfo.StockCode = code
		result = append(result, tradeInfo)
	}
	ch <- result
}

func CrawTransactions(db *gorm.DB) {
	var allStockTypes []types.StockType
	db.Find(&allStockTypes)
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan []types.RawTradeInfo, 5)
	resetCh := make(chan bool, 1)
	go writeToDb(db, ch, &wg, resetCh)
	for _, stockType := range allStockTypes {
		var info types.TradeInfo
		db.Where("stock_code = ?", stockType.Code).First(&info)
		//if info.StockCode != 0 {
		//	select {
		//	case resetCh <- true:
		//		fmt.Println("write to reset")
		//	default:
		//		//fmt.Printf("%d exists, pass\n", stockType.Code)
		//	}
		//	continue
		//}
		if stockType.AtStockMarket {
			go fetchTradeInfoAtMarket(stockType.Code, 2020, 2, ch)
			time.Sleep(4000*time.Millisecond)
		} else {
			go fetchTradeInfo(stockType.Code, 2020, 2, ch)
			time.Sleep(4000*time.Millisecond)
		}
	}
	fmt.Println("waiting")
	wg.Wait()
}

/*
上市
{
   "stat":"OK",
   "date":"20200330",
   "title":"109年03月 1102 亞泥             各日成交資訊",
   "fields":[
      "日期",
      "成交股數",
      "成交金額",
      "開盤價",
      "最高價",
      "最低價",
      "收盤價",
      "漲跌價差",
      "成交筆數"
   ],
   "data":[
      [
         "109/03/02",
         "15,205,546",
         "665,152,354",
         "43.00",
         "44.45",
         "43.00",
         "43.90",
         "-0.45",
         "6,515"
      ]
   ],
   "notes":[
      "符號說明:+/-/X表示漲/跌/不比價",
      "當日統計資訊含一般、零股、盤後定價、鉅額交易，不含拍賣、標購。",
      "ETF證券代號第六碼為K、M、S、C者，表示該ETF以外幣交易。"
   ]
}

上櫃
{
   "stkNo":"5483",
   "stkName":"\u4e2d\u7f8e\u6676",
   "showListPriceNote":false,
   "showListPriceLink":false,
   "reportDate":"109\/04",
   "iTotalRecords":1,
   "aaData":[
      [
         "109\/04\/01",
         "10,712",
         "825,370",
         "77.60",
         "77.90",
         "76.50",
         "77.20",
         "-1.00",
         "7,010"
      ]
   ]
}
 */