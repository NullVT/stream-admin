package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nullvt/stream-admin/internal/api"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/nullvt/stream-admin/internal/secrets"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set global zerolog configuration
	zerolog.SetGlobalLevel(zerolog.InfoLevel)                                                // Global log level
	zerolog.TimeFieldFormat = time.RFC3339                                                   // Human-readable timestamps
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}) // Human-readable log format with console writer
	log.Logger = log.With().Caller().Logger()

	// Load config
	if err := config.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	// Open channel for sub process comms
	msgChan := make(chan livechat.Message)

	// Start API server
	server, err := api.Start(msgChan)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start API server")
		os.Exit(1)
	}

	// Start chat listeners
	twitchToken, err := secrets.Get("twitch_token")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load twitch_token")
		os.Exit(1)
	}
	twitchAuth := twitch.AuthConfig{
		ClientID:  config.Cfg.Twitch.ClientID,
		AuthToken: twitchToken,
		// TODO: replace with user from secrets
		UserID:        "1113117444",
		BroadcasterID: "1113117444",
	}
	twitch.StartListener(context.TODO(), msgChan, twitchAuth)

	// Graceful shutdown on SIGINT and SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to gracefully shut down server")
	} else {
		log.Info().Msg("Server shut down gracefully")
	}
}
