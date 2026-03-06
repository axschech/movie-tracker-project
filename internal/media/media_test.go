package media

import (
	"context"
	"net/http"
	"testing"

	"github.com/axschech/rockbot-backend/internal/database"
	"github.com/axschech/rockbot-backend/internal/database/repository"

	"github.com/pashagolub/pgxmock/v5"
)

type fakeSourcer struct {
	Response http.Response
	Err      error
}

func (f *fakeSourcer) Fetch(title string) (http.Response, error) {
	return f.Response, f.Err
}

func TestGetOrSaveMediaWithoutFetch(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	db := database.Database{P: mock}
	rep := repository.NewRepository(ctx, &db)

	fakeSourcer := &fakeSourcer{
		Response: http.Response{},
	}

	mock.ExpectQuery("SELECT \\* FROM media WHERE title ILIKE \\$1").WithArgs("%Test Movie%").WillReturnRows(pgxmock.NewRows([]string{"id", "title", "runtime", "type", "image_url", "year"}).AddRow("1", "Test Movie", "120 min", "movie", "http://example.com/image.jpg", "2021"))

	media := NewMedia(*rep, fakeSourcer)

	got, err := media.GetOrSaveMedia("Test Movie", "movie")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if got == nil {
		t.Fatal("Expected media to be returned, got nil")
	}

	if len(got) != 1 {
		t.Fatalf("Expected 1 media to be returned, got %d", len(got))
	}

	// use assert library here
	if got[0].Title != "Test Movie" {
		t.Errorf("Expected title to be 'Test Movie', got '%s'", got[0].Title)
	}

	if got[0].Runtime != "120 min" {
		t.Errorf("Expected runtime to be '120 min', got '%s'", got[0].Runtime)
	}

	if got[0].Type != "movie" {
		t.Errorf("Expected type to be 'movie', got '%s'", got[0].Type)
	}

	if got[0].ImageURL != "http://example.com/image.jpg" {
		t.Errorf("Expected image URL to be 'http://example.com/image.jpg', got '%s'", got[0].ImageURL)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expected no unfulfilled expectations, got %v", err)
	}
}

func TestGetOrSaveMediaWithFetchError(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	db := database.Database{P: mock}
	rep := repository.NewRepository(ctx, &db)

	fakeSourcer := &fakeSourcer{
		Err: http.ErrHandlerTimeout,
	}

	media := NewMedia(*rep, fakeSourcer)

	mock.ExpectQuery("SELECT \\* FROM media WHERE title ILIKE \\$1").WithArgs("%Test Movie%").WillReturnRows(pgxmock.NewRows([]string{"id", "title", "runtime", "type", "image_url", "year"}))

	got, err := media.GetOrSaveMedia("Test Movie", "movie")

	if err == nil {
		t.Logf("error: %v", err)
		t.Fatal("Expected error, got nil")
	}

	if got != nil {
		t.Fatalf("Expected media to be nil, got %v", got)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Expected no unfulfilled expectations, got %v", err)
	}
}
