package service

import (
    "context"
    "errors"

    "github.com/FUADIKAMIL/taskify/internal/model"
    "github.com/FUADIKAMIL/taskify/internal/repository"
    "github.com/FUADIKAMIL/taskify/pkg/auth"
    "golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
    Register func(ctx context.Context, username, password string) (model.User, error)
    Login    func(ctx context.Context, username, password string) (string, error)
}

func NewAuthService(userRepo *repository.UserRepo) *AuthService {

    registerFn := func(ctx context.Context, username, password string) (model.User, error) {
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            return model.User{}, err
        }
        u := model.User{Username: username, Password: string(hashed)}
        return userRepo.Create(ctx, u)
    }

    loginFn := func(ctx context.Context, username, password string) (string, error) {
        u, err := userRepo.GetByUsername(ctx, username)
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

    return &AuthService{
        Register: registerFn,
        Login:    loginFn,
    }
}

