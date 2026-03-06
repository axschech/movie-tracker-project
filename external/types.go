package external

import (
	"net/http"

	"github.com/axschech/rockbot-backend/internal/config"
)

type Sourcer interface {
	// was trying to figure out how to return a struct here instead
	// I think structs are just going to be too different.
	Fetch(title string) (http.Response, error)
}

type TVDBSource struct {
	Config config.Source
	Client *http.Client
	Type   string
}

type TokenRequest struct {
	APIKey string `json:"apikey"`
	PIN    string `json:"pin"`
}

type TokenData struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Data TokenData `json:"data"`
}

type TVSearchData struct {
	ImageURL string `json:"image_url"`
	Title    string `json:"title"`
	Runtime  string `json:"runtime"`
	Name     string `json:"name"`
	Year     string `json:"year"`
}

type TVSearchResponse struct {
	Data []TVSearchData `json:"data"`
}
