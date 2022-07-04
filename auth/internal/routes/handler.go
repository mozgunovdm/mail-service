package routes

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

var (
	ErrBadRequest    = errors.New("bad request")
	ErrWrongLogin    = errors.New("wrong login")
	ErrWrongPassword = errors.New("wrong password")
	ErrTokenExpired  = errors.New("token expired")
)

const (
	AccessToken  = "accessToken"
	RefreshToken = "refreshToken"
	RedirectUri  = "redirect_uri"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().Timestamp().
		Str("Url", r.URL.String()).
		Str("Host", r.Host).
		Str("Method", r.Method).
		Msg("Getting request:")
	h(w, r)
}
