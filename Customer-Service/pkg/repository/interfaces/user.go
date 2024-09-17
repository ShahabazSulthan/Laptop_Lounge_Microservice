package interfaces

import (
	"user-service/pkg/domain"
	"user-service/pkg/models"
)

type UserRepository interface {
	CheckUserExistsByEmail(email string) (*domain.User, error)
	CheckUserExistsByPhone(phone string) (*domain.User, error)
	UserSignUp(user models.UserSignUp) (models.UserDetails, error)
	FindUserByEmail(user models.UserLogin) (models.UserDetail, error)

	AddAddress(address domain.Address) (domain.Address, error)
	GetAddressByID(id uint) (domain.Address, error)
	UpdateAddress(address domain.Address) (domain.Address, error)
	DeleteAddress(id uint) error
}
