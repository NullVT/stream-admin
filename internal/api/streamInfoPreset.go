package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/rs/zerolog/log"
)

func (h *Handler) StreamInfoPresetGet(ctx echo.Context) error {
	return ctx.JSON(200, config.Cfg.StreamInfoPresets)
}

func (h *Handler) StreamInfoPresetPost(ctx echo.Context) error {
	// read body
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		return echo.NewHTTPError(500, "failed to read request body")
	}
	var preset config.StreamInfoPreset
	if err := json.Unmarshal(body, &preset); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		return echo.NewHTTPError(500, "failed to unmarshal request body")
	}

	// validate body
	if preset.ID != "" {
		return echo.NewHTTPError(400, "cannot create StreamInfoPreset with manually set ID")
	}
	preset.ID = uuid.NewString()

	// save
	newPresets := append(config.Cfg.StreamInfoPresets, preset)
	if err := config.SetConfigValue("streamInfoPresets", newPresets); err != nil {
		log.Error().Err(err).Msg("failed to persist StreamInfoPresets")
		return echo.NewHTTPError(500, "failed to save presets")
	}

	return ctx.JSON(200, newPresets)
}

func (h *Handler) StreamInfoPresetPut(ctx echo.Context) error {
	// find preset in config
	presetID := ctx.Param("id")
	if presetID == "" {
		return echo.NewHTTPError(400, "invalid preset id")
	}

	// Read existing preset
	oldPresetIndex := -1
	for i, preset := range config.Cfg.StreamInfoPresets {
		if presetID == preset.ID {
			oldPresetIndex = i
			break
		}
	}
	if oldPresetIndex == -1 {
		return echo.NewHTTPError(404, "preset not found")
	}

	// read body
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		return echo.NewHTTPError(500, "failed to read request body")
	}

	// unmarshal body into a preset
	var preset config.StreamInfoPreset
	if err := json.Unmarshal(body, &preset); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		return echo.NewHTTPError(400, "failed to unmarshal request body")
	}

	// validate: preset ID must match the existing one
	if preset.ID != presetID {
		return echo.NewHTTPError(400, "you cannot change the ID of a preset")
	}

	// Update the existing preset in the slice
	config.Cfg.StreamInfoPresets[oldPresetIndex] = preset

	// Persist updated presets
	if err := config.SetConfigValue("streamInfoPresets", config.Cfg.StreamInfoPresets); err != nil {
		log.Error().Err(err).Msg("failed to persist StreamInfoPresets")
		return echo.NewHTTPError(500, "failed to save presets")
	}

	return ctx.JSON(200, config.Cfg.StreamInfoPresets)
}

func (h *Handler) StreamInfoPresetDelete(ctx echo.Context) error {
	// find preset in config
	presetID := ctx.Param("id")
	if presetID == "" {
		return echo.NewHTTPError(400, "invalid preset id")
	}

	// Find the index of the preset to delete
	oldPresetIndex := -1
	for i, preset := range config.Cfg.StreamInfoPresets {
		if presetID == preset.ID {
			oldPresetIndex = i
			break
		}
	}

	// If no matching preset was found, return 404
	if oldPresetIndex == -1 {
		return echo.NewHTTPError(404, "preset not found")
	}

	// Remove the preset from the slice
	config.Cfg.StreamInfoPresets = append(
		config.Cfg.StreamInfoPresets[:oldPresetIndex],
		config.Cfg.StreamInfoPresets[oldPresetIndex+1:]...,
	)

	// Persist the updated list of presets
	if err := config.SetConfigValue("streamInfoPresets", config.Cfg.StreamInfoPresets); err != nil {
		log.Error().Err(err).Msg("failed to persist StreamInfoPresets")
		return echo.NewHTTPError(500, "failed to save presets")
	}

	return ctx.JSON(200, config.Cfg.StreamInfoPresets)
}

func (h *Handler) StreamInfoPresetApply(ctx echo.Context) error {
	// find preset in config
	presetID := ctx.Param("id")
	if presetID == "" {
		return echo.NewHTTPError(400, "invalid preset id")
	}

	// Read existing preset
	var oldPresetIndex = -1
	var preset config.StreamInfoPreset
	for i, p := range config.Cfg.StreamInfoPresets {
		if presetID == p.ID {
			oldPresetIndex = i
			preset = p
			break
		}
	}
	if oldPresetIndex == -1 {
		return echo.NewHTTPError(404, "preset not found")
	}

	// get twitch auth
	twitchAuth, err := helpers.GetTwitchAuth()
	if err != nil {
		log.Error().Err(err).Msg("failed to get Twitch auth")
		return echo.NewHTTPError(500)
	}

	// build request
	reqURL, _ := url.Parse("https://api.twitch.tv/helix/channels")
	reqQuery := reqURL.Query()
	reqQuery.Add("broadcaster_id", twitchAuth.UserID)
	reqURL.RawQuery = reqQuery.Encode()
	requestBody := TwitchModifyChannelInformationRequest{
		GameID: preset.Category.ID,
		Title:  preset.Title,
		Tags:   preset.Tags,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal request body")
		return echo.NewHTTPError(500, "Failed to create request")
	}
	req, err := http.NewRequest("PATCH", reqURL.String(), bytes.NewBuffer(bodyBytes))
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
	if res.StatusCode != 204 {
		log.Error().Any("responseCode", res.StatusCode).Msg("failed to update Twitch channel info")
		return echo.NewHTTPError(500, "twitch error")
	}

	return ctx.JSON(200, config.Cfg.StreamInfoPresets)
}
