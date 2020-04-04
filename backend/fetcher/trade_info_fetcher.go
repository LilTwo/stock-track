package fetcher

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"stock-tracker.com/types"
	"time"
)

type TradeInfoFetcher struct {
	Db *gorm.DB
	tradeInfos []types.TradeInfo
}

func (fetcher *TradeInfoFetcher) QueryByCode(code int) *TradeInfoFetcher {
	fetcher.Db = fetcher.Db.Where("stock_code = ?", code)
	return fetcher
}

func (fetcher *TradeInfoFetcher) QueryByStartDate(startDate string) *TradeInfoFetcher {
	layout := "2006-01-02"
	date, err := time.Parse(layout, startDate)
	if err != nil {
		fmt.Println("wrong date format")
		panic(err)
	}
	fetcher.Db = fetcher.Db.Where("date >= ?", date)
	return fetcher
}

func (fetcher *TradeInfoFetcher) QueryByEndDate(endDate string) *TradeInfoFetcher {
	layout := "2006-01-02"
	date, err := time.Parse(layout, endDate)
	if err != nil {
		fmt.Println("wrong date format")
		panic(err)
	}
	fetcher.Db = fetcher.Db.Where("date <= ?", date)
	return fetcher
}

func (fetcher *TradeInfoFetcher) Get() []types.TradeInfo {
	fetcher.Db = fetcher.Db.Order("date asc").Find(&fetcher.tradeInfos)
	return fetcher.tradeInfos
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d", date.Year(), int(date.Month()), date.Day())
}

func MapToMinPriceAndDate(tradeInfos []types.TradeInfo) []struct{ Date string; Data float64 } {
	var result []struct{ Date string; Data float64 }
	for _, tradeInfo := range tradeInfos {
		result = append(result, struct {
			Date     string
			Data float64
		}{Date: formatDate(tradeInfo.Date), Data: tradeInfo.MinPrice})
	}
	return result
}
