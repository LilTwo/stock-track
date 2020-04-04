package fetcher

import (
	"github.com/jinzhu/gorm"
	"stock-tracker.com/types"
)

type StockTypeFetcher struct {
	Db *gorm.DB
	stockTypes []types.StockType
}

func (fetcher *StockTypeFetcher) QueryByCode(code string) *StockTypeFetcher {
	fetcher.Db = fetcher.Db.Or("code LIKE ?", "%" + code +"%")
	return fetcher
}

func (fetcher *StockTypeFetcher) QueryName(name string) *StockTypeFetcher {
	fetcher.Db = fetcher.Db.Or("name LIKE ?", "%" + name +"%")
	return fetcher
}

func (fetcher *StockTypeFetcher) Get(limit int) []types.StockType {
	fetcher.Db.Limit(limit).Find(&fetcher.stockTypes)
	return fetcher.stockTypes
}
