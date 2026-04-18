package service

import (
	"context"
	"errors"
	"inventory-app/internal/config"
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	cnf      *config.Config
	authRepo model.UserAuthRepository
}

func NewAuth(cnf *config.Config, authRepo model.UserAuthRepository) model.AuthService {
	return &authService{
		cnf:      cnf,
		authRepo: authRepo,
	}
}

// Login implements [model.AuthService].
func (a *authService) Login(ctx context.Context, req dto.AuthRequest) (resp dto.AuthResponse, err error) {
	user, err := a.authRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if user.Status != model.UserStatusAktif {
		return dto.AuthResponse{}, errors.New("User sudah tidak aktif")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		_ = a.authRepo.HandleLoginResult(ctx, req.Username, false)
		return dto.AuthResponse{}, errors.New("Invalid Credentials")
	}

	_ = a.authRepo.HandleLoginResult(ctx, req.Username, true)

	expTime := time.Now().Add(time.Duration(a.cnf.Jwt.Exp) * time.Minute)

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"status":   user.Status,
		"exp":      expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.cnf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		TypeToken: "Bearer",
		Token:     tokenString,
	}, nil
}
