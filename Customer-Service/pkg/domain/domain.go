package domain

import (
	"errors"
	"user-service/pkg/models"
)

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Firstname string `json:"firstname" gorm:"not null"`
	Lastname  string `json:"lastname" gorm:"not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
	Phone     string `json:"phone" gorm:"not null"`
}

type TokenUser struct {
	User         models.UserDetails
	AccessToken  string
	RefreshToken string
}

type Address struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID  uint   `json:"user_id" gorm:"not null"`
	Street  string `json:"street" gorm:"not null"`
	City    string `json:"city" gorm:"not null"`
	State   string `json:"state" gorm:"not null"`
	ZipCode string `json:"zip_code" gorm:"not null"`
	Country string `json:"country" gorm:"not null"`
}

var (
	ErrAddressNotFound    = errors.New("address not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
