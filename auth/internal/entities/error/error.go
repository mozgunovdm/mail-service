package error

import (
	"mts/auth-service/internal/entities/response"
)

type Err struct {
	Message string `json:"message"`
}

type ErrorStruct struct {
	response.Response

	Error Err `json:"error"`
}

func NewErrorResponse(message string) *ErrorStruct {
	e := &Err{Message: message}
	return &ErrorStruct{Error: *e}
}
