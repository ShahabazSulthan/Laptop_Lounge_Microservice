package handler

import (
	"api-gateway/pkg/client/interfaces"
	"api-gateway/pkg/utils/models"
	"api-gateway/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	GRPC_Client interfaces.UserClient
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(UserClient interfaces.UserClient) *UserHandler {
	return &UserHandler{
		GRPC_Client: UserClient,
	}
}

// UserSignup handles user signup by making a call to the gRPC service
func (ur *UserHandler) UserSignup(c *gin.Context) {
	var signupDetail models.UserSignUp
	if err := c.ShouldBindJSON(&signupDetail); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Validate the struct fields
	err := validator.New().Struct(signupDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Call the gRPC client for user signup
	user, err := ur.GRPC_Client.UsersSignUp(signupDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Error during signup", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusCreated, "User successfully signed up", user, nil)
	c.JSON(http.StatusCreated, success)
}

// UserLogin handles user login by calling the gRPC service
func (ur *UserHandler) Userlogin(c *gin.Context) {
	var userLoginDetail models.UserLogin
	if err := c.ShouldBindJSON(&userLoginDetail); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Validate the login details
	err := validator.New().Struct(userLoginDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Call the gRPC client for user login
	user, err := ur.GRPC_Client.UserLogin(userLoginDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Error during login", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "User successfully logged in", user, nil)
	c.JSON(http.StatusOK, success)
}

// AddAddress handles the addition of a new address
func (ur *UserHandler) AddAddress(c *gin.Context) {
	var addressDetail models.Address
	if err := c.ShouldBindJSON(&addressDetail); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Validate the address details
	err := validator.New().Struct(addressDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Call the gRPC client to add the address
	address, err := ur.GRPC_Client.AddAddress(addressDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Error adding address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusCreated, "Address successfully added", address, nil)
	c.JSON(http.StatusCreated, success)
}

// GetAddress handles retrieval of an address by ID
func (ur *UserHandler) GetAddress(c *gin.Context) {
	id := c.Param("id")
	addressID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Invalid address ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Call the gRPC client to get the address
	address, err := ur.GRPC_Client.GetAddress(uint(addressID))
	if err != nil {
		errs := response.ClientResponse(http.StatusNotFound, "Address not found", nil, err.Error())
		c.JSON(http.StatusNotFound, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Address retrieved successfully", address, nil)
	c.JSON(http.StatusOK, success)
}

// UpdateAddress handles the update of an existing address
func (ur *UserHandler) UpdateAddress(c *gin.Context) {
	var addressDetail models.Address
	if err := c.ShouldBindJSON(&addressDetail); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Validate the address details
	err := validator.New().Struct(addressDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Call the gRPC client to update the address
	address, err := ur.GRPC_Client.UpdateAddress(addressDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Error updating address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Address successfully updated", address, nil)
	c.JSON(http.StatusOK, success)
}

// DeleteAddress handles removal of an address by ID
func (ur *UserHandler) DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	addressID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Invalid address ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	// Call the gRPC client to delete the address
	err = ur.GRPC_Client.DeleteAddress(uint(addressID))
	if err != nil {
		errs := response.ClientResponse(http.StatusNotFound, "Error deleting address", nil, err.Error())
		c.JSON(http.StatusNotFound, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Address successfully deleted", nil, nil)
	c.JSON(http.StatusOK, success)
}
