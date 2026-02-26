package user

import "github.com/axschech/rockbot-backend/internal/entities"

func GetUserByID(ID string) entities.User {
	return entities.User{
		ID:       ID,
		Username: "testuser",
		Email:    "testuser@example.com",
	}
}
