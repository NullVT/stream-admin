package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SubscriptionRequest struct {
	Type      string                       `json:"type"`
	Version   string                       `json:"version"`
	Condition SubscriptionRequestCondition `json:"condition"`
	Transport SubscriptionRequestTransport `json:"transport"`
}

type SubscriptionRequestCondition struct {
	UserID            string `json:"user_id"`
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type SubscriptionRequestTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

// Subscription represents a subscription to a Twitch event.
type Subscription struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
		UserID            string `json:"user_id"`
	} `json:"condition"`
	Transport struct {
		Method    string `json:"method"`
		SessionID string `json:"session_id"`
	} `json:"transport"`
	CreatedAt string `json:"created_at"` // ISO 8601 timestamp
	Cost      int    `json:"cost"`
}

type EventSubSubscriptionsResponse struct {
	Data []Subscription `json:"data"`
}

func Subscribe(authConfig AuthConfig, sessionID string, subType string) (string, error) {
	// create and marshal request body
	body, err := json.Marshal(SubscriptionRequest{
		Type:    subType,
		Version: "1",
		Condition: SubscriptionRequestCondition{
			UserID:            authConfig.UserID,
			BroadcasterUserID: authConfig.BroadcasterID,
		},
		Transport: SubscriptionRequestTransport{
			Method:    "websocket",
			SessionID: sessionID,
		},
	})
	if err != nil {
		return "", err
	}

	// create http req
	req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", authConfig.Bearer())
	req.Header.Set("Client-Id", authConfig.ClientID)
	req.Header.Set("Content-Type", "application/json")

	// send req
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 202 {
		return "", fmt.Errorf("failed to subscribe to chat events (%d)", res.StatusCode)
	}

	// parse the response
	var resBody EventSubSubscriptionsResponse
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", err
	}

	// find correct subscription and return the ID
	for _, sub := range resBody.Data {
		if sub.Transport.Method == "websocket" && sub.Transport.SessionID == sessionID && sub.Condition.UserID == authConfig.UserID && sub.Condition.BroadcasterUserID == authConfig.BroadcasterID {
			return sub.ID, nil
		}
	}

	return "", errors.New("no matching subscription found")
}
