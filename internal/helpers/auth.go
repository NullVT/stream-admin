package helpers

import (
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/nullvt/stream-admin/internal/secrets"
	"github.com/rs/zerolog/log"
)

func GetTwitchAuth() (twitch.AuthConfig, error) {
	twitchToken, err := secrets.Get("twitch_token")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load twitch_token")
		return twitch.AuthConfig{}, err
	}

	return twitch.AuthConfig{
		ClientID:  config.Cfg.Twitch.ClientID,
		AuthToken: twitchToken,
		// TODO: replace with user from secrets
		UserID:        "1113117444",
		BroadcasterID: "1113117444",
	}, nil
}
