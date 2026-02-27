package repository

import (
	"context"

	"github.com/axschech/rockbot-backend/internal/database"
	"github.com/axschech/rockbot-backend/internal/entities"
)

type Repository struct {
	db  *database.Database
	Ctx context.Context
}

func NewRepository(ctx context.Context, db *database.Database) *Repository {
	return &Repository{db: db, Ctx: ctx}
}

// putting all repositories here for now but would probably separate them out
func (r *Repository) GetUser(user entities.UserEntity) (entities.UserEntity, error) {
	err := r.db.P.QueryRow(r.Ctx, "SELECT id, username, email FROM users WHERE id=$1 or username=$2 or email=$3", user.ID, user.Username, user.Email).Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return entities.UserEntity{}, err
	}

	return user, nil
}

func (r *Repository) CreateUser(user entities.UserEntity) (entities.UserEntity, error) {
	err := r.db.P.QueryRow(r.Ctx, "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id", user.Username, user.Email).Scan(&user.ID)

	if err != nil {
		return entities.UserEntity{}, err
	}

	return user, nil
}

func (r *Repository) GetMedia(media entities.MediaEntity) (entities.MediaEntity, error) {
	err := r.db.P.QueryRow(r.Ctx, "SELECT id, title FROM media WHERE id=$1 or title=$2", media.ID, media.Title).Scan(&media.ID, &media.Title)

	if err != nil {
		return entities.MediaEntity{}, err
	}

	return media, nil
}

func (r *Repository) CreateMedia(media entities.MediaEntity) (entities.MediaEntity, error) {
	err := r.db.P.QueryRow(r.Ctx, "INSERT INTO media (title, runtime, type, image_url) VALUES ($1, $2, $3, $4) RETURNING id", media.Title, media.Runtime, media.Type, media.ImageURL).Scan(&media.ID)

	if err != nil {
		return entities.MediaEntity{}, err
	}

	return media, nil
}
