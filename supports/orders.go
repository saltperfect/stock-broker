package supports

import (
	"design_patterns/stock_broker/models"
	"sync"
)

type IOrderService interface {
	AddOrder()
}

type OrderService struct {
	orders map[int]models.IOrder
	mu     sync.Mutex

	orderProcessor *OrderProcessor
}

// NewOrderService creates a new OrderService instance
func NewOrderService(op *OrderProcessor) *OrderService {
	return &OrderService{
		orders:         make(map[int]models.IOrder),
		orderProcessor: op,
	}
}

// AddOrder adds a new order to the service
func (s *OrderService) AddOrder(userID string, orderType models.OrderType, stockSymbol string, quantity int, price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order := models.NewOrder(userID, orderType, stockSymbol, quantity, price)
	s.orders[order.ID()] = order
	if order.Type() == models.Buy {
		s.orderProcessor.updateBuy(order)
	} else {
		s.orderProcessor.updateSell(order)
	}
}

// GetOrder retrieves an order by its ID
func (s *OrderService) GetOrder(orderID int) (models.IOrder, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, exists := s.orders[orderID]
	return order, exists
}
