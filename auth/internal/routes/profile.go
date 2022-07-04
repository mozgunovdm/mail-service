package routes

import (
	"net/http"
	"runtime"

	"mts/auth-service/internal/entities/message"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/hlog"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	var msg string
	if runtime.MemProfileRate == 0 {
		log.Info().Timestamp().Msg("Profile started")
		runtime.MemProfileRate = 2048
	} else {
		log.Info().Timestamp().Msg("Profile stoped")
		runtime.MemProfileRate = 0
	}
	render.Render(w, r, message.NewMessageResponse(msg))
}
