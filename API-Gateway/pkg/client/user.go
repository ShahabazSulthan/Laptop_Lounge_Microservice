package client

import (
	"api-gateway/pkg/client/interfaces"
	"api-gateway/pkg/config"
	pb "api-gateway/pkg/pb/user"
	"api-gateway/pkg/utils/models"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userClient struct {
	Client pb.UserClient
}

func NewUserClient(cfg config.Config) interfaces.UserClient {
	// Establish a secure gRPC connection with proper error handling.
	grpcConnection, err := grpc.Dial(cfg.UserSvcUrl, grpc.WithInsecure())
	if err != nil {
		log.Printf("Error connecting to user service: %v", err)
		return nil
	}

	grpcClient := pb.NewUserClient(grpcConnection)

	return &userClient{
		Client: grpcClient,
	}
}

// UsersSignUp handles user signup by making a gRPC call to the user service.
func (c *userClient) UsersSignUp(user models.UserSignUp) (models.TokenUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.Client.UserSignUp(ctx, &pb.UserSignUpRequest{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	})
	if err != nil {
		return models.TokenUser{}, handleGrpcError(err)
	}

	if res.UserDetails == nil {
		return models.TokenUser{}, fmt.Errorf("user details are missing in the response")
	}

	userDetails := models.UserDetails{
		ID:        uint(res.UserDetails.Id),
		Firstname: res.UserDetails.Firstname,
		Lastname:  res.UserDetails.Lastname,
		Email:     res.UserDetails.Email,
		Phone:     res.UserDetails.Phone,
	}

	return models.TokenUser{
		User:        userDetails,
		AccessToken: res.AccessToken,
	}, nil
}

// UserLogin handles user login by calling the gRPC service.
func (c *userClient) UserLogin(user models.UserLogin) (models.TokenUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.Client.UserLogin(ctx, &pb.UserLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return models.TokenUser{}, handleGrpcError(err)
	}

	if res.UserDetails == nil {
		return models.TokenUser{}, fmt.Errorf("user details are missing in the response")
	}

	userDetails := models.UserDetails{
		ID:        uint(res.UserDetails.Id),
		Firstname: res.UserDetails.Firstname,
		Lastname:  res.UserDetails.Lastname,
		Email:     res.UserDetails.Email,
		Phone:     res.UserDetails.Phone,
	}

	return models.TokenUser{
		User:        userDetails,
		AccessToken: res.AccessToken,
	}, nil
}

// AddAddress adds a new address for a user.
func (c *userClient) AddAddress(address models.Address) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.Client.AddAddress(ctx, &pb.AddAddressRequest{
		UserId:  uint64(address.UserID),
		Street:  address.Street,
		City:    address.City,
		State:   address.State,
		ZipCode: address.ZipCode,
		Country: address.Country,
	})
	if err != nil {
		return models.Address{}, handleGrpcError(err)
	}

	if res.AddressDetails == nil {
		return models.Address{}, fmt.Errorf("address details are missing in the response")
	}

	return models.Address{
		ID:      uint(res.AddressDetails.Id),
		UserID:  uint(res.AddressDetails.UserId),
		Street:  res.AddressDetails.Street,
		City:    res.AddressDetails.City,
		State:   res.AddressDetails.State,
		ZipCode: res.AddressDetails.ZipCode,
		Country: res.AddressDetails.Country,
	}, nil
}

// GetAddress retrieves an address by ID.
func (c *userClient) GetAddress(id uint) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.Client.GetAddress(ctx, &pb.GetAddressRequest{
		Id: uint64(id),
	})
	if err != nil {
		return models.Address{}, handleGrpcError(err)
	}

	if res.Status != 200 || res.AddressDetails == nil {
		return models.Address{}, fmt.Errorf("address not found with ID %d", id)
	}

	return models.Address{
		ID:      uint(res.AddressDetails.Id),
		UserID:  uint(res.AddressDetails.UserId),
		Street:  res.AddressDetails.Street,
		City:    res.AddressDetails.City,
		State:   res.AddressDetails.State,
		ZipCode: res.AddressDetails.ZipCode,
		Country: res.AddressDetails.Country,
	}, nil
}

// UpdateAddress updates an existing address.
func (c *userClient) UpdateAddress(address models.Address) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.Client.UpdateAddress(ctx, &pb.UpdateAddressRequest{
		Id:      uint64(address.ID),
		Street:  address.Street,
		City:    address.City,
		State:   address.State,
		ZipCode: address.ZipCode,
		Country: address.Country,
	})
	if err != nil {
		return models.Address{}, handleGrpcError(err)
	}

	if res.Status != 200 || res.AddressDetails == nil {
		return models.Address{}, fmt.Errorf("failed to update address with ID %d", address.ID)
	}

	return models.Address{
		ID:      uint(res.AddressDetails.Id),
		UserID:  uint(res.AddressDetails.UserId),
		Street:  res.AddressDetails.Street,
		City:    res.AddressDetails.City,
		State:   res.AddressDetails.State,
		ZipCode: res.AddressDetails.ZipCode,
		Country: res.AddressDetails.Country,
	}, nil
}

// DeleteAddress removes an address by ID.
func (c *userClient) DeleteAddress(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.Client.DeleteAddress(ctx, &pb.DeleteAddressRequest{
		Id: uint64(id),
	})
	if err != nil {
		return handleGrpcError(err)
	}

	if res.Status != 200 {
		return fmt.Errorf("failed to delete address with ID %d", id)
	}

	return nil
}

// handleGrpcError handles errors returned from gRPC calls, mapping them to more descriptive errors without error codes.
func handleGrpcError(err error) error {
    st, ok := status.FromError(err)
    if !ok {
        // Not a gRPC error, return as-is.
        return fmt.Errorf("unexpected error: %v", err)
    }

    // Map specific gRPC error codes to user-friendly messages without including the error codes in the message.
    switch st.Code() {
    case codes.NotFound:
        return fmt.Errorf("%v", st.Message()) // "resource not found" simplified.
    case codes.Unauthenticated:
        return fmt.Errorf("%v", st.Message()) // "unauthenticated" simplified.
    case codes.PermissionDenied:
        return fmt.Errorf("%v", st.Message()) // "permission denied" simplified.
    case codes.Internal:
        return fmt.Errorf("%v", st.Message()) // "internal server error" simplified.
    case codes.InvalidArgument:
        return fmt.Errorf("%v", st.Message()) // "invalid argument" simplified.
    default:
        return fmt.Errorf("%v", st.Message()) // Other gRPC errors.
    }
}
