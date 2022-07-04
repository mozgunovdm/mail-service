package message

import (
	"mts/auth-service/internal/entities/response"
)

type MessageResponse struct {
	response.Response

	Message string `json:"message"`
}

func NewMessageResponse(message string) *MessageResponse {
	return &MessageResponse{Message: message}
}
