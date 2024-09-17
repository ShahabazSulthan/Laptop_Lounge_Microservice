package interfaceUseCase

import (
	"user-service/pkg/domain"
	"user-service/pkg/models"
)

type UserUseCase interface {
	UsersSignUp(user models.UserSignUp) (domain.TokenUser, error)
	UsersLogin(user models.UserLogin) (domain.TokenUser, error)

	AddAddress(address models.AddressDetail) (models.AddressDetail, error)
	GetAddress(id uint) (domain.Address, error)
	UpdateAddress(address domain.Address) (domain.Address, error) 
	DeleteAddress(id uint) error

}
