package api

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetEmote(ctx echo.Context) error {
	// lookup ID in cache index
	emoteID := ctx.Param("id")
	emote := h.emotesCache.FindByID(emoteID)
	if emote == nil {
		return echo.NewHTTPError(404)
	}

	// load file
	if _, err := os.Stat(emote.FilePath); os.IsNotExist(err) {
		log.Error().Err(err).Msg("emote file does not exist")
		return echo.NewHTTPError(404, "file not found")
	}

	// read the file content
	return ctx.File(emote.FilePath)
}

func (h *Handler) EmoteWhitelistGet(ctx echo.Context) error {
	whitelist := config.Cfg.EmotesWhitelist
	return ctx.JSON(200, whitelist)
}

func (h *Handler) EmoteWhitelistPost(ctx echo.Context) error {
	// read body
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		return echo.NewHTTPError(500, "failed to read request body")
	}
	channelName := ""
	if err := json.Unmarshal(body, &channelName); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		return echo.NewHTTPError(500, "failed to unmarshal request body")
	}

	// find channel ID
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		return echo.NewHTTPError(500, "failed to load twitch auth")
	}
	users, err := twitch.GetUsers(twitchAuth, []string{strings.ToLower(channelName)})
	if err != nil {
		log.Error().Err(err).Msg("failed to get twitch users")
		return echo.NewHTTPError(500, "failed to get twitch users")
	}
	if len(users) != 1 {
		log.Error().Msg("twitch GetUsers returned unexpected length")
	}

	// persist to config
	whitelist := config.Cfg.EmotesWhitelist
	whitelist[users[0].ID] = users[0].DisplayName
	if err := config.SetConfigValue("emotesWhitelist", whitelist); err != nil {
		log.Error().Err(err).Msg("failed to persist EmotesWhitelist")
		return echo.NewHTTPError(500, "failed to update whitelist")
	}

	return ctx.JSON(200, whitelist)
}

func (h *Handler) EmoteWhitelistDelete(ctx echo.Context) error {
	// read body
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		return echo.NewHTTPError(500, "failed to read request body")
	}
	var channelID string
	if err := json.Unmarshal(body, &channelID); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		return echo.NewHTTPError(500, "failed to unmarshal request body")
	}

	// remove channel and persist
	whitelist := config.Cfg.EmotesWhitelist
	delete(whitelist, string(channelID))
	if err := config.SetConfigValue("emotesWhitelist", whitelist); err != nil {
		log.Error().Err(err).Msg("failed to persist EmotesWhitelist")
		return echo.NewHTTPError(500, "failed to update whitelist")
	}

	return ctx.JSON(200, whitelist)
}
