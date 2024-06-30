package supports

import (
	"design_patterns/stock_broker/models"
	"fmt"
	"sync"
	"time"
)

// TradeService is responsible for managing trades
type TradeService struct {
	trades map[int]*models.Trade
	mu     sync.Mutex
}

// NewTradeService creates a new TradeService instance
func NewTradeService() *TradeService {
	return &TradeService{
		trades: make(map[int]*models.Trade),
	}
}

func (s *TradeService) Show() {
	for range time.NewTicker(time.Second).C {
		for _, trade := range s.trades {
			fmt.Printf("trade: %v\n", trade)
		}
	}
}

// CreateTrade creates a new trade between a buyer and seller order
func (s *TradeService) CreateTrade(buyerOrder, sellerOrder models.IOrder) (*models.Trade, error) {
	if buyerOrder.Type() != models.Buy || sellerOrder.Type() != models.Sell {
		return nil, fmt.Errorf("invalid order types: buyer must be Buy and seller must be Sell")
	}
	if buyerOrder.Symbol() != sellerOrder.Symbol() {
		return nil, fmt.Errorf("stock symbols do not match")
	}
	if buyerOrder.StockQuantity() != sellerOrder.StockQuantity() {
		return nil, fmt.Errorf("quantities do not match")
	}
	if buyerOrder.StockPrice() != sellerOrder.StockPrice() {
		return nil, fmt.Errorf("prices do not match")
	}

	trade := &models.Trade{
		TradeID:        len(s.trades) + 1,
		BuyerOrderID:   buyerOrder.ID(),
		SellerOrderID:  sellerOrder.ID(),
		StockSymbol:    buyerOrder.Symbol(),
		Quantity:       buyerOrder.StockQuantity(),
		Price:          buyerOrder.StockPrice(),
		TradeTimestamp: time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.trades[trade.TradeID] = trade

	return trade, nil
}
