package service

import (
	"errors"

	"auth-service/model"
	"auth-service/repository"

	"golang.org/x/crypto/bcrypt"

	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: repo}
}

func (s *AuthService) Register(input model.RegisterRequest) (*model.User, error) {
	existing, _ := s.UserRepo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("el correo ya está registrado")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(input.Email, string(hashed))
	if err := s.UserRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input model.LoginRequest, redisClient *redis.Client) (string, error) {
	user, err := s.UserRepo.FindByEmail(input.Email)
	if err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return "", errors.New("usuario o contraseña inválidos")
	}

	token := uuid.NewString()
	ttl := time.Duration(3600) * time.Second
	if customTTL := os.Getenv("SESSION_TTL"); customTTL != "" {
		if seconds, err := time.ParseDuration(customTTL + "s"); err == nil {
			ttl = seconds
		}
	}

	err = redisClient.Set(context.Background(), "session:"+token, user.ID.String(), ttl).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}
