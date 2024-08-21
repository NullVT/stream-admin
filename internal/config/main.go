package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Twitch TwitchConfig `mapstructure:"twitch"`
	Server ServerConfig `mapstructure:"server"`
}
type TwitchConfig struct {
	ClientID string `mapstructure:"clientId"`
}
type ServerConfig struct {
	Host    string
	Port    uint16
	BaseURL string `mapstructure:"baseUrl"`
	Keyring bool
}

// Global variable to hold the loaded config.
var Cfg = &Config{}

// Load handles loading the configuration from the appropriate file.
func Load() error {
	// Initialize zerolog with default settings
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Determine the appropriate OS-specific config directory
	var configDir string
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("APPDATA"), "nullvt")
	} else {
		configDir = "/etc/nullvt"
	}
	configFileName := "stream-admin.toml"

	// Add the working directory and the OS-specific config directory to the search paths
	viper.SetConfigName("stream-admin") // This will be used for both working directory and OS-specific directory files
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")       // Check the working directory first
	viper.AddConfigPath(configDir) // Check the OS-specific directory

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// If the config file is not found, create it in the OS-specific directory
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn().Msg("Config file not found. Creating new config file...")
			if createErr := createConfigFile(configDir, configFileName); createErr != nil {
				log.Error().Err(createErr).Msg("failed to create config file")
				return fmt.Errorf("failed to create config file: %w", createErr)
			}

			// Reattempt to load the newly created config file
			if err = viper.ReadInConfig(); err != nil {
				log.Error().Err(err).Msg("failed to read newly created config file")
				return fmt.Errorf("failed to read newly created config file: %w", err)
			}
		} else {
			// Handle other errors that might occur during reading the config
			log.Error().Err(err).Msg("failed to read config file")
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	log.Debug().Interface("AllSettings", viper.AllSettings()).Msg("Viper loaded settings")

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Error().Err(err).Msg("unable to decode into struct")
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	log.Debug().Interface("Config", Cfg).Msg("Config loaded successfully")
	return nil
}

// createConfigFile creates a new config file with default values in the specified directory.
func createConfigFile(configDir, configFileName string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Error().Err(err).Msg("failed to create config directory")
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Define the full path for the config file
	fullPath := filepath.Join(configDir, configFileName)

	// set config values
	setDefaults()

	// Write the config file
	viper.SetConfigFile(fullPath)
	if err := viper.WriteConfigAs(fullPath); err != nil {
		log.Error().Err(err).Msg("failed to write config file")
		return fmt.Errorf("failed to write config file: %w", err)
	}

	log.Info().Str("path", fullPath).Msg("Config file created")
	return nil
}

// SetConfigValue sets a configuration value and persists it to the config file.
func SetConfigValue(key string, value interface{}) error {
	// Set the new value in Viper's in-memory configuration
	viper.Set(key, value)

	// Persist the updated configuration to the file
	if err := viper.WriteConfig(); err != nil {
		log.Error().Err(err).Msg("failed to write config to file")
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	log.Info().Str("key", key).Msg("Config value updated successfully")
	return nil
}
