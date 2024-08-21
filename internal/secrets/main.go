package secrets

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/nullvt/stream-admin/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
)

type Secrets map[string]string

const keyringService = "nullvt_stream_admin"
const secretsFile = "secrets.json"

var mu sync.Mutex

func Set(key string, value string) error {
	if config.Cfg.Server.Keyring {
		return keyring.Set(keyringService, key, value)
	}

	mu.Lock()
	defer mu.Unlock()

	secrets := make(Secrets)

	// Load existing secrets
	file, err := os.Open(secretsFile)
	if err == nil {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&secrets); err != nil {
			log.Error().Err(err).Msg("Failed to decode secrets file")
			return err
		}
		file.Close()
	}

	// Set the new key-value pair
	secrets[key] = value

	// Save the updated secrets
	file, err = os.Create(secretsFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create secrets file")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(secrets)
}

func Get(key string) (string, error) {
	if config.Cfg.Server.Keyring {
		return keyring.Get(keyringService, key)
	}

	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(secretsFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open secrets file")
		return "", err
	}
	defer file.Close()

	secrets := make(Secrets)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&secrets); err != nil {
		log.Error().Err(err).Msg("Failed to decode secrets file")
		return "", err
	}

	value, exists := secrets[key]
	if !exists {
		return "", nil // Or return an appropriate error if the key is not found
	}

	return value, nil
}
