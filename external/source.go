package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/axschech/rockbot-backend/internal/config"
)

func (s *TVDBSource) GetToken() (string, error) {
	fmt.Printf("Getting token with API key: %s and PIN: %s\n", s.Config.APIKey, s.Config.PIN)
	tr := TokenRequest{APIKey: s.Config.APIKey, PIN: s.Config.PIN}
	trBytes, err := json.Marshal(tr)
	if err != nil {
		return "", err
	}
	fmt.Printf("Token request body: %s\n", string(trBytes))
	req, err := http.NewRequest("POST", s.Config.BaseURL+"/login", bytes.NewBuffer(trBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// handle status code

	var tokenRes TokenResponse
	err = json.NewDecoder(res.Body).Decode(&tokenRes)
	if err != nil {
		return "", err
	}
	fmt.Printf("Token response: %v\n", tokenRes)
	return tokenRes.Data.Token, nil
}

func (s *TVDBSource) Fetch(title string) (http.Response, error) {
	token, err := s.GetToken()
	if err != nil {
		return http.Response{}, err
	}

	req, err := http.NewRequest("GET", s.Config.BaseURL+"/search?language=en&type=series&query="+title, nil)
	if err != nil {
		return http.Response{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	// handle status code

	return *resp, nil
}

func NewTVDBSource(cfg config.Source, client *http.Client) *TVDBSource {
	return &TVDBSource{Config: cfg, Client: client}
}

func GetSource(client *http.Client, cfg config.Source, source string) Sourcer {
	// validate source?
	sourceMap := map[string]Sourcer{
		"tvdb": NewTVDBSource(cfg, client),
	}

	return sourceMap[source]
}
