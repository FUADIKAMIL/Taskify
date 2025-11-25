package service

import (
	"context"
	"errors"

	"github.com/yourname/taskify/internal/model"
	"github.com/yourname/taskify/internal/repository"
	"github.com/yourname/taskify/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct{ userRepo *repository.UserRepo }

func NewAuthService(u *repository.UserRepo) *AuthService { return &AuthService{userRepo: u} }

func (s *AuthService) Register(ctx context.Context, username, password string) (model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	u := model.User{Username: username, Password: string(hashed)}
	return s.userRepo.Create(ctx, u)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	tok, err := auth.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}
	return tok, nil
}
