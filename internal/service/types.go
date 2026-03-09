package service

type PostUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PostMediaUserRequest struct {
	UserID  int    `json:"userId"`
	MediaID int    `json:"mediaId"`
	Status  string `json:"status"`
}

type PutMediaUserRequest struct {
	MediaUserID int    `json:"mediaUserId"`
	Status      string `json:"status"`
}
