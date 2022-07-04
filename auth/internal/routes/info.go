package routes

import (
	"net/http"

	"github.com/go-chi/render"

	"mts/auth-service/internal/entities/user"
	"mts/auth-service/internal/middleware"
)

func Info(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserCtx).(*user.UserResponse)
	render.Render(w, r, user)
}
