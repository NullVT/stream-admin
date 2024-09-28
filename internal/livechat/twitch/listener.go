package twitch

import (
	"context"
	"time"

	"github.com/nullvt/stream-admin/internal/livechat"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

func StartListener(ctx context.Context, msgChan chan livechat.Message, authConfig AuthConfig) <-chan livechat.Message {
	go func() {
		defer close(msgChan)
		var sessionID string
		var subscriptions map[string]string = map[string]string{}

		// connect to twitch
		url := "wss://eventsub.wss.twitch.tv/ws"
		conn, err := websocket.Dial(url, "", "http://localhost/")
		if err != nil {
			log.Error().Err(err).Msg("failed to dial Twitch WebSocket")
			return
		}
		defer conn.Close()

		// handle message type
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var message string
				if err := websocket.Message.Receive(conn, &message); err != nil {
					log.Error().Err(err).Msg("error receiving TwitchWS message")
					return
				}
				log.Debug().Any("websocket msg", message).Msg("Twitch WS Message received")

				// parse the message
				parsedMsg, err := parseTwitchWebsocketMessage([]byte(message))
				if err != nil || parsedMsg == nil {
					log.Error().Err(err).Msg("failed to parse TwitchWS message")
					return
				}

				// handle welcome message
				if parsedMsg.SessionWelcome != nil {
					sessionID = parsedMsg.SessionWelcome.Payload.Session.ID
					chatSubType := "channel.chat.message"
					subscriptionID, err := Subscribe(authConfig, sessionID, chatSubType)
					if err != nil {
						log.Error().Err(err).Msg("failed to subscribe to chat messages")
					}
					subscriptions[chatSubType] = subscriptionID
				}

				// handle chat message
				if parsedMsg.Chat != nil {
					msgChan <- livechat.Message{
						Platform:    "twitch",
						ID:          parsedMsg.Chat.Metadata.MessageID,
						Body:        parsedMsg.Chat.Payload.Event.Message.Text,
						ReceivedAt:  time.Now().UTC(),
						PublishedAt: parsedMsg.Chat.Metadata.MessageTimestamp,
						Sender: livechat.User{
							ID:            parsedMsg.Chat.Payload.Event.ChatterUserID,
							Name:          parsedMsg.Chat.Payload.Event.ChatterUserName,
							Broadcaster:   parsedMsg.Chat.HasBadge("broadcaster"),
							Moderator:     parsedMsg.Chat.HasBadge("moderator"),
							TwitchVIP:     parsedMsg.Chat.HasBadge("vip"),
							YouTubeMember: false,
						},
					}
				}
			}
		}
	}()

	return msgChan
}
