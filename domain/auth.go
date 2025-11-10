package domain

import (
	"context"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthUsecase interface {
	LoginUserByEmail(ctx context.Context, email string, password string, accessSecret string, accessExpiry int, refreshSecret string, refreshExpiry int) (loginResponse *LoginResponse, err error)
	LoginGuestUser(ctx context.Context, accessSecret string, accessExpiry int, refreshSecret string, refreshExpiry int) (loginResponse *LoginResponse, err error)
	CreateAccessToken(user *User, accessSecret string, accessExpiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, refreshSecret string, refreshExpiry int) (refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string, refreshSecret string, accessSecret string, accessExpiry int, refreshExpiry int) (refreshResponse *RefreshTokenResponse, err error)

	ForgotPassword(ctx context.Context, email string) (err error)
	ResetPassword(ctx context.Context, email string, newPassword string) (err error)
}
