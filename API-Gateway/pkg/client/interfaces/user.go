package interfaces

import "api-gateway/pkg/utils/models"

type UserClient interface {
	UsersSignUp(user models.UserSignUp) (models.TokenUser, error)
	UserLogin(user models.UserLogin) (models.TokenUser, error)

	AddAddress(address models.Address) (models.Address, error)
	GetAddress(id uint) (models.Address, error)
	UpdateAddress(address models.Address) (models.Address, error)
	DeleteAddress(id uint) error
}
