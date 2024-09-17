package di

import (
	server "order-service/pkg/api"
	"order-service/pkg/api/services"
	"order-service/pkg/client"
	"order-service/pkg/config"
	"order-service/pkg/db"
	"order-service/pkg/repository"
	"order-service/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.Server, error) {

	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}

	orderRepository := repository.NewOrderRepository(gormDB)
	cartClient := client.NewCartClient(&cfg)
	productClient := client.NewProductClient(&cfg)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, cartClient, productClient)

	orderServiceServer := services.NewOrderServer(orderUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, orderServiceServer)

	if err != nil {
		return &server.Server{}, err
	}

	return grpcServer, nil

}
