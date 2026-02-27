package entities

type UserEntity struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
