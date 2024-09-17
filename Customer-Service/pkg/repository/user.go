package repository

import (
	"errors"
	"user-service/pkg/domain"
	"user-service/pkg/models"
	"user-service/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		DB: DB,
	}
}

// CheckUserExistsByEmail checks if a user exists by email.
func (ur *userRepository) CheckUserExistsByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CheckUserExistsByPhone checks if a user exists by phone.
func (ur *userRepository) CheckUserExistsByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := ur.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UserSignUp signs up a new user.
func (ur *userRepository) UserSignUp(user models.UserSignUp) (models.UserDetails, error) {
	var signupDetail models.UserDetails

	err := ur.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO users (firstname, lastname, email, password, phone)
			VALUES (?, ?, ?, ?, ?)
		`, user.Firstname, user.Lastname, user.Email, user.Password, user.Phone).Error; err != nil {
			return err
		}

		if err := tx.Raw(`
			SELECT id, firstname, lastname, email, phone FROM users WHERE email = ?
		`, user.Email).Scan(&signupDetail).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.UserDetails{}, err
	}

	return signupDetail, nil
}

// FindUserByEmail finds a user by email.
func (ur *userRepository) FindUserByEmail(user models.UserLogin) (models.UserDetail, error) {
	var userDetails models.UserDetail
	query := `SELECT id, firstname, lastname, email, password, phone FROM users WHERE email = ?`
	err := ur.DB.Raw(query, user.Email).Scan(&userDetails).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserDetail{}, errors.New("user not found")
		}
		return models.UserDetail{}, err
	}
	return userDetails, nil
}


// AddAddress adds a new address.
func (ur *userRepository) AddAddress(address domain.Address) (domain.Address, error) {
	err := ur.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO addresses (user_id, street, city, state, zip_code, country)
			VALUES (?, ?, ?, ?, ?, ?)
		`, address.UserID, address.Street, address.City, address.State, address.ZipCode, address.Country).Error; err != nil {
			return err
		}

		if err := tx.Raw(`
			SELECT id, user_id, street, city, state, zip_code, country FROM addresses
			WHERE user_id = ? AND street = ? AND city = ? AND state = ? AND zip_code = ? AND country = ?
		`, address.UserID, address.Street, address.City, address.State, address.ZipCode, address.Country).Scan(&address).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return domain.Address{}, err
	}
	return address, nil
}

// GetAddressByID retrieves an address by ID.
func (ur *userRepository) GetAddressByID(id uint) (domain.Address, error) {
	var address domain.Address
	err := ur.DB.First(&address, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Address{}, nil
		}
		return domain.Address{}, err
	}
	return address, nil
}

// UpdateAddress updates an existing address.
func (ur *userRepository) UpdateAddress(address domain.Address) (domain.Address, error) {
	err := ur.DB.Model(&address).Where("id = ?", address.ID).Updates(domain.Address{
		Street:   address.Street,
		City:     address.City,
		State:    address.State,
		ZipCode:  address.ZipCode,
		Country:  address.Country,
	}).Error

	if err != nil {
		return domain.Address{}, err
	}

	// Refresh the updated address details
	err = ur.DB.First(&address, address.ID).Error
	if err != nil {
		return domain.Address{}, err
	}

	return address, nil
}

// DeleteAddress deletes an address by ID.
func (ur *userRepository) DeleteAddress(id uint) error {
	// Step 1: Check if the address exists
	var address domain.Address
	err := ur.DB.First(&address, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrAddressNotFound // Address not found
		}
		return err // Other errors
	}

	// Step 2: Delete the address
	err = ur.DB.Delete(&domain.Address{}, id).Error
	if err != nil {
		return err // Return error if deletion fails
	}
	return nil
}
