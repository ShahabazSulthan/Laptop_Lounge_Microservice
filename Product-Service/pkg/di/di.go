package di

import (
	server "product-service/pkg/api"
	"product-service/pkg/api/service"
	"product-service/pkg/config"
	"product-service/pkg/db"
	"product-service/pkg/repository"
	"product-service/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*server.Server, error) {

	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}

	adminRepository := repository.NewProductRepository(gormDB)
	adminUseCase := usecase.NewProductUseCase(adminRepository)

	adminServiceServer := service.NewProductServer(adminUseCase)
	grpcServer, err := server.NewGRPCServer(cfg, adminServiceServer)

	if err != nil {
		return &server.Server{}, err
	}

	return grpcServer, nil

}
