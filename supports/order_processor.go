package supports

import (
	"container/heap"
	"design_patterns/stock_broker/models"
	"fmt"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type key struct {
	symbol   string
	price    float64
	quantity int
}

type OrderProcessor struct {
	mutex sync.Mutex
	Buy   map[key]*models.OrderHeap
	Sell  map[key]*models.OrderHeap

	tradeService *TradeService

	tradeTicker time.Ticker

	getmatchingOrderStratergy IMatchingStratergy
}

func NewOrderProcessor(ts *TradeService) *OrderProcessor {
	return &OrderProcessor{
		Buy:                       make(map[key]*models.OrderHeap),
		Sell:                      make(map[key]*models.OrderHeap),
		tradeService:              ts,
		tradeTicker:               *time.NewTicker(100 * time.Millisecond),
		getmatchingOrderStratergy: getOldestOrderStratergy{},
	}
}

func getKey(symbol string, price float64, quantity int) key {
	return key{
		symbol:   symbol,
		price:    price,
		quantity: quantity,
	}
}

func (op *OrderProcessor) updateSell(o models.IOrder) {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.Sell[getKey(o.Symbol(), o.StockPrice(), o.StockQuantity())] == nil {
		h := &models.OrderHeap{}
		heap.Init(h)
		op.Sell[getKey(o.Symbol(), o.StockPrice(), o.StockQuantity())] = h
	}
	heap.Push(op.Sell[getKey(o.Symbol(), o.StockPrice(), o.StockQuantity())], o)
}

func (op *OrderProcessor) updateBuy(o models.IOrder) {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.Buy[getKey(o.Symbol(), o.StockPrice(), o.StockQuantity())] == nil {
		h := &models.OrderHeap{}
		heap.Init(h)
		op.Buy[getKey(o.Symbol(), o.StockPrice(), o.StockQuantity())] = h
	}
	heap.Push(op.Buy[getKey(o.Symbol(), o.StockPrice(), o.StockQuantity())], o)
}

func (op *OrderProcessor) GetMatches() {

	op.mutex.Lock()
	for key, buyHeap := range op.Buy {
		if sellHeap, ok := op.Sell[key]; ok {
			buyOrder := (*buyHeap)[0]
			sellOrder := (*sellHeap)[0]
			trade, err := op.tradeService.CreateTrade(buyOrder, sellOrder)
			if err != nil {
				fmt.Printf("unable to trade buy: %v, sell: %v\n", buyOrder, sellOrder)
				continue
			}
			buyHeap.Pop()
			sellHeap.Pop()
			buyOrder.Execute()
			sellOrder.Execute()
			spew.Printf("trade done: %v\n", trade)
		}
	}
	op.mutex.Unlock()
}
