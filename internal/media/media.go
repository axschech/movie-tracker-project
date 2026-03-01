package media

import (
	"github.com/axschech/rockbot-backend/external"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/entities"
)

type Media struct {
	R repository.Repository
	S external.Sourcer
}

func NewMedia(title string, mediaType string, r repository.Repository, s external.Sourcer) *Media {
	return &Media{
		R: r,
		S: s,
	}
}

func (m *Media) GetOrSaveMedia(title string, t string) []entities.MediaEntity {
	// var sourceType string
	// switch t {
	// case "tv":
	// 	sourceType = "tvdb"
	// }

	// rows, err := m.R.QueryMedia(entities.MediaEntity{Title: title})
	// if err != nil && !strings.Contains(err.Error(), "no rows") {
	// 	fmt.Printf("Failed to query media from database: %v\n", err)
	// 	return nil
	// }

	return nil
}
