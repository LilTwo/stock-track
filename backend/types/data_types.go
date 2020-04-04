package types

import (
	"time"
)

type RawStockType struct {
	Code int
	Name string
	ISN string
	AtStockMarket bool //上市
}

type RawTradeInfo struct {
	Date time.Time
	NumStock int
	TurnOver uint64
	OpeningPrice float64
	ClosingPrice float64
	MaxPrice float64
	MinPrice float64
	NumTransaction int
	StockCode int
}