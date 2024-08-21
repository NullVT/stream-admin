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
	}

	viper.Set("twitch.clientId", defaultConfig.Twitch.ClientID)
	viper.Set("server.host", defaultConfig.Server.Host)
	viper.Set("server.port", defaultConfig.Server.Port)
	viper.Set("server.baseUrl", defaultConfig.Server.BaseURL)
}
