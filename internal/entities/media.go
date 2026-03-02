package entities

// need to add more fields here
// should be distinguished between tv and movies
// year + name + type is how we check if it exists already
// had to use year, since names can be the same but it can be a relelease or a remake, etc.
type MediaEntity struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Runtime  string `json:"runtime"`
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
	Year     string `json:"year"`
}
