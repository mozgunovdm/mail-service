package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/hlog"

	"mts/auth-service/internal/db/repository"
	"mts/auth-service/internal/entities/error"
	"mts/auth-service/internal/entities/user"
	"mts/auth-service/internal/service"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var bodyUser user.UserRequest
	log := hlog.FromRequest(r)
	err := json.NewDecoder(r.Body).Decode(&bodyUser)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg(ErrBadRequest.Error())
		w.WriteHeader(http.StatusInternalServerError)
		newError := error.NewErrorResponse("")
		render.Render(w, r, newError)
		return
	}
	defer r.Body.Close()

	dbUser, ok := repository.Repository[bodyUser.Login]
	if !ok {
		log.Error().Err(err).Timestamp().Msg(ErrWrongLogin.Error())
		newError := error.NewErrorResponse(ErrWrongLogin.Error())
		render.Render(w, r, newError)
		return
	}

	isPasswordValid := service.CheckPasswordHash(bodyUser.Password, dbUser.Password)
	if !isPasswordValid {
		log.Error().Err(err).Timestamp().Msg(ErrWrongPassword.Error())
		w.WriteHeader(http.StatusForbidden)
		newError := error.NewErrorResponse(ErrWrongPassword.Error())
		render.Render(w, r, newError)
		return
	}

	resUser := user.NewUserIdResponse(dbUser.Id, dbUser.Login)

	refreshToken, err := service.GenerateRefreshToken(resUser)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg(ErrTokenExpired.Error())
		newError := error.NewErrorResponse(ErrTokenExpired.Error())
		render.Render(w, r, newError)
		return
	}
	refreshTokenCookie := http.Cookie{Name: RefreshToken, Value: refreshToken, Expires: time.Now().Add(time.Hour), HttpOnly: true}
	http.SetCookie(w, &refreshTokenCookie)

	accessToken, _ := service.GenerateAccessToken(resUser)
	accessTokenCookie := http.Cookie{Name: AccessToken, Value: accessToken, Expires: time.Now().Add(time.Minute), HttpOnly: true}
	http.SetCookie(w, &accessTokenCookie)

	redirectUri := r.URL.Query().Get(RedirectUri)
	if redirectUri != "" {
		ref := r.Header.Get("Referer")
		log.Info().Timestamp().Msgf("Redirect to %s", ref+redirectUri)
		http.Redirect(w, r, ref+redirectUri, http.StatusSeeOther)
	}

	render.Render(w, r, resUser)
}
