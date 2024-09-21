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

// LoadSecrets loads the secrets from the file, creating a new map if the file doesn't exist
func loadSecrets() (Secrets, error) {
	var secrets Secrets

	file, err := os.Open(secretsFile)
	if err != nil {
		if os.IsNotExist(err) {
			secrets = make(Secrets) // Create a new Secrets map if the file doesn't exist
			return secrets, nil
		}
		log.Error().Err(err).Msg("Failed to open secrets file")
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close secrets file")
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&secrets); err != nil {
		log.Error().Err(err).Msg("Failed to decode secrets file")
		return nil, err
	}

	return secrets, nil
}

// SaveSecrets saves the secrets map to the file
func saveSecrets(secrets Secrets) error {
	file, err := os.Create(secretsFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create secrets file")
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close secrets file after writing")
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(secrets); err != nil {
		log.Error().Err(err).Msg("Failed to encode secrets to file")
		return err
	}

	return nil
}

func Set(key string, value string) error {
	if config.Cfg.Server.Keyring {
		return keyring.Set(keyringService, key, value)
	}

	mu.Lock()
	defer mu.Unlock()

	// Load existing secrets
	secrets, err := loadSecrets()
	if err != nil {
		return err
	}

	// Set the new key-value pair
	secrets[key] = value

	// Save the updated secrets
	if err := saveSecrets(secrets); err != nil {
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	if config.Cfg.Server.Keyring {
		return keyring.Get(keyringService, key)
	}

	mu.Lock()
	defer mu.Unlock()

	// Load existing secrets
	secrets, err := loadSecrets()
	if err != nil {
		return "", err
	}

	value, exists := secrets[key]
	if !exists {
		return "", nil // Return an empty string if the key is not found
	}

	return value, nil
}
