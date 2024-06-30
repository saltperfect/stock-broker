package models

import (
	"math/rand"
	"time"
)

type IOrder interface {
	ID() int
	OrderStatus() OrderStatus
	OrderTimestamp() time.Time
	Symbol() string
	StockPrice() float64
	StockQuantity() int
	Type() OrderType
	Execute()
}

// OrderType represents the type of order: Buy or Sell
type OrderType string

const (
	Buy  OrderType = "Buy"
	Sell OrderType = "Sell"
)

// OrderStatus represents the status of the order: ACCEPTED, REJECTED, or CANCELED
type OrderStatus string

const (
	ACCEPTED OrderStatus = "ACCEPTED"
	EXPIRED  OrderStatus = "EXPIRED"
	EXECUTED OrderStatus = "EXECUTED"
	CANCELED OrderStatus = "CANCELED"
)

// Order represents a stock order
type Order struct {
	OrderID                int
	UserID                 string
	OrderType              OrderType
	StockSymbol            string
	Quantity               int
	Price                  float64
	OrderAcceptedTimestamp time.Time
	Status                 OrderStatus
}

func randInt() int {
	return rand.Intn(100)
}

// NewOrder is a constructor function to create a new Order instance
func NewOrder(userID string, orderType OrderType, stockSymbol string, quantity int, price float64) IOrder {
	return &Order{
		OrderID:                randInt(),
		UserID:                 userID,
		OrderType:              orderType,
		StockSymbol:            stockSymbol,
		Quantity:               quantity,
		Price:                  price,
		OrderAcceptedTimestamp: time.Now(),
		Status:                 ACCEPTED,
	}
}

func (o *Order) OrderStatus() OrderStatus {
	return o.Status
}
func (o *Order) ID() int {
	return o.OrderID
}

func (o *Order) OrderTimestamp() time.Time {
	return o.OrderAcceptedTimestamp
}

func (o *Order) StockPrice() float64 {
	return o.Price
}

func (o *Order) StockQuantity() int {
	return o.Quantity
}

func (o *Order) Symbol() string {
	return o.StockSymbol
}

func (o *Order) Type() OrderType {
	return o.OrderType
}

func (o *Order) Execute() {
	o.OrderType = OrderType(EXECUTED)
}

type OrderHeap []IOrder

func (h OrderHeap) Len() int           { return len(h) }
func (h OrderHeap) Less(i, j int) bool { return h[i].OrderTimestamp().Before(h[j].OrderTimestamp()) }
func (h OrderHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *OrderHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(IOrder))
}

func (h *OrderHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
