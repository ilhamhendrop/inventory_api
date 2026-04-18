package model

import (
	"context"
	"inventory-app/internal/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthRequest) (resp dto.AuthResponse, err error)
}
