package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/rs/zerolog"
)

type Handler struct {
	msgChan     chan livechat.Message
	emotesCache *livechat.EmoteCache
}

func Start(msgChan chan livechat.Message, emc *livechat.EmoteCache) (*echo.Echo, error) {
	// Setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	// Use zerolog for request logging
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("Status", v.Status).
				Msg("request")
			return nil
		},
	}))

	// Fuck CORS with a rake
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// API routes
	apiGroup := e.Group("/api")
	handler := &Handler{
		msgChan:     msgChan,
		emotesCache: emc,
	}

	// messages
	apiGroup.GET("/messages", handler.MessageWebsocket)
	apiGroup.DELETE("/messages/:id", handler.MessageDelete)

	// emotes
	apiGroup.GET("/emotes/:id", handler.GetEmote)
	apiGroup.GET("/emotes/whitelist", handler.EmoteWhitelistGet)
	apiGroup.POST("/emotes/whitelist", handler.EmoteWhitelistPost)
	apiGroup.DELETE("/emotes/whitelist", handler.EmoteWhitelistDelete)

	// stream info
	apiGroup.GET("/stream-info", handler.TwitchGetStreamInfo)
	apiGroup.GET("/stream-info-presets", handler.StreamInfoPresetGet)
	apiGroup.POST("/stream-info-presets", handler.StreamInfoPresetPost)
	apiGroup.PUT("/stream-info-presets/:id", handler.StreamInfoPresetPut)
	apiGroup.DELETE("/stream-info-presets/:id", handler.StreamInfoPresetDelete)
	apiGroup.POST("/stream-info-presets/:id/apply", handler.StreamInfoPresetApply)

	// Twitch routes
	apiGroup.GET("/auth/twitch", handler.TwitchLogin)
	apiGroup.POST("/auth/twitch", handler.TwitchCallback)
	apiGroup.DELETE("/auth/twitch", handler.TwitchLogout)
	apiGroup.GET("/auth/twitch/valid", handler.TwitchValidateAuth)
	apiGroup.POST("/twitch/link-filtering", handler.TwitchLinkFiltering)
	apiGroup.GET("/twitch/categories", handler.TwitchCategorySearch)
	apiGroup.POST("/twitch/ban-user", handler.TwitchBanUser)

	// Start server in a goroutine
	go func() {
		address := fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// start background tasks
	go pruneOldMessages()

	return e, nil
}
