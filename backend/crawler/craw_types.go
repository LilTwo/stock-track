package crawler

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"regexp"
	"stock-tracker.com/types"
	"strconv"
	"time"
)

func WriteToDb(db *gorm.DB, ch chan *types.RawStockType) {
	db.AutoMigrate(&types.StockType{})
	count := 0
	for {
		select {
		case stock := <- ch:
			db.Create(&types.StockType{RawStockType: *stock, StockCode: stock.Code})
			count = 0
		default:
			count += 1
			fmt.Println("waiting", count)
			time.Sleep(time.Second)
			if count > 3 {
				return
			}
		}
	}
}

func ParseTokenizer(ch chan *types.RawStockType, tokenizer *html.Tokenizer, atStockMarket bool) {
	inTd := true
	nextIsISIN := false
	end := false
	codePattern := regexp.MustCompile("^([0-9]{4})[\u3000]+(.*)$")
	var t types.RawStockType
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.TextToken:
			token := tokenizer.Token()
			if inTd {
				if !nextIsISIN {
					matchResult := codePattern.FindStringSubmatch(token.Data)
					// all, code, name
					if len(matchResult) == 3 {
						code := matchResult[1]
						name := matchResult[2]
						convertedCode, _ := strconv.Atoi(code)
						t = types.RawStockType{Code: convertedCode, Name: name, AtStockMarket: atStockMarket}
						nextIsISIN = true
					}
				} else {
					t.ISN = token.Data
					ch <- &t
					nextIsISIN = false
				}
			}
		case html.StartTagToken:
			tag, _ := tokenizer.TagName()
			if string(tag) == "td" {
				inTd = true
			}
		case html.EndTagToken:
			tag, _ := tokenizer.TagName()
			if string(tag) == "td" {
				inTd = false
			}
		case html.ErrorToken:
			end = true
		}
		if end {
			break
		}
	}
}

func CrawStockTypes(db *gorm.DB) {
	ch := make(chan *types.RawStockType, 5)
	resp, _ := http.Get("https://isin.twse.com.tw/isin/C_public.jsp?strMode=2") //上市
	content, _ := ioutil.ReadAll(resp.Body)
	rawString, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), string(content))
	tokenizer := html.NewTokenizer(bytes.NewReader([]byte(rawString)))
	go ParseTokenizer(ch, tokenizer, true)

	resp, _ = http.Get("https://isin.twse.com.tw/isin/C_public.jsp?strMode=4") //上櫃
	content, _ = ioutil.ReadAll(resp.Body)
	rawString, _, _ = transform.String(traditionalchinese.Big5.NewDecoder(), string(content))
	tokenizer = html.NewTokenizer(bytes.NewReader([]byte(rawString)))
	go ParseTokenizer(ch, tokenizer, false)

	WriteToDb(db, ch)
}
