package service

import (
	"context"
	"errors"
	"user-service/pkg/domain"
	"user-service/pkg/models"
	"user-service/pkg/pb"
	interfaceUseCase "user-service/pkg/usecase/interface"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userUseCase interfaceUseCase.UserUseCase
	pb.UnimplementedUserServer
}

func NewUserServer(useCase interfaceUseCase.UserUseCase) pb.UserServer {
	return &UserServer{
		userUseCase: useCase,
	}
}

// UserSignUp handles the signup process for new users.
func (s *UserServer) UserSignUp(ctx context.Context, userSignUpDetails *pb.UserSignUpRequest) (*pb.UserSignUpResponse, error) {
	userCreateDetails := models.UserSignUp{
		Firstname: userSignUpDetails.Firstname,
		Lastname:  userSignUpDetails.Lastname,
		Email:     userSignUpDetails.Email,
		Phone:     userSignUpDetails.Phone,
		Password:  userSignUpDetails.Password,
	}

	data, err := s.userUseCase.UsersSignUp(userCreateDetails)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign up user: %v", err)
	}

	userDetails := &pb.UserDetails{
		Id:        uint64(data.User.ID),
		Firstname: data.User.Firstname,
		Lastname:  data.User.Lastname,
		Email:     data.User.Email,
		Phone:     data.User.Phone,
	}

	return &pb.UserSignUpResponse{
		Status:      201,
		UserDetails: userDetails,
		AccessToken: data.AccessToken,
	}, nil
}

// UserLogin handles the login process for users.
func (s *UserServer) UserLogin(ctx context.Context, loginDetails *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	login := models.UserLogin{
		Email:    loginDetails.Email,
		Password: loginDetails.Password,
	}

	data, err := s.userUseCase.UsersLogin(login)
	if err != nil {
		// If specific error handling is required, check for that error type
		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
	}

	userDetails := &pb.UserDetails{
		Id:        uint64(data.User.ID),
		Firstname: data.User.Firstname,
		Lastname:  data.User.Lastname,
		Email:     data.User.Email,
		Phone:     data.User.Phone,
	}

	return &pb.UserLoginResponse{
		Status:      200,
		UserDetails: userDetails,
		AccessToken: data.AccessToken,
	}, nil
}

// AddAddress adds a new address for a user.
func (s *UserServer) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	addressDetails := models.AddressDetail{
		UserID:  uint(req.UserId),
		Street:  req.Street,
		City:    req.City,
		State:   req.State,
		ZipCode: req.ZipCode,
		Country: req.Country,
	}

	newAddress, err := s.userUseCase.AddAddress(addressDetails)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add address: %v", err)
	}

	return &pb.AddAddressResponse{
		Status: 201,
		AddressDetails: &pb.AddressDetails{
			Id:      uint64(newAddress.ID),
			UserId:  uint64(newAddress.UserID),
			Street:  newAddress.Street,
			City:    newAddress.City,
			State:   newAddress.State,
			ZipCode: newAddress.ZipCode,
			Country: newAddress.Country,
		},
	}, nil
}

// GetAddress retrieves a user's address by its ID.
func (s *UserServer) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.GetAddressResponse, error) {
	address, err := s.userUseCase.GetAddress(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get address: %v", err)
	}

	return &pb.GetAddressResponse{
		Status: 200,
		AddressDetails: &pb.AddressDetails{
			Id:      uint64(address.ID),
			UserId:  uint64(address.UserID),
			Street:  address.Street,
			City:    address.City,
			State:   address.State,
			ZipCode: address.ZipCode,
			Country: address.Country,
		},
	}, nil
}

// UpdateAddress updates an existing address.
func (s *UserServer) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	// Convert protobuf request to domain model
	addressDetails := domain.Address{
		ID:      uint(req.Id),
		Street:  req.Street,
		City:    req.City,
		State:   req.State,
		ZipCode: req.ZipCode,
		Country: req.Country,
	}

	// Call the use case to update the address
	updatedAddress, err := s.userUseCase.UpdateAddress(addressDetails)
	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			// Use a specific gRPC error code for not found
			return nil, status.Errorf(codes.NotFound, "address with ID %d not found", req.Id)
		}
		// General internal error
		return nil, status.Errorf(codes.Internal, "failed to update address: %v", err)
	}

	// Prepare response with updated address details
	return &pb.UpdateAddressResponse{
		Status: 200,
		AddressDetails: &pb.AddressDetails{
			Id:      uint64(updatedAddress.ID),
			UserId:  uint64(updatedAddress.UserID),
			Street:  updatedAddress.Street,
			City:    updatedAddress.City,
			State:   updatedAddress.State,
			ZipCode: updatedAddress.ZipCode,
			Country: updatedAddress.Country,
		},
	}, nil
}

// DeleteAddress deletes a user's address by its ID.
func (s *UserServer) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	err := s.userUseCase.DeleteAddress(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete address: %v", err)
	}

	return &pb.DeleteAddressResponse{
		Status: 200,
	}, nil
}
