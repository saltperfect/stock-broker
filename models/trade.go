package models

import "time"

// Trade represents a stock trade
type Trade struct {
	TradeID        int
	BuyerOrderID   int
	SellerOrderID  int
	StockSymbol    string
	Quantity       int
	Price          float64
	TradeTimestamp time.Time
}
