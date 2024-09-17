package usecase

import (
	"errors"
	"user-service/pkg/domain"
	"user-service/pkg/helper"
	"user-service/pkg/models"
	"user-service/pkg/repository/interfaces"
	interfaceUseCase "user-service/pkg/usecase/interface"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUseCase struct {
	userRepository interfaces.UserRepository
}

func NewUserUseCase(repository interfaces.UserRepository) interfaceUseCase.UserUseCase {
	return &userUseCase{
		userRepository: repository,
	}
}

// UsersSignUp handles user registration.
func (uc *userUseCase) UsersSignUp(user models.UserSignUp) (domain.TokenUser, error) {
	// Check if email exists
	if existingUser, err := uc.userRepository.CheckUserExistsByEmail(user.Email); err != nil {
		return domain.TokenUser{}, errors.New("server error during email check")
	} else if existingUser != nil {
		return domain.TokenUser{}, errors.New("user with this email already exists")
	}

	// Check if phone exists
	if existingUser, err := uc.userRepository.CheckUserExistsByPhone(user.Phone); err != nil {
		return domain.TokenUser{}, errors.New("server error during phone check")
	} else if existingUser != nil {
		return domain.TokenUser{}, errors.New("user with this phone already exists")
	}

	// Hash password
	hashedPassword, err := helper.PasswordHash(user.Password)
	if err != nil {
		return domain.TokenUser{}, errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Create user
	userData, err := uc.userRepository.UserSignUp(user)
	if err != nil {
		return domain.TokenUser{}, errors.New("failed to sign up user")
	}

	// Generate tokens
	accessToken, err := helper.GenerateAccessToken(userData)
	if err != nil {
		return domain.TokenUser{}, errors.New("failed to generate access token")
	}
	refreshToken, err := helper.GenerateRefreshToken(userData)
	if err != nil {
		return domain.TokenUser{}, errors.New("failed to generate refresh token")
	}

	return domain.TokenUser{
		User:         userData,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ur *userUseCase) UsersLogin(user models.UserLogin) (domain.TokenUser, error) {
	email, err := ur.userRepository.CheckUserExistsByEmail(user.Email)
	if err != nil {
		return domain.TokenUser{}, errors.New("error with server")
	}
	if email == nil {
		return domain.TokenUser{}, errors.New("email doesn't exist")
	}
	userdeatils, err := ur.userRepository.FindUserByEmail(user)
	if err != nil {
		return domain.TokenUser{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userdeatils.Password), []byte(user.Password))
	if err != nil {
		return domain.TokenUser{}, errors.New("password not matching")
	}
	var userDetails models.UserDetails
	err = copier.Copy(&userDetails, &userdeatils)
	if err != nil {
		return domain.TokenUser{}, err
	}
	accessToken, err := helper.GenerateAccessToken(userDetails)
	if err != nil {
		return domain.TokenUser{}, errors.New("couldn't create accesstoken due to internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userDetails)
	if err != nil {
		return domain.TokenUser{}, errors.New("counldn't create refreshtoken due to internal error")
	}
	return domain.TokenUser{
		User:         userDetails,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// AddAddress adds an address to the user's account.
func (uc *userUseCase) AddAddress(address models.AddressDetail) (models.AddressDetail, error) {
	addressData := domain.Address{}
	if err := copier.Copy(&addressData, &address); err != nil {
		return models.AddressDetail{}, err
	}

	newAddress, err := uc.userRepository.AddAddress(addressData)
	if err != nil {
		return models.AddressDetail{}, err
	}

	var createdAddress models.AddressDetail
	if err := copier.Copy(&createdAddress, &newAddress); err != nil {
		return models.AddressDetail{}, err
	}

	return createdAddress, nil
}

// GetAddress retrieves an address by ID through the use case
func (au *userUseCase) GetAddress(id uint) (domain.Address, error) {
	address, err := au.userRepository.GetAddressByID(id)
	if err != nil {
		return domain.Address{}, errors.New("could not retrieve address: " + err.Error())
	}
	if address.ID == 0 {
		return domain.Address{}, errors.New("address not found")
	}
	return address, nil
}

// UpdateAddress updates an address through the use case
func (au *userUseCase) UpdateAddress(address domain.Address) (domain.Address, error) {
	updatedAddress, err := au.userRepository.UpdateAddress(address)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return a specific error if the address is not found
			return domain.Address{}, domain.ErrAddressNotFound
		}
		// General error message
		return domain.Address{}, errors.New("could not update address: " + err.Error())
	}
	return updatedAddress, nil
}

// DeleteAddress deletes an address by ID through the use case
func (au *userUseCase) DeleteAddress(id uint) error {
	err := au.userRepository.DeleteAddress(id)
	if err != nil {
		return errors.New("could not delete address: " + err.Error())
	}
	return nil
}
