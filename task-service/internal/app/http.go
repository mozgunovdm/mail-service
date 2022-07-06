package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5/middleware"

	"mts/task-service/internal/config"
	"mts/task-service/internal/db"
	"mts/task-service/internal/router"

	"github.com/rs/zerolog"
)

type RestService struct {
	Log *zerolog.Logger
	Cfg *config.Config
	Db  db.IDatabase
}

func (s RestService) Start() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}))

	r.Use(middleware.Logger)

	r.Mount("/", router.Router{
		Db:  s.Db,
		Log: s.Log,
		Cfg: s.Cfg,
	}.Routes())

	restAddr := fmt.Sprintf("%s:%v", s.Cfg.Rest.Host, s.Cfg.Rest.Port)
	s.Log.Debug().Msg(restAddr)
	srv := http.Server{
		Addr:    restAddr,
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.Log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			s.Log.Fatal().Err(err)
		}
		serverStopCtx()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Log.Fatal().Timestamp().Err(err).Msgf("HTTP server ListenAndServe: %w", err)
	} else {
		s.Log.Info().Timestamp().Msg("HTTP server closed")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
