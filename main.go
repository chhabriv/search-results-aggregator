package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/chhabriv/search-results-aggregator/serverconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO: make log level and format configurable
	// throught properties file or flags
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Info().Interface("system_call", oscall.String()).Msg("Received system call")
		cancel()
	}()

	serverconfig.StartServer(ctx)
}
