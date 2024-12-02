package config

import "github.com/spf13/viper"

func setDefaults() {
	defaultConfig := Config{
		Twitch: TwitchConfig{
			ClientID: "",
		},
		Server: ServerConfig{
			Host:    "localhost",
			Port:    8080,
			BaseURL: "http://localhost:8080/",
		},
		EmotesWhitelist:   map[string]string{},
		StreamInfoPresets: []StreamInfoPreset{},
	}

	viper.SetDefault("twitch.clientId", defaultConfig.Twitch.ClientID)
	viper.SetDefault("server.host", defaultConfig.Server.Host)
	viper.SetDefault("server.port", defaultConfig.Server.Port)
	viper.SetDefault("server.baseUrl", defaultConfig.Server.BaseURL)
	viper.SetDefault("emotesWhitelist", defaultConfig.EmotesWhitelist)
	viper.SetDefault("streamInfoPresets", defaultConfig.StreamInfoPresets)
}
