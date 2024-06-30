package supports

import (
	"design_patterns/stock_broker/models"
	"time"
)

type IMatchingStratergy interface {
	getLatestOrder([]models.IOrder) models.IOrder
}

type getOldestOrderStratergy struct {
}

func (g getOldestOrderStratergy) getLatestOrder(orders []models.IOrder) models.IOrder {
	oldestOrderTime := time.Now()
	var oldestOrder models.IOrder

	for _, order := range orders {
		if oldestOrderTime.Before(order.OrderTimestamp()) {
			oldestOrderTime = order.OrderTimestamp()
			oldestOrder = order
		}
	}
	return oldestOrder
}
