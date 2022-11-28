package main

import (
	"github.com/lus/kratos-readonly-traits/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Set up the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.IsDevEnv() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Warn().Msg("The service was started in development mode. Please change the 'ENVIRONMENT' variable to 'prod' in production!")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Warn().Msg("An invalid log level was provided. Using the 'info' fallback value.")
			logLevel = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(logLevel)
	}

	// Wait for a Ctrl-C signal
	log.Info().Msg("The application has been started. To stop it press Ctrl-C.")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
