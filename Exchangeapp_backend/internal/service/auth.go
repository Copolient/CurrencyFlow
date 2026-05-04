package service

import (
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
		return "", err
	}

	user := &model.User{
		Username: username,
		Password: hashedPwd,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	return s.jwt.GenerateToken(username)
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	if !s.jwt.CheckPassword(password, user.Password) {
		return "", err
	}

	return s.jwt.GenerateToken(user.Username)
}
