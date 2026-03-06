package service

type PostUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PostMediaUserRequest struct {
	UserID  int    `json:"user_id"`
	MediaID int    `json:"media_id"`
	Status  string `json:"status"`
}
