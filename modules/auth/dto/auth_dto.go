package dto

import (
	"errors"
)

const (
	MESSAGE_FAILED_REFRESH_TOKEN  = "failed refresh token"
	MESSAGE_SUCCESS_REFRESH_TOKEN = "success refresh token"
	MESSAGE_FAILED_LOGOUT         = "failed logout"
	MESSAGE_SUCCESS_LOGOUT        = "success logout"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrPasswordResetToken   = errors.New("password reset token invalid")
)

type (
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required=Refresh token is required"`
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
