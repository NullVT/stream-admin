package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/rs/zerolog/log"
)

type BanUserRequest struct {
	UserID    string  `json:"user_id"`
	Permanent bool    `json:"permanent"`
	Duration  *uint   `json:"duration,omitempty"`
	Reason    *string `json:"reason,omitempty"`
}

type TwitchBanUserRequestData struct {
	UserID   string  `json:"user_id"`
	Duration *uint   `json:"duration,omitempty"`
	Reason   *string `json:"reason,omitempty"`
}
type TwitchBanUserRequest struct {
	Data TwitchBanUserRequestData `json:"data"`
}

func (h *Handler) TwitchBanUser(ctx echo.Context) error {
	// unmarshal request
	body := new(BanUserRequest)
	if err := ctx.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		return echo.NewHTTPError(500, err.Error())
	}

	// get twitch auth
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		log.Error().Err(err).Msg("failed to get Twitch auth")
		return echo.NewHTTPError(500)
	}

	// build request
	reqURL, _ := url.Parse("https://api.twitch.tv/helix/moderation/bans")
	reqQuery := reqURL.Query()
	reqQuery.Add("broadcaster_id", twitchAuth.UserID)
	reqQuery.Add("moderator_id", twitchAuth.UserID)
	reqURL.RawQuery = reqQuery.Encode()
	requestBody := TwitchBanUserRequest{
		Data: TwitchBanUserRequestData{
			UserID:   body.UserID,
			Duration: body.Duration,
			Reason:   body.Reason,
		},
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal request body")
		return echo.NewHTTPError(500, "Failed to create request")
	}
	req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Error().Err(err).Msg("failed to init request")
		return echo.NewHTTPError(500)
	}
	req.Header.Set("Client-Id", twitchAuth.ClientID)
	req.Header.Set("Authorization", twitchAuth.Bearer())
	req.Header.Set("Content-Type", "application/json")

	// send req
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to update Twitch channel info")
		return echo.NewHTTPError(500, "twitch error")
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 200 {
		log.Error().Any("responseCode", res.StatusCode).Msg("failed to update Twitch channel info")
		return echo.NewHTTPError(500, "twitch error")
	}

	return ctx.JSON(204, nil)
}
