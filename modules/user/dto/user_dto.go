package dto

import (
	"errors"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER      = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER      = "failed get list user"
	MESSAGE_FAILED_TOKEN_NOT_VALID    = "token not valid"
	MESSAGE_FAILED_TOKEN_NOT_FOUND    = "token not found"
	MESSAGE_FAILED_GET_USER           = "failed get user"
	MESSAGE_FAILED_LOGIN              = "failed login"
	MESSAGE_FAILED_UPDATE_USER        = "failed update user"
	MESSAGE_FAILED_DELETE_USER        = "failed delete user"
	MESSAGE_FAILED_PROSES_REQUEST     = "failed proses request"
	MESSAGE_FAILED_DENIED_ACCESS      = "denied access"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER = "success get list user"
	MESSAGE_SUCCESS_GET_USER      = "success get user"
	MESSAGE_SUCCESS_LOGIN         = "success login"
	MESSAGE_SUCCESS_UPDATE_USER   = "success update user"
	MESSAGE_SUCCESS_DELETE_USER   = "success delete user"
)

var (
	ErrCreateUser        = errors.New("failed to create user")
	ErrGetUserById       = errors.New("failed to get user by id")
	ErrGetUserByName     = errors.New("failed to get user by name")
	ErrNameAlreadyExists = errors.New("user already exist")
	ErrUpdateUser        = errors.New("failed to update user")
	ErrUserNotFound      = errors.New("user not found")
	ErrDeleteUser        = errors.New("failed to delete user")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrTokenExpired      = errors.New("token expired")
)

type (
	UserCreateRequest struct {
		Name     string `json:"name" form:"name" binding:"required,min=2,max=100"`
		Password string `json:"password" form:"password" binding:"required,min=4"`
	}

	UserResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	UserUpdateRequest struct {
		Name string `json:"name" form:"name" binding:"omitempty,min=2,max=100"`
	}

	UserUpdateResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	UserLoginRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
)
