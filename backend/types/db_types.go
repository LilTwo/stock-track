package types

import (
	"github.com/jinzhu/gorm"
)

type StockType struct {
	gorm.Model
	RawStockType
	StockCode int `gorm:"UNIQUE_INDEX;NOT NULL"` // user Code will make AutoMigrate failed due to name collision with embedded field
}

type TradeInfo struct {
	gorm.Model
	RawTradeInfo
	StockType StockType `gorm:"foreignkey:StockCode;association_foreignkey:StockCode"`
	CodeDate string `gorm:"UNIQUE_INDEX;NOT NULL"`
}