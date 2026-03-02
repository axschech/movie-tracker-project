package repository

import (
	"context"
	"fmt"

	"github.com/axschech/rockbot-backend/internal/database"
	"github.com/axschech/rockbot-backend/internal/entities"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) QueryMedia(media entities.MediaEntity) ([]entities.MediaEntity, error) {
	rows, err := r.db.P.Query(r.Ctx, "SELECT * FROM media WHERE title ILIKE $1", "%"+media.Title+"%")

	if err != nil {
		return nil, err
	}

	var medias []entities.MediaEntity
	for rows.Next() {
		var m entities.MediaEntity
		if err := rows.Scan(&m.ID, &m.Title, &m.Runtime, &m.Type, &m.ImageURL, &m.Year); err != nil {
			return nil, err
		}
		medias = append(medias, m)
	}

	return medias, nil
}

func (r *Repository) CreateMedia(media []entities.MediaEntity) ([]entities.MediaEntity, error) {
	// need to add last updated
	batch := &pgx.Batch{}
	for _, m := range media {
		batch.Queue(`INSERT INTO media (title, runtime, type, image_url, year) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, runtime, type, image_url, year`, m.Title, m.Runtime, m.Type, m.ImageURL, m.Year)
	}

	br := r.db.P.SendBatch(r.Ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		err := br.QueryRow().Scan(&media[i].ID, &media[i].Title, &media[i].Runtime, &media[i].Type, &media[i].ImageURL, &media[i].Year)
		fmt.Printf("Inserted media into database: %+v\n", media[i])
		if err != nil {
			br.Close()
			return nil, fmt.Errorf("batch execution failed on index %d: %w", i, err)
		}
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch results: %w", err)
	}

	return media, nil
}

func (r *Repository) CreateMediaUser(mediaUser entities.MediaUserEntity) (entities.MediaUserEntity, error) {
	// keeping this single update for now
	err := r.db.P.QueryRow(r.Ctx, "INSERT INTO media_user (user_id, media_id, status) VALUES ($1, $2, $3) RETURNING id", mediaUser.UserID, mediaUser.MediaID, mediaUser.Status).Scan(&mediaUser.ID)

	if err != nil {
		return entities.MediaUserEntity{}, err
	}

	return mediaUser, nil
}

func (r *Repository) GetMediaUsers(mediaUser entities.MediaUserEntity, withMedia bool) ([]entities.MediaUserWithMediaEntity, error) {
	sql := "SELECT id, user_id, media_id, status FROM media_user WHERE user_id=$1 or media_id=$2"

	if withMedia {
		sql = "SELECT mu.id, mu.user_id, mu.media_id, mu.status, m.id, m.title, m.runtime, m.type, m.image_url, m.year FROM media_user mu JOIN media m ON mu.media_id = m.id WHERE mu.user_id=$1 or mu.media_id=$2"
	}
	rows, err := r.db.P.Query(r.Ctx, sql, mediaUser.UserID, mediaUser.MediaID)

	if err != nil {
		return nil, err
	}

	var mediaUsers []entities.MediaUserWithMediaEntity
	var m entities.MediaEntity
	for rows.Next() {
		var mu entities.MediaUserWithMediaEntity
		if err := rows.Scan(&mu.MediaUser.ID, &mu.MediaUser.UserID, &mu.MediaUser.MediaID, &mu.MediaUser.Status, &m.ID, &m.Title, &m.Runtime, &m.Type, &m.ImageURL, &m.Year); err != nil {
			return nil, err
		}
		mu.Medias = append(mu.Medias, m)
		mediaUsers = append(mediaUsers, mu)
	}

	return mediaUsers, nil
}
