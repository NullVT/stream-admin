package twitch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Metadata represents the metadata of a Twitch WebSocket message.
type Metadata struct {
	MessageID           string    `json:"message_id"`
	MessageType         string    `json:"message_type"`
	MessageTimestamp    time.Time `json:"message_timestamp"`
	SubscriptionType    *string   `json:"subscription_type,omitempty"`
	SubscriptionVersion *string   `json:"subscription_version,omitempty"`
}

// Session represents the session details in a session welcome message.
type Session struct {
	ID                      string  `json:"id"`
	Status                  string  `json:"status"`
	ConnectedAt             string  `json:"connected_at"`
	KeepaliveTimeoutSeconds int     `json:"keepalive_timeout_seconds"`
	ReconnectURL            *string `json:"reconnect_url,omitempty"`
	RecoveryURL             *string `json:"recovery_url,omitempty"`
}

// SessionWelcomeMessage represents a welcome message from the Twitch WebSocket.
type SessionWelcomeMessage struct {
	Metadata Metadata `json:"metadata"`
	Payload  struct {
		Session Session `json:"session"`
	} `json:"payload"`
}

// KeepAliveMessage represents a keep-alive message from the Twitch WebSocket.
type KeepAliveMessage struct {
	Metadata Metadata `json:"metadata"`
	Payload  struct{} `json:"payload"`
}

// ChatMessage represents a chat message from the Twitch WebSocket.
type ChatMessage struct {
	Metadata Metadata `json:"metadata"`
	Payload  struct {
		Subscription Subscription `json:"subscription"`
		Event        struct {
			BroadcasterUserID    string `json:"broadcaster_user_id"`
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			BroadcasterUserName  string `json:"broadcaster_user_name"`
			ChatterUserID        string `json:"chatter_user_id"`
			ChatterUserLogin     string `json:"chatter_user_login"`
			ChatterUserName      string `json:"chatter_user_name"`
			MessageID            string `json:"message_id"`
			Message              struct {
				Text      string `json:"text"`
				Fragments []struct {
					Type      string       `json:"type"`
					Text      string       `json:"text"`
					Cheermote *interface{} `json:"cheermote,omitempty"` // Replace with appropriate type if needed
					Emote     *interface{} `json:"emote,omitempty"`     // Replace with appropriate type if needed
					Mention   *interface{} `json:"mention,omitempty"`   // Replace with appropriate type if needed
				} `json:"fragments"`
			} `json:"message"`
			Color  string `json:"color"`
			Badges []struct {
				SetID string `json:"set_id"`
				ID    string `json:"id"`
				Info  string `json:"info"`
			} `json:"badges"`
			MessageType                 string       `json:"message_type"`
			Cheer                       *interface{} `json:"cheer,omitempty"`                           // Replace with appropriate type if needed
			Reply                       *interface{} `json:"reply,omitempty"`                           // Replace with appropriate type if needed
			ChannelPointsCustomRewardID *interface{} `json:"channel_points_custom_reward_id,omitempty"` // Replace with appropriate type if needed
			ChannelPointsAnimationID    *interface{} `json:"channel_points_animation_id,omitempty"`     // Replace with appropriate type if needed
		} `json:"event"`
	} `json:"payload"`
}

func (cm *ChatMessage) HasBadge(name string) bool {
	for _, badge := range cm.Payload.Event.Badges {
		if badge.SetID == name {
			return true
		}
	}

	return false
}

// TwitchWebsocketMessage represents a union of possible messages received from the Twitch WebSocket.
type TwitchWebsocketMessage struct {
	SessionWelcome *SessionWelcomeMessage
	KeepAlive      *KeepAliveMessage
	Chat           *ChatMessage
}

func parseTwitchWebsocketMessage(rawJSON []byte) (*TwitchWebsocketMessage, error) {
	var base struct {
		Metadata Metadata `json:"metadata"`
	}

	if err := json.Unmarshal(rawJSON, &base); err != nil {
		return nil, err
	}

	msg := &TwitchWebsocketMessage{}

	switch base.Metadata.MessageType {

	case "session_welcome":
		var welcomeMsg SessionWelcomeMessage
		if err := json.Unmarshal(rawJSON, &welcomeMsg); err != nil {
			return nil, err
		}
		msg.SessionWelcome = &welcomeMsg

	case "session_keepalive":
		var keepAliveMsg KeepAliveMessage
		if err := json.Unmarshal(rawJSON, &keepAliveMsg); err != nil {
			return nil, err
		}
		msg.KeepAlive = &keepAliveMsg

	case "notification":
		if base.Metadata.SubscriptionType != nil && *base.Metadata.SubscriptionType == "channel.chat.message" {
			var chatMsg ChatMessage
			if err := json.Unmarshal(rawJSON, &chatMsg); err != nil {
				return nil, err
			}
			msg.Chat = &chatMsg
		}
	default:
		return nil, fmt.Errorf("unknown message type: %s", base.Metadata.MessageType)
	}

	return msg, nil
}

func DeleteMessage(auth AuthConfig, messageID string) error {
	// set URL and query
	reqURL, _ := url.Parse("https://api.twitch.tv/helix/moderation/chat")
	reqQuery := reqURL.Query()
	reqQuery.Add("broadcaster_id", auth.BroadcasterID)
	reqQuery.Add("moderator_id", auth.UserID)
	reqQuery.Add("message_id", messageID)
	reqURL.RawQuery = reqQuery.Encode()

	// create http req
	req, err := http.NewRequest("DELETE", reqURL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Client-Id", auth.ClientID)
	req.Header.Set("Authorization", auth.Bearer())
	req.Header.Set("Content-Type", "application/json")

	// send req
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 204 {
		return fmt.Errorf("failed to delete Twitch chat message (%d)", res.StatusCode)
	}

	return nil
}
