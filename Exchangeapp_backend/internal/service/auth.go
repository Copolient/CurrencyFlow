package service

import (
	"fmt"
	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	"exchangeapp/pkg/auth"
)

type AuthService struct {
	userRepo repository.UserRepository
	jwt      *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{userRepo: userRepo, jwt: jwt}
}

func (s *AuthService) Register(username, password string) (string, error) {
	hashedPwd, err := s.jwt.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Username: username,
		Password: hashedPwd,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	token, err := s.jwt.GenerateToken(username)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}

	if !s.jwt.CheckPassword(password, user.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := s.jwt.GenerateToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
