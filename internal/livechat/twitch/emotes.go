package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nullvt/stream-admin/internal/livechat"
)

type EmotesResponse interface {
	GetTemplate() string
	GetEmotes() *[]Emote
	GetSegment() string
}

type Emote struct {
	Segment   string
	ID        string
	Name      string
	Images    EmoteImages
	Format    []string
	Scale     []string
	ThemeMode []string
}

type EmoteImages struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

type ChannelEmotesResponse struct {
	Data     []ChannelEmoteData `json:"data"`
	Template string             `json:"template"`
	Segment  string             `json:"-"`
}

type ChannelEmoteData struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Images     EmoteImages `json:"images"`
	Tier       string      `json:"tier"`
	EmoteType  string      `json:"emote_type"`
	EmoteSetID string      `json:"emote_set_id"`
	Format     []string    `json:"format"`
	Scale      []string    `json:"scale"`
	ThemeMode  []string    `json:"theme_mode"`
}

func (cer *ChannelEmotesResponse) GetTemplate() string {
	return cer.Template
}

func (cer *ChannelEmotesResponse) GetSegment() string {
	return cer.Segment
}

func (cer *ChannelEmotesResponse) GetEmotes() *[]Emote {
	var emotes []Emote

	for _, emoteData := range cer.Data {
		emotes = append(emotes, Emote{
			ID:        emoteData.ID,
			Name:      emoteData.Name,
			Images:    emoteData.Images,
			Format:    emoteData.Format,
			Scale:     emoteData.Scale,
			ThemeMode: emoteData.ThemeMode,
		})
	}

	return &emotes
}

type GlobalEmotesResponse struct {
	Data     []GlobalEmoteData `json:"data"`
	Template string            `json:"template"`
}

type GlobalEmoteData struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Images    EmoteImages `json:"images"`
	Format    []string    `json:"format"`
	Scale     []string    `json:"scale"`
	ThemeMode []string    `json:"theme_mode"`
}

func (ger *GlobalEmotesResponse) GetTemplate() string {
	return ger.Template
}

func (ger *GlobalEmotesResponse) GetSegment() string {
	return "__global"
}

func (ger *GlobalEmotesResponse) GetEmotes() *[]Emote {
	var emotes []Emote

	for _, emoteData := range ger.Data {
		emotes = append(emotes, Emote{
			ID:        emoteData.ID,
			Name:      emoteData.Name,
			Images:    emoteData.Images,
			Scale:     emoteData.Scale,
			ThemeMode: emoteData.ThemeMode,
		})
	}

	return &emotes
}

func ListChannelEmotes(auth AuthConfig, channelID string) (*ChannelEmotesResponse, error) {
	// set URL and query
	reqURL, _ := url.Parse("https://api.twitch.tv/helix/chat/emotes")
	reqQuery := reqURL.Query()
	reqQuery.Add("broadcaster_id", channelID)
	reqURL.RawQuery = reqQuery.Encode()

	// create http req
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", auth.ClientID)
	req.Header.Set("Authorization", auth.Bearer())
	req.Header.Set("Content-Type", "application/json")

	// send req
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to list Twitch channel emotes (%d)", res.StatusCode)
	}

	// parse the response
	var resBody ChannelEmotesResponse
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	resBody.Segment = channelID

	return &resBody, nil
}

func ListGlobalEmotes(auth AuthConfig) (*GlobalEmotesResponse, error) {
	// create http req
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/chat/emotes/global", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", auth.ClientID)
	req.Header.Set("Authorization", auth.Bearer())
	req.Header.Set("Content-Type", "application/json")

	// send req
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// check response code
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to list Twitch global emotes (%d)", res.StatusCode)
	}

	// parse the response
	var resBody GlobalEmotesResponse
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}

	return &resBody, nil
}

func CacheEmotes(emoteCache *livechat.EmoteCache, emoteReq EmotesResponse) error {
	basePath := "./emotecache/twitch"

	// ensure dir exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err := os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// add new emote
	for _, emoteData := range *emoteReq.GetEmotes() {
		// get image URL "https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}"
		imgUrl := replaceMultiple(emoteReq.GetTemplate(), map[string]string{
			"{{id}}":         emoteData.ID,
			"{{format}}":     getPreferredFormat(emoteData.Format),
			"{{theme_mode}}": getPreferredThemeMode(emoteData.ThemeMode), // TODO: support lightmode
			"{{scale}}":      getPreferredScale(emoteData.Scale),
		})

		// fetch image from CDN
		res, err := http.Get(imgUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// Read the first 512 bytes to detect the content type
		buffer := make([]byte, 512)
		_, err = res.Body.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		contentType := http.DetectContentType(buffer)
		var ext string
		switch contentType {
		case "image/jpeg":
			ext = "jpg"
		case "image/png":
			ext = "png"
		case "image/gif":
			ext = "gif"
		default:
			// TODO: proper error handling
			panic("Unsupported content type: " + contentType)
		}

		// write image to disk
		filename := fmt.Sprintf("%s/%s.%s", basePath, emoteData.ID, ext)
		file, err := os.Create(fmt.Sprintf("%s/%s.%s", basePath, emoteData.ID, ext))
		if err != nil {
			return err
		}
		defer file.Close()

		// Reset the response body reader and copy the image data to the file
		res.Body = io.NopCloser(io.MultiReader(strings.NewReader(string(buffer)), res.Body))
		_, err = io.Copy(file, res.Body)
		if err != nil {
			return err
		}

		// add emote to map
		emoteCache.Update(livechat.Twitch, emoteReq.GetSegment(), emoteData.Name, filename, contentType)
	}

	// remove old emotes
	for idx := 0; idx < len(*emoteCache); {
		cachedEmote := (*emoteCache)[idx]
		keep := true

		// Check if the cached emote matches the platform and segment in emoteReq
		if cachedEmote.Platform == livechat.Twitch && cachedEmote.Segment == emoteReq.GetSegment() {
			for _, emote := range *emoteReq.GetEmotes() {
				keep = false
				if cachedEmote.Name == emote.Name {
					keep = true
					break
				}
			}
		}

		// delete the emote
		if !keep {
			if err := emoteCache.Delete(cachedEmote.ID); err != nil {
				fmt.Printf("Error removing emote: %v\n", err)
			}
			continue
		} else {
			idx++
		}
	}

	return nil
}

func SyncEmotes(emoteCache *livechat.EmoteCache, auth AuthConfig, channelIDs []string) error {
	// sync global emotes
	globalEmotes, err := ListGlobalEmotes(auth)
	if err != nil {
		return err
	}
	if err := CacheEmotes(emoteCache, globalEmotes); err != nil {
		return err
	}

	// sync channel emotes
	for _, channelID := range channelIDs {
		channelEmotes, err := ListChannelEmotes(auth, channelID)
		if err != nil {
			return err
		}
		if err := CacheEmotes(emoteCache, channelEmotes); err != nil {
			return err
		}
	}

	log.Print(channelIDs)

	// remove channels
	for idx := 0; idx < len(*emoteCache); {
		cachedEmote := (*emoteCache)[idx]

		// Only operate on Twitch platform emotes
		if cachedEmote.Platform == livechat.Twitch {
			keep := false

			// Check if the emote belongs to the global segment
			if cachedEmote.Segment == globalEmotes.GetSegment() {
				keep = true
			}

			// If not global, check if the segment matches any channelID
			for _, channelID := range channelIDs {
				if cachedEmote.Segment == channelID {
					keep = true
					break
				}
			}

			// If the emote should not be kept then delete it
			if !keep {
				if err := emoteCache.Delete(cachedEmote.ID); err != nil {
					return err
				}
				continue
			}
		}

		idx++
	}

	return nil
}
