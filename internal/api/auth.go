package api

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/nullvt/stream-admin/internal/secrets"
	"github.com/rs/zerolog/log"
)

func (h *Handler) TwitchLogin(ctx echo.Context) error {
	scopes := []string{"user:bot", "user:read:chat", "moderator:manage:chat_messages"}

	redirectURL := config.Cfg.Server.BaseURL + "/oauth/twitch"
	url, err := twitch.OAuthLogin(config.Cfg.Twitch.ClientID, redirectURL, scopes)
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	return ctx.JSON(200, map[string]string{"url": url})
}

func (h *Handler) TwitchCallback(ctx echo.Context) error {
	// parse the body
	rawBody, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to read body")
	}
	body := string(rawBody)

	// parse the auth token
	authToken, err := twitch.OAuthCallback(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// store the auth token
	if err := secrets.Set("twitch_token", authToken); err != nil {
		log.Error().Err(err).Msg("Failed to persist token")
		return echo.NewHTTPError(500, "Failed to persist token")
	}

	// get extended user info
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		return echo.NewHTTPError(500, "failed to load twitch auth")
	}
	userInfo, err := twitch.GetUsers(twitchAuth, []string{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get UserInfo for authToken")
		return echo.NewHTTPError(500, "Failed to get UserInfo")
	}

	// persist user info
	userJson, err := userInfo[0].MarshalString()
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal UserInfo to string")
		return echo.NewHTTPError(500, "Failed to persist UserInfo")
	}
	if err := secrets.Set("twitch_user", userJson); err != nil {
		log.Error().Err(err).Msg("Failed to persist UserInfo")
		return echo.NewHTTPError(500, "Failed to persist UserInfo")
	}

	return ctx.NoContent(200)
}

func (h *Handler) TwitchValidateAuth(ctx echo.Context) error {
	token, err := secrets.Get("twitch_token")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get twitch token from secrets")
		return echo.NewHTTPError(500, "Failed to get twitch token from secrets")
	}
	if token == "" {
		ctx.JSON(200, map[string]bool{"isValid": false})
	}

	// init validation request
	isValid, err := twitch.OAuthValidateToken(token)
	if err != nil {
		return echo.NewHTTPError(500, err)
	}
	return ctx.JSON(200, map[string]bool{"isValid": isValid})
}
