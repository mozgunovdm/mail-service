package router

import (
	"mts/task-service/internal/config"
	"mts/task-service/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
)

const (
	AccessToken  = "accessToken"
	RefreshToken = "refreshToken"
	RedirectUri  = "redirect_uri"
)

type Router struct {
	Db  db.IDatabase
	Log *zerolog.Logger
	Cfg *config.Config
}

func NewRouter(d db.IDatabase, log *zerolog.Logger, cfg *config.Config) *Router {
	return &Router{
		Db:  d,
		Log: log,
		Cfg: cfg,
	}
}

// Routes creates a REST router for the task service
func (rs Router) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	r.Get("/healthz", rs.healthzHandler)
	r.Mount("/", rs.taskHandlers())

	return r
}
