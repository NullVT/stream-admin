package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nullvt/stream-admin/internal/helpers"
	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/nullvt/stream-admin/internal/livechat/twitch"
	"github.com/rs/zerolog/log"
)

var (
	wsClients   map[*websocket.Conn]bool // Track active WebSocket clients
	wsClientsMu sync.Mutex
	msgCache    []livechat.Message
	msgCacheMu  sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// You can customize this to allow different origins
		return true
	},
}

func init() {
	msgCache = make([]livechat.Message, 0, 100)
	wsClients = make(map[*websocket.Conn]bool)
}

func findMessageByID(msgID string) *livechat.Message {
	for _, msg := range msgCache {
		if msg.ID == msgID {
			return &msg
		}
	}
	return nil
}

func pruneOldMessages() {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			sixHoursAgo := now.Add(-6 * time.Hour)

			msgCacheMu.Lock()
			var newCache []livechat.Message
			for _, msg := range msgCache {
				if msg.ReceivedAt.After(sixHoursAgo) {
					newCache = append(newCache, msg)
				}
			}
			msgCache = newCache
			msgCacheMu.Unlock()
		}
	}
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
	wsClientsMu.Lock()
	wsClients[ws] = true
	wsClientsMu.Unlock()
	log.Info().Any("ws", ws).Msg("WS client connected")

	// Remove WebSocket connection when the function returns
	defer func() {
		wsClientsMu.Lock()
		delete(wsClients, ws)
		wsClientsMu.Unlock()
		log.Info().Any("ws", ws).Msg("WS client removed")
	}()

	// Broadcast messages to all connected WebSocket clients
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range h.msgChan {
			// add the message to the cache
			msgCacheMu.Lock()
			if len(msgCache) >= 30000 { // Adjust the threshold as needed
				msgCache = msgCache[1:] // Drop the oldest message
			}
			msgCache = append(msgCache, msg)
			msgCacheMu.Unlock()

			// Convert the message to JSON
			msgJson, err := json.Marshal(msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal message to JSON")
				continue
			}

			// Broadcast the message to all connected WebSocket clients
			wsClientsMu.Lock()
			for client := range wsClients {
				err := client.WriteMessage(websocket.TextMessage, msgJson)
				if err != nil {
					log.Error().Err(err).Msg("failed to send message to WebSocket client")
					client.Close()
					delete(wsClients, client) // Remove clients that fail to receive the message
				}
			}
			wsClientsMu.Unlock()
		}
	}()

	// Ensure the channel goroutine is closed properly
	<-done
	return nil
}

func (h *Handler) MessageDelete(ctx echo.Context) error {
	// get id
	msgID := ctx.Param("id")
	if msgID == "" {
		return echo.NewHTTPError(404)
	}

	// get message from cache
	msg := findMessageByID(msgID)
	if msg == nil {
		return echo.NewHTTPError(404, "message not found")
	}

	// delete twitch message
	if string(msg.Platform) == string(livechat.Twitch) {
		twitchAuth, err := helpers.GetTwitchAuth()
		if err != nil {
			return echo.NewHTTPError(500, "failed to load twitch auth")
		}

		if err := twitch.DeleteMessage(twitchAuth, msg.ID); err != nil {
			log.Error().Err(err).Msg("failed to delete twitch chat message")
			return echo.NewHTTPError(500, "failed to delete message")
		}

		return ctx.NoContent(204)
	}

	// should be unreachable
	return echo.NewHTTPError(500, "unreachable")
}
