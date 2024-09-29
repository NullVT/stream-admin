package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nullvt/stream-admin/internal/api"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
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

	// load emotes
	emc := &livechat.EmoteCache{}
	emotesIndexFile := "./emotecache/index.json"
	if err := emc.LoadFromFile(emotesIndexFile); err != nil {
		log.Error().Err(err).Msg("failed to load emotes")
	}

	// Start API server
	server, err := api.Start(msgChan, emc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start API server")
		os.Exit(1)
	}

	// Start chat listeners
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		os.Exit(1)
	}
	twitch.StartListener(context.TODO(), msgChan, twitchAuth, emc)

	// sync emotes
	// TODO: setup proper background task processing
	emotesChannels := helpers.MapKeys(config.Cfg.EmotesWhitelist)
	if err := twitch.SyncEmotes(emc, twitchAuth, emotesChannels); err != nil {
		log.Error().Err(err).Msg("Failed to sync Twitch Emotes")
	}
	emc.SaveToFile(emotesIndexFile)

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
