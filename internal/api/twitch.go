package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/nullvt/stream-admin/internal/secrets"
	"github.com/rs/zerolog/log"
)

const (
	twitchGqlUrl = "https://gql.twitch.tv/gql"
)

type TwitchUpdateChatSettingsRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Extensions    map[string]interface{} `json:"extensions"`
}

type TwitchUpdateChatSettingsResponse struct {
	Data struct {
		UpdateChatSettings struct {
			ChatSettings struct {
				HideLinks bool `json:"hideLinks"`
			} `json:"chatSettings"`
		} `json:"updateChatSettings"`
	} `json:"data"`
}

type TwitchModifyChannelInformationRequest struct {
	GameID string   `json:"game_id"`
	Title  string   `json:"title"`
	Tags   []string `json:"tags"`
}

type LinkFilteringRequest struct {
	Enabled bool `json:"enabled"`
}

func (h *Handler) TwitchLinkFiltering(ctx echo.Context) error {
	// unmarshal request
	body := new(LinkFilteringRequest)
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

	// get twitch auth
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		log.Error().Err(err).Msg("failed to get Twitch auth")
		return echo.NewHTTPError(500)
	}

	// create GQL query
	requestBody := []TwitchUpdateChatSettingsRequest{
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
	req.Header.Set("Authorization", twitchAuth.Bearer())
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

	var response TwitchUpdateChatSettingsResponse
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

func (h *Handler) TwitchCategorySearch(ctx echo.Context) error {
	perPage := "20"
	query := ctx.Request().URL.Query().Get("query")
	if query == "" {
		return echo.NewHTTPError(400, "query param required")
	}

	// get twitch auth
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		log.Error().Err(err).Msg("failed to get Twitch auth")
		return echo.NewHTTPError(500)
	}

	// set URL and query
	reqURL, _ := url.Parse("https://api.twitch.tv/helix/search/categories")
	reqQuery := reqURL.Query()
	reqQuery.Add("query", query)
	reqQuery.Add("first", perPage) // WTF Twitch
	reqURL.RawQuery = reqQuery.Encode()

	// create http req
	req, err := http.NewRequest("GET", reqURL.String(), nil)
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
		log.Error().Err(err).Msg("failed to search Twitch categories")
		return echo.NewHTTPError(500, "twitch error")
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 200 {
		log.Error().Any("responseCode", res.StatusCode).Msg("failed to search Twitch categories")
		return echo.NewHTTPError(500, "twitch error")
	}

	// parse the response
	var resBody struct {
		Data []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			BoxArtUrl string `json:"box_art_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal response")
		return echo.NewHTTPError(500)
	}

	return ctx.JSON(200, resBody.Data)
}

func (h *Handler) TwitchGetStreamInfo(ctx echo.Context) error {
	// get twitch auth
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		log.Error().Err(err).Msg("failed to get Twitch auth")
		return echo.NewHTTPError(500)
	}

	// set URL and query
	reqURL, _ := url.Parse("https://api.twitch.tv/helix/channels")
	reqQuery := reqURL.Query()
	reqQuery.Add("broadcaster_id", twitchAuth.BroadcasterID)
	reqURL.RawQuery = reqQuery.Encode()

	// create http req
	req, err := http.NewRequest("GET", reqURL.String(), nil)
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
		log.Error().Err(err).Msg("failed to get channel information")
		return echo.NewHTTPError(500, "twitch error")
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 200 {
		log.Error().Any("responseCode", res.StatusCode).Msg("failed to get Twitch channel information")
		return echo.NewHTTPError(500, "twitch error")
	}

	// parse the response
	var resBody struct {
		Data []struct {
			Title    string   `json:"title"`
			GameName string   `json:"game_name"`
			GameID   string   `json:"game_id"`
			Tags     []string `json:"tags"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal response")
		return echo.NewHTTPError(500)
	}

	if len(resBody.Data) != 1 {
		log.Error().Any("dataLength", len(resBody.Data)).Msg("unexpected number of results for Twitch get stream info")
		return echo.NewHTTPError(500, "malformed Twitch GetStreamInfo")
	}

	return ctx.JSON(200, resBody.Data[0])
}
