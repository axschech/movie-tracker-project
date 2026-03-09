package entities

// linking users with media when they set a status
type MediaUserEntity struct {
	ID      int `json:"id"`
	MediaID int `json:"media"`
	UserID  int `json:"user"`
	// tried using an enum here but it was taking too much time, might come back
	Status string `json:"status"`
}

type MediaUserWithMediaEntity struct {
	MediaUser MediaUserEntity `json:"mediaUser"`
	Media     MediaEntity     `json:"media"`
}
