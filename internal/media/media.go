package media

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/axschech/rockbot-backend/external"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/entities"
)

type Media struct {
	R repository.Repository
	S external.Sourcer
}

func NewMedia(r repository.Repository, s external.Sourcer) *Media {
	return &Media{
		R: r,
		S: s,
	}
}

func (m *Media) GetOrSaveMedia(title string, mediaType string) ([]entities.MediaEntity, error) {
	var (
		medias []entities.MediaEntity
		err    error
	)

	// should also check media type
	medias, err = m.R.QueryMedia(entities.MediaEntity{Title: title})
	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Printf("Failed to query media from database: %v\n", err)
		return nil, fmt.Errorf("failed to query media: %w", err)
	}
	fmt.Printf("Queried media from database: %v\n", medias)
	// could update this to get both results and check if the length is equal?
	// or could do a more refined checked instead of doing things in batches
	// might also just allow user to manually insert media
	if len(medias) > 0 {
		fmt.Printf("Found media in database: %v\n", medias)
		return medias, nil
	}

	resp, err := m.S.Fetch(title)
	if err != nil {
		fmt.Printf("Failed to fetch media from external source: %v\n", err)
		return nil, fmt.Errorf("failed to fetch media from external source: %w", err)
	}

	// not sure if this belongs here, might need to be imported from source package
	var tvSearchResponse external.TVSearchResponse
	switch mediaType {
	case "tvdb":
		err = json.NewDecoder(resp.Body).Decode(&tvSearchResponse)
		defer resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to decode response from external source: %w", err)
		}
		for _, data := range tvSearchResponse.Data {
			medias = append(medias, entities.MediaEntity{
				Title:    data.Name,
				Runtime:  data.Runtime,
				Type:     "tv",
				ImageURL: data.ImageURL,
				Year:     data.Year,
			})
		}
	}

	if err != nil {
		fmt.Printf("Failed to convert response to media structs: %v\n", err)
		return nil, fmt.Errorf("failed to convert response to media structs: %w", err)
	}

	fmt.Printf("Inserting media into database: %v\n", medias)

	medias, err = m.R.CreateMedia(medias)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			fmt.Printf("Media already exists in database: %v\n", err)
			return medias, nil
		}

		return nil, fmt.Errorf("failed to insert media into database: %w", err)
	}

	return medias, nil
}
