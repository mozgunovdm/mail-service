package routes

import (
	"mts/auth-service/internal/entities/error"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie := http.Cookie{Name: AccessToken, Value: "", Expires: time.Now()}
	http.SetCookie(w, &accessTokenCookie)

	refreshTokenCookie := http.Cookie{Name: RefreshToken, Value: "", Expires: time.Now()}
	http.SetCookie(w, &refreshTokenCookie)

	redirectUri := r.URL.Query().Get(RedirectUri)
	if redirectUri != "" {
		ref := r.Header.Get("Referer")
		http.Redirect(w, r, ref+redirectUri, http.StatusSeeOther)
	}

	newError := error.NewErrorResponse("")
	render.Render(w, r, newError)
}
