package client

import (
	"context"
	"fmt"
	"order-service/pkg/config"
	"order-service/pkg/models"
	pbc "order-service/pkg/pb/cart"

	"google.golang.org/grpc"
)

type cartClient struct {
	Client pbc.CartClient
}

func NewCartClient(cfg *config.Config) *cartClient {

	grpcConnection, err := grpc.Dial(cfg.CartSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pbc.NewCartClient(grpcConnection)

	return &cartClient{
		Client: grpcClient,
	}

}

func (c *cartClient) GetAllItemsFromCart(userID int) ([]models.Cart, error) {
	res, err := c.Client.GetAllItemsFromCart(context.Background(), &pbc.GetAllItemsFromCartRequest{
		UserID: int64(userID),
	})
	if err != nil {
		return []models.Cart{}, err
	}
	var result []models.Cart

	for _, v := range res.Cart {
		result = append(result, models.Cart{
			ProductID:  uint(v.ProductID),
			Quantity:   float64(v.Quantity),
			TotalPrice: float64(v.TotalPrice),
		})
	}
	return result, nil
}
func (c *cartClient) DoesCartExist(userID int) (bool, error) {
	res, err := c.Client.DoesCartExist(context.Background(), &pbc.DoesCartExistRequest{
		UserID: int64(userID),
	})
	if err != nil {
		return false, err
	}
	return res.Data, nil
}
func (c *cartClient) TotalAmountInCart(userID int) (float64, error) {
	res, err := c.Client.TotalAmountInCart(context.Background(), &pbc.TotalAmountInCartRequest{
		UserID: int64(userID),
	})
	if err != nil {
		return 0.0, err
	}
	return float64(res.Data), err
}
func (c *cartClient) UpdateCartAfterOrder(userID, productID int, quantity float64) error {
	_, err := c.Client.UpdateCartAfterOrder(context.Background(), &pbc.UpdateCartAfterOrderRequest{
		UserID:    int64(userID),
		ProductID: int64(productID),
		Quantity:  int64(quantity),
	})
	if err != nil {
		return err
	}
	return nil
}
