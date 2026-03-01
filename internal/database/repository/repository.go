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

func (r *Repository) CreateMedia(media []entities.MediaEntity) error {
	// need to add last updated
	batch := &pgx.Batch{}
	for _, m := range media {
		batch.Queue(`INSERT INTO media (title, runtime, type, image_url, year) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, runtime, type, image_url, year`, m.Title, m.Runtime, m.Type, m.ImageURL, m.Year)
	}

	br := r.db.P.SendBatch(r.Ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		row, err := br.Query()
		fmt.Printf("Batch result for index %d: %+v\n", i, row)
		if err != nil {
			br.Close()
			return fmt.Errorf("batch execution failed on index %d: %w", i, err)
		}

	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close batch results: %w", err)
	}

	return nil
}
