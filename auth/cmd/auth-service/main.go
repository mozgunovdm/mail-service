package main

import (
	"context"
	"mts/auth-service/internal/app"
	"mts/auth-service/internal/config"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

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

	ctx = context.WithValue(ctx, "config", *cfg)
	log.Info().Timestamp().Msg("Start service")

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Debug().Timestamp().Msg("Get sigterm")
		cancelFunc()
		close(idleConnsClosed)
	}()

	app.InitApp(ctx, log)

	<-idleConnsClosed
}
