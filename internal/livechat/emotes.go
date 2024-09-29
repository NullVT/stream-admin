package livechat

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Emote struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Platform Platform `json:"platform"`
	Segment  string   `json:"segment"`
	MimeType string   `json:"mimetype"`
	FilePath string   `json:"filepath"`
}

type EmoteCache []Emote

func (ec *EmoteCache) FindByName(name string, platform Platform) *Emote {
	for i := range *ec {
		emote := &(*ec)[i]
		if platform == "" && emote.Name == name {
			return emote
		}
		if platform == emote.Platform && emote.Name == name {
			return emote
		}
	}
	return nil
}

func (ec *EmoteCache) FindByID(id string) *Emote {
	for i := range *ec {
		emote := &(*ec)[i]
		if emote.ID == id {
			return emote
		}
	}
	return nil
}

func (ec *EmoteCache) Update(platform Platform, segment string, name string, filepath string, mimetype string) {
	for i := range *ec {
		emote := &(*ec)[i]
		if emote.Platform == platform && emote.Name == name {
			emote.FilePath = filepath
			emote.MimeType = mimetype
			return
		}
	}
	// If not found, add new emote
	newEmote := Emote{
		ID:       uuid.New().String(),
		Name:     name,
		Platform: platform,
		Segment:  segment,
		FilePath: filepath,
		MimeType: mimetype,
	}
	*ec = append(*ec, newEmote)
}

func (ec *EmoteCache) Delete(id string) error {
	for idx, emote := range *ec {
		if emote.ID == id {
			// Remove the file associated with the emote
			if err := os.Remove(emote.FilePath); err != nil {
				return err // Return error if the file couldn't be deleted
			}

			// Remove the emote from the cache
			*ec = append((*ec)[:idx], (*ec)[idx+1:]...)

			// Return after the emote is found and removed
			return nil
		}
	}

	// If the emote ID was not found, return an error
	return fmt.Errorf("emote with ID %s not found", id)
}

func (ec *EmoteCache) SaveToFile(fileName string) error {
	// Ensure the parent directory exists
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// convert data into JSON
	content, err := json.Marshal(ec)
	if err != nil {
		return fmt.Errorf("failed to marshal EmoteCache to JSON: %w", err)
	}

	// Write content to the file
	if err := os.WriteFile(fileName, content, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func (ec *EmoteCache) LoadFromFile(fileName string) error {
	// Ensure the file exists before trying to load it
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %w", err)
	}

	// Read the file content
	content, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the JSON content into the EmoteCache struct
	if err := json.Unmarshal(content, ec); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
