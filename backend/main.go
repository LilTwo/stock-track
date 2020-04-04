package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"os"
	"stock-tracker.com/fetcher"
	"strconv"
)

func QueryStockType(w http.ResponseWriter, req *http.Request, db *gorm.DB) {
	if req.Method == "GET" {
		search := req.URL.Query()["search"][0]
		stockTypeFetcher := fetcher.StockTypeFetcher{
			Db: db,
		}
		stockTypeFetcher.QueryByCode(search).QueryName(search)
		result := stockTypeFetcher.Get(10)
		jsonResult, _ := json.Marshal(result)
		w.Write(jsonResult)
	}
}

func QueryTradeInfo(w http.ResponseWriter, req *http.Request, db *gorm.DB) {
	if req.Method == "GET" {
		//field := req.URL.Query()["field"][0]
		startDateQuery := req.URL.Query()["startDate"]
		endDateQuery := req.URL.Query()["endDate"]
		codeQuery := req.URL.Query()["code"]
		tradeInfoFetcher := fetcher.TradeInfoFetcher{
			Db: db,
		}
		if len(codeQuery) != 1 || codeQuery[0] == ""{
			w.WriteHeader(403)
			return
		}
		code := codeQuery[0]
		intCode, _ := strconv.Atoi(code)
		tradeInfoFetcher.QueryByCode(intCode)
		if len(startDateQuery) == 1 && startDateQuery[0] != "" {
			tradeInfoFetcher.QueryByStartDate(startDateQuery[0])
		}

		if len(endDateQuery) == 1 && endDateQuery[0] != "" {
			tradeInfoFetcher.QueryByEndDate(endDateQuery[0])
		}

		tradeInfos := tradeInfoFetcher.Get()
		result := fetcher.MapToMinPriceAndDate(tradeInfos)
		jsonResult, _ := json.Marshal(result)
		w.Write(jsonResult)
	}
}

func main() {
	rootPassword := os.Getenv("ROOTPASS")
	db, err := gorm.Open("mysql", "root:"+rootPassword+"@(localhost)/stock?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	autoCompleteHandler := func(w http.ResponseWriter, req *http.Request) { QueryStockType(w, req, db)}
	http.HandleFunc("/api/auto-complete", autoCompleteHandler)
	tradeInfoHandler := func(w http.ResponseWriter, req *http.Request) { QueryTradeInfo(w, req, db)}
	http.HandleFunc("/api/trade-info", tradeInfoHandler)

	/*
	run these functions to crawl data
	crawler.CrawStockTypes(db)
	crawler.CrawTransactions(db)
	 */
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
