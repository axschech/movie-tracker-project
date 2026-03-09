package media_user

import (
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/entities"
)

var AllowedStatuses = []string{"not watched", "will watch", "watching", "watched", "wont watch"}

type MediaUser struct {
	R repository.Repository
}

func NewMediaUser(r repository.Repository) *MediaUser {
	return &MediaUser{R: r}
}

func (um *MediaUser) CreateMediaUser(mediaUser entities.MediaUserEntity) (entities.MediaUserEntity, error) {
	// may want to update this to insert or update
	return um.R.CreateMediaUser(mediaUser)
}

func (um *MediaUser) UpdateMediaUser(mediaUser entities.MediaUserEntity) (entities.MediaUserEntity, error) {
	// may want to update this to insert or update
	return um.R.UpdateMediaUser(mediaUser)
}

func (um *MediaUser) QueryMediaUsersWithUserID(userID int, withMedia bool) ([]entities.MediaUserWithMediaEntity, error) {
	mediaUser := entities.MediaUserEntity{UserID: userID}

	return um.R.GetMediaUsers(mediaUser, withMedia)
}
