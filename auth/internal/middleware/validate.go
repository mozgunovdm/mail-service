package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/render"

	"mts/auth-service/internal/entities/error"
	"mts/auth-service/internal/entities/user"
	"mts/auth-service/internal/service"
)

var (
	ErrUserNotAuth = error.NewErrorResponse("user not authorised")
)

type userCtx string

const UserCtx userCtx = "user"

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refreshToken")
		if err != nil {
			render.Render(w, r, ErrUserNotAuth)
			return
		}

		claims := service.ParseRefreshToken(refreshToken.Value)
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			render.Render(w, r, ErrUserNotAuth)
			return
		}

		user := user.NewUserResponse(claims.Login)
		ctx := context.WithValue(r.Context(), UserCtx, user)

		_, err = r.Cookie("accessToken")
		if err != nil {
			accessToken, _ := service.GenerateAccessToken(user)
			accessTokenCookie := http.Cookie{Name: "accessToken", Value: accessToken, Expires: time.Now().Add(time.Minute)}
			http.SetCookie(w, &accessTokenCookie)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
