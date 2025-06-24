package service

import (
	"auth-service/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) CreateUser(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) DeleteUserByEmail(email string) error {
	result := s.DB.Where("email = ?", email).Delete(&model.User{})
	return result.Error
}
