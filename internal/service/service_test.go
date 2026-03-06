package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/database"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/routing"
	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v5"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

const defaultTokenResponse = `{"data":{"token":"testtoken"}}`
const defaultSearchResponse = `{"data":[{"image_url":"http://example.com/image.jpg","title":"Test Show","runtime":"60 min","name":"Test Show","year":"2021"}]}`

func TestGetUserHandler(t *testing.T) {
	cfg, err := config.MakeConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery("SELECT id, username, email FROM users WHERE id=\\$1 or username=\\$2 or email=\\$3").WithArgs(1, "", "").WillReturnRows(pgxmock.NewRows([]string{"id", "username", "email"}).AddRow(1, "testuser", "testuser@example.com"))

	db := database.Database{P: mock}
	rep := repository.NewRepository(ctx, &db)
	r := routing.NewRouter("")
	service := NewService(cfg, *rep, r, &http.Client{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/user/1", nil)
	// im not sure why but chi.URLParam is not working in the test, so we have to add the route context manually
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1") // Add with an empty string
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	service.GetUserHandler(w, req)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	want := `{"id":1,"username":"testuser","email":"testuser@example.com"}
`
	got := w.Body.String()
	if got != want {
		t.Errorf("unexpected response body: got %q, want %q", got, want)
	}
}

func TestPostUserHandler(t *testing.T) {
	cfg, err := config.MakeConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id")).WithArgs("newuser", "newuser@example.com").WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))

	db := database.Database{P: mock}
	rep := repository.NewRepository(ctx, &db)
	r := routing.NewRouter("")
	service := NewService(cfg, *rep, r, &http.Client{})

	w := httptest.NewRecorder()

	body := PostUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
	}

	bts, err := json.Marshal(body)
	if err != nil {
		log.Fatal("json.Marshal failed:", err)
	}

	req := httptest.NewRequest("POST", "/api/user", bytes.NewReader(bts))

	service.PostUserHandler(w, req)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	want := `{"id":2,"username":"newuser","email":"newuser@example.com"}
`
	got := w.Body.String()
	if got != want {
		t.Errorf("unexpected response body: got %q, want %q", got, want)
	}
}

func TestQueryMediaHandlerValidation(t *testing.T) {
	tests := []struct {
		name string
		q    string
	}{
		{
			name: "empty query",
			q:    "/api/media/query?query=&type=tv",
		},
		{
			name: "invalid type",
			q:    "/api/media/query?query=test&type=invalid",
		},
		{
			name: "missing type",
			q:    "/api/media/query?query=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.MakeConfig()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			db := database.Database{P: mock}
			rep := repository.NewRepository(ctx, &db)
			r := routing.NewRouter("")
			service := NewService(cfg, *rep, r, &http.Client{})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.q, nil)

			service.QueryMediaHandler(w, req)
			want := http.StatusBadRequest

			if w.Result().StatusCode != want {
				t.Fatalf("unexpected status code: got %d, want %d", w.Result().StatusCode, want)
			}
		})
	}
}

func TestQueryMediaHandlerWithInsert(t *testing.T) {
	cfg, err := config.MakeConfig()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM media WHERE title ILIKE $1")).WithArgs("%test%").WillReturnRows(pgxmock.NewRows([]string{"id", "title", "runtime", "type", "image_url", "year"}))

	batch := mock.ExpectBatch()
	batch.ExpectQuery(regexp.QuoteMeta("INSERT INTO media (title, runtime, type, image_url, year) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, runtime, type, image_url, year")).WithArgs("Test Show", "60 min", "tv", "http://example.com/image.jpg", "2021").WillReturnRows(pgxmock.NewRows([]string{"id", "title", "runtime", "type", "image_url", "year"}).AddRow("1", "Test Show", "60 min", "tv", "http://example.com/image.jpg", "2021"))

	db := database.Database{P: mock}
	rep := repository.NewRepository(ctx, &db)
	r := routing.NewRouter("")

	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		var rawRes string
		switch req.URL.Path {
		case "/login":
			rawRes = defaultTokenResponse
		case "/search":
			rawRes = defaultSearchResponse
		default:
			return nil, http.ErrNotSupported
		}

		res := io.NopCloser(strings.NewReader(rawRes))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       res,
		}, nil
	})}
	service := NewService(cfg, *rep, r, httpClient)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/media/query?query=test&type=tv", nil)

	service.QueryMediaHandler(w, req)

	want := `[{"id":"1","title":"Test Show","runtime":"60 min","type":"tv","image_url":"http://example.com/image.jpg","year":"2021"}]
`
	got := w.Body.String()
	if got != want {
		t.Errorf("unexpected response body: got %q, want %q", got, want)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
