package entities

// need to add more fields here
// should be distinguished between tv and movies
// also need stuff like genre, release date, etc. but for now just keeping it simple
type MediaEntity struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Runtime  string `json:"runtime"`
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
}
