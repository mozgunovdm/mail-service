package app

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5/middleware"

	"mts/auth-service/internal/config"
	mdw "mts/auth-service/internal/middleware"
	"mts/auth-service/internal/routes"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func InitApp(ctx context.Context, log zerolog.Logger) {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	//r.Use(middleware.Logger)

	c := alice.New()
	c = c.Append(hlog.NewHandler(log))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	r.Route("/auth/v1/i", func(r chi.Router) {
		r.Use(mdw.ValidateToken)
		r.Method("GET", "/", c.Then(routes.Handler(routes.Info)))
	})
	r.Mount("/debug/", middleware.Profiler())
	r.Method("POST", "/auth/v1/login", c.Then(routes.Handler(routes.Login)))
	r.Method("POST", "/auth/v1/logout", c.Then(routes.Handler(routes.Logout)))
	r.Method("GET", "/profile", c.Then(routes.Handler(routes.Profile)))

	cfg := ctx.Value("config").(config.Config)
	os.Setenv("ACCESS_SECRET", cfg.Auth.Token.AccessSecret)
	os.Setenv("REFRESH_SECRET", cfg.Auth.Token.RefreshSecret)

	port := cfg.Auth.Port
	if len(port) < 2 {
		port = "3000"
	}

	var sb strings.Builder
	sb.WriteString(cfg.Auth.Host)
	sb.WriteString(":")
	sb.WriteString(port)
	log.Debug().Msg(sb.String())
	srv := http.Server{
		Addr:    sb.String(),
		Handler: r,
	}

	servChan := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Timestamp().Err(err).Msgf("HTTP server ListenAndServe: %w", err)
		} else {
			log.Info().Timestamp().Msg("HTTP server closed")
		}
		close(servChan)
	}()

	select {
	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Error().Timestamp().Err(err).Msgf("HTTP server Shutdown: %v", err)
		}
	case <-servChan:
		log.Debug().Timestamp().Msg("HTTP server Shutdown by channel")
	}
}
