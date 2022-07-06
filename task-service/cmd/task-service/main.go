package main

import (
	"context"
	"mts/task-service/internal/app"
	"mts/task-service/internal/config"
	"mts/task-service/internal/db/mongo"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().
		Str("Service", "Auth").
		Logger()

	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal().Err(err).Timestamp().Msg("Error config path")
	}
	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Timestamp().Msg("Error open config file")
	}

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mongo.NewMongo(initCtx, cfg, &log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init database")
		return
	}
	defer db.Close(initCtx)

	app.RestService{
		Log: &log,
		Cfg: cfg,
		Db:  db,
	}.Start()
}
