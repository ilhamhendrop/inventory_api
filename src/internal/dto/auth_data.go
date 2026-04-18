package dto

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	TypeToken string `json:"token_type"`
	Token     string `json:"token"`
}
