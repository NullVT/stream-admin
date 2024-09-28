package livechat

import (
	"time"
)

type User struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Broadcaster   bool   `json:"broadcaster"`
	Moderator     bool   `json:"moderator"`
	TwitchVIP     bool   `json:"twitch_vip"`
	YouTubeMember bool   `json:"youtube_member"`
}

type Message struct {
	ID          string         `json:"id"`
	Platform    Platform       `json:"platform"`
	Body        string         `json:"body"`
	Emotes      []MessageEmote `json:"emotes"`
	Sender      User           `json:"sender"`
	ReceivedAt  time.Time      `json:"received_at"`
	PublishedAt time.Time      `json:"published_at"`
}

type MessageEmote struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
