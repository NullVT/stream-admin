package twitch

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type User struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	BroadcasterType string `json:"broadcaster_type"`
}

type UsersResponse struct {
	Data []User `json:"data"`
}

func GetUsers(auth AuthConfig, userNames []string) ([]User, error) {
	// set userName
	queryParams := url.Values{}
	if len(userNames) >= 1 {
		for _, userName := range userNames {
			queryParams.Add("login", userName)
		}
	}

	baseURL := url.URL{
		Scheme:   "https",
		Host:     "api.twitch.tv",
		Path:     "/helix/users",
		RawQuery: queryParams.Encode(),
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth.Bearer())
	req.Header.Set("Client-Id", auth.ClientID)

	// send request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// check if the response status is not 200 OK
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch current user: status code " + res.Status)
	}

	// decode the response body
	var usersResponse UsersResponse
	err = json.NewDecoder(res.Body).Decode(&usersResponse)
	if err != nil {
		return nil, err
	}

	// check if the data array is empty
	if len(usersResponse.Data) == 0 {
		return nil, errors.New("no user data found")
	}

	return usersResponse.Data, nil
}

func (usr *User) MarshalString() (string, error) {
	jsonBytes, err := json.Marshal(usr)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (usr *User) UnmarshalString(jsonString string) error {
	if err := json.Unmarshal([]byte(jsonString), usr); err != nil {
		return err
	}

	return nil
}
