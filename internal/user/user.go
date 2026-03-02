package user

import (
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/entities"
)

type User struct {
	R repository.Repository
}

func NewUser(r repository.Repository) *User {
	return &User{R: r}
}

// TODO: Add more generic GetUser function?
func (u *User) GetUserByID(id int) (entities.UserEntity, error) {
	user, err := u.R.GetUser(entities.UserEntity{ID: id})
	if err != nil {
		return entities.UserEntity{}, err
	}
	return user, nil
}

func (u *User) Register(username, email string) (entities.UserEntity, error) {
	user, err := u.R.CreateUser(entities.UserEntity{Username: username, Email: email})
	if err != nil {
		return entities.UserEntity{}, err
	}

	return user, nil
}
