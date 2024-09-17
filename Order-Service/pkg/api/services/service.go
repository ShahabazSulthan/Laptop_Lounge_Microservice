package services

import (
	"context"
	"order-service/pkg/models"
	pb "order-service/pkg/pb/order"
	interfaceUse "order-service/pkg/usecase/interfaces"
)

type OrderServer struct {
	UseCase interfaceUse.OrderUseCase
	pb.UnimplementedOrderServer
}

func NewOrderServer(useCase interfaceUse.OrderUseCase) pb.OrderServer {
	return &OrderServer{
		UseCase: useCase,
	}
}

func (or *OrderServer) OrderItemsFromCart(ctx context.Context, req *pb.OrderItemsFromCartRequest) (*pb.OrderItemsFromCartResponse, error) {
	model := &models.OrderFromCart{
		AddressID: uint(req.OrderFromCart.AddressID),
		PaymentID: uint(req.OrderFromCart.PaymentID),
	}
	userID := int(req.UserID)
	result, err := or.UseCase.OrderItemsFromCart(*model, userID)
	if err != nil {
		return &pb.OrderItemsFromCartResponse{}, err
	}
	return &pb.OrderItemsFromCartResponse{
		OrderID:        int64(result.OrderID),
		Shipmentstatus: result.ShipmentStatus,
	}, nil
}
func (or *OrderServer) GetOrderDetails(ctx context.Context, req *pb.GetOrderDetailsRequest) (*pb.GetOrderDetailsResponse, error) {
	userID := int(req.UserID)
	page := int(req.Page)
	count := int(req.Count)

	details, err := or.UseCase.GetOrderDetails(userID, page, count)
	if err != nil {
		return nil, err
	}

	var result pb.GetOrderDetailsResponse

	for _, v := range details {
		var orderDetails pb.OrderDetails
		orderDetails.OrderID = int64(v.OrderDetails.OrderId)
		orderDetails.Price = float32(v.OrderDetails.FinalPrice)
		orderDetails.Shipmentstatus = v.OrderDetails.ShipmentStatus
		orderDetails.Paymentstatus = v.OrderDetails.PaymentStatus

		var orderProductDetails []*pb.OrderProductDetails
		for _, product := range v.OrderProductDetails {
			orderProduct := &pb.OrderProductDetails{
				ProductID: int64(product.ProductID),
				Quantity:  int64(product.Quantity),
				Price:     float32(product.TotalPrice),
			}
			orderProductDetails = append(orderProductDetails, orderProduct)
		}

		fullOrderDetails := &pb.FullOrderDetails{
			Orderdetails:        &orderDetails,
			OrderProductDetails: orderProductDetails,
		}

		result.Details = append(result.Details, fullOrderDetails)
	}

	return &result, nil
}
