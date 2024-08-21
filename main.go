package main

import (
	"os"

	"github.com/nullvt/stream-admin/internal/api"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	// start API server
	_, err := api.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start API server")
		os.Exit(1)
	}
}
