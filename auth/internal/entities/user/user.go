package user

import (
	"mts/auth-service/internal/entities/response"
)

type UserRequest struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type UserResponse struct {
	response.Response

	Id    string `json:"id,omitempty"`
	Login string `json:"login"`
}

func NewUserIdResponse(id, login string) *UserResponse {
	return &UserResponse{Id: id, Login: login}
}

func NewUserResponse(login string) *UserResponse {
	return &UserResponse{Login: login}
}
