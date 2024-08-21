package twitch

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// generate a login URL to request an Oauth token
func OAuthLogin(clientId string, redirectURI string, scopes []string) (string, error) {
	if clientId == "" {
		return "", errors.New("invalid ClientID")
	}

	queryParams := url.Values{}
	queryParams.Add("response_type", "token")
	queryParams.Add("client_id", clientId)
	queryParams.Add("redirect_uri", redirectURI)
	queryParams.Add("scope", strings.Join(scopes, " "))

	baseURL := url.URL{
		Scheme:   "https",
		Host:     "id.twitch.tv",
		Path:     "/oauth2/authorize",
		RawQuery: queryParams.Encode(),
	}

	return baseURL.String(), nil
}

// parse the auth token from the callback URL
func OAuthCallback(callbackURL string) (string, error) {
	parsedURL, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}

	// get the params from the query... or the fragment because Twitch is bloody weird.
	rawParams := parsedURL.RawQuery
	if parsedURL.RawQuery == "" && parsedURL.RawFragment != "" {
		rawParams = parsedURL.Fragment
	}

	// parse the query params
	if rawParams == "" {
		return "", errors.New("oauth callback contains no parameters")
	}
	params, err := url.ParseQuery(rawParams)
	if err != nil {
		return "", err
	}

	// validate the params
	if params.Has("status") {
		return "", fmt.Errorf("OAuth call back error: %s", params.Get("error"))
	}
	if !params.Has("access_token") {
		return "", errors.New("OAuth call back did not contain an access_token")
	}

	return params.Get("access_token"), nil
}

func OAuthValidateToken(token string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return false, err
	}

	// validate token
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	// token is valid
	if res.StatusCode == 200 {
		return true, nil
	}

	// invalid token
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	return false, errors.New(string(body))
}
