package interfaces

import (
	"order-service/pkg/domain"
	"order-service/pkg/models"
)

type OrderRepository interface {
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	PaymentExist(orderBody models.OrderIncoming) (bool, error)
	PaymentStatus(orderID int) (string, error)
	OrderItems(ob models.OrderIncoming, price float64) (int, error)
	AddOrderProducts(order_id int, cart []models.Cart) error
	GetBriefOrderDetails(orderID int) (domain.OrderSuccessResponse, error)
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
}
