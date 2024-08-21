package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/nullvt/stream-admin/internal/secrets"
	"github.com/rs/zerolog/log"
)

const (
	twitchGqlUrl = "https://gql.twitch.tv/gql"
)

type UpdateChatSettingsRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Extensions    map[string]interface{} `json:"extensions"`
}

type UpdateChatSettingsResponse struct {
	Data struct {
		UpdateChatSettings struct {
			ChatSettings struct {
				HideLinks bool `json:"hideLinks"`
			} `json:"chatSettings"`
		} `json:"updateChatSettings"`
	} `json:"data"`
}

type TwitchLinkFilteringRequest struct {
	Enabled bool `json:"enabled"`
}

func (h *Handler) TwitchLinkFiltering(ctx echo.Context) error {
	// unmarshal request
	body := new(TwitchLinkFilteringRequest)
	if err := ctx.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		return echo.NewHTTPError(500, err.Error())
	}

	// get channel ID
	userInfoJson, err := secrets.Get("twitch_user")
	if err != nil {
		log.Error().Err(err).Msg("failed to get Twitch UserInfo from secrets")
		return echo.NewHTTPError(500, "Failed to load UserInfo")
	}
	userInfo := &twitch.User{}
	if err := userInfo.UnmarshalString(userInfoJson); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal Twitch UserInfo")
		return echo.NewHTTPError(500, "Failed to load UserInfo")
	}

	// get twitch secret
	authToken, err := secrets.Get("twitch_token")
	if err != nil || authToken == "" {
		log.Error().Err(err).Msg("failed to get twitch token from keyring")
		return echo.NewHTTPError(500, "Failed to load Twitch auth token")
	}

	// create GQL query
	requestBody := []UpdateChatSettingsRequest{
		{
			OperationName: "UpdateChatSettings",
			Variables: map[string]interface{}{
				"input": map[string]interface{}{
					"channelID": userInfo.ID,
					"hideLinks": body.Enabled,
				},
			},
			Extensions: map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"version":    1,
					"sha256Hash": "6d8b11f4e29f87be5e2397dd54b2df669e9a5aacd831252d88b7b7a6616dc170",
				},
			},
		},
	}

	// create request
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal request body")
		return echo.NewHTTPError(500, "Failed to create request")
	}
	req, err := http.NewRequest("POST", twitchGqlUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Error().Err(err).Msg("failed to init request")
		return echo.NewHTTPError(500, "Failed to create request")
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Client-ID", config.Cfg.Twitch.ClientID)
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("GQL request failed")
		return echo.NewHTTPError(500, "Failed query Twitch GQL")
	}
	defer resp.Body.Close()

	var response UpdateChatSettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal GQL response")
		return echo.NewHTTPError(500, "Twitch GQL request failed")
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Any("responseCode", resp.StatusCode).Any("body", response).Msg("GQL request failed")
		return echo.NewHTTPError(500, "Failed query Twitch GQL")
	}

	if response.Data.UpdateChatSettings.ChatSettings.HideLinks {
		return ctx.NoContent(200)
	} else {
		log.Error().Any("response", response.Data).Msg("unexpected response from Twitch GQL")
		return echo.NewHTTPError(500, "Twitch GQL request failed")
	}
}
