package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nullvt/stream-admin/internal/config"
	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	msgChan     chan livechat.Message
	wsClients   map[*websocket.Conn]bool // Track active WebSocket clients
	wsClientsMu sync.Mutex               // Ensure thread-safe access to clients
	emotesCache *livechat.EmoteCache
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// You can customize this to allow different origins
		return true
	},
}

func (h *Handler) MessageWebsocket(ctx echo.Context) error {
	// Upgrade the HTTP connection to a WebSocket
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		ctx.Logger().Errorf("failed to upgrade connection: %v", err)
		return err
	}
	defer ws.Close()

	// Add new WebSocket connection to the clients map
	h.wsClientsMu.Lock()
	h.wsClients[ws] = true
	h.wsClientsMu.Unlock()
	log.Info().Any("ws", ws).Msg("WS client connected")

	// Remove WebSocket connection when the function returns
	defer func() {
		h.wsClientsMu.Lock()
		delete(h.wsClients, ws)
		h.wsClientsMu.Unlock()
		log.Info().Any("ws", ws).Msg("WS client removed")
	}()

	// Create a go-routine to send messages from the channel to the WebSocket
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range h.msgChan {
			// Convert the message to JSON
			msgJson, err := json.Marshal(msg)
			if err != nil {
				log.Printf("failed to marshal message to JSON: %v\n", err)
				continue
			}

			// Broadcast the message to all connected WebSocket clients
			h.wsClientsMu.Lock()
			for client := range h.wsClients {
				err := client.WriteMessage(websocket.TextMessage, msgJson)
				if err != nil {
					log.Printf("failed to send message to WebSocket client: %v\n", err)
					client.Close()
					delete(h.wsClients, client) // Remove clients that fail to receive the message
				}
			}
			h.wsClientsMu.Unlock()
		}
	}()

	// Ensure the channel goroutine is closed properly
	<-done
	return nil
}

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

func Start(msgChan chan livechat.Message, emc *livechat.EmoteCache) (*echo.Echo, error) {
	// Setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
		wsClients:   make(map[*websocket.Conn]bool),
		emotesCache: emc,
	}
	apiGroup.GET("/messages", handler.MessageWebsocket)
	apiGroup.GET("/emotes/:id", handler.GetEmote)

	// Twitch routes
	apiGroup.GET("/auth/twitch", handler.TwitchLogin)
	apiGroup.POST("/auth/twitch", handler.TwitchCallback)
	apiGroup.GET("/auth/twitch/valid", handler.TwitchValidateAuth)
	apiGroup.POST("/twitch/link-filtering", handler.TwitchLinkFiltering)

	// Start server in a goroutine
	go func() {
		address := fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	return e, nil
}
