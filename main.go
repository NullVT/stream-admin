package main

import (
	"context"
	"os"

	"github.com/nullvt/stream-admin/internal/api"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/rs/zerolog/log"
)

func main() {
	// load config
	if err := config.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	// open channel for sub process comms
	msgChan := make(chan livechat.Message)

	// start API server
	_, err := api.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start API server")
		os.Exit(1)
	}

	// start chat listeners
	twitchAuth := twitch.AuthConfig{}
	twitch.StartListener(context.TODO(), msgChan, twitchAuth)
}
