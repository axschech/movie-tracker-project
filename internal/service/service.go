package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/axschech/rockbot-backend/external"
	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/entities"
	"github.com/axschech/rockbot-backend/internal/routing"
	"github.com/axschech/rockbot-backend/internal/user"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Config     config.Config
	Repository repository.Repository
	Router     routing.Router
}

func NewService(
	cfg config.Config,
	r repository.Repository,
	router routing.Router,
) *Service {
	return &Service{
		Config:     cfg,
		Repository: r,
		Router:     router,
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func (s *Service) Run() error {
	s.Router.R.Get("/ping", Ping)
	s.Router.R.Route("/api", func(r chi.Router) {
		r.Get("/user/{id}", s.GetUserHandler)
		r.Post("/user", s.PostUserHandler)
		r.Get("/search/media", s.QueryMediaHandler)
	})

	return s.Router.Listen()
}

// handlers should probably be their own structs, with an interface called Handlerer
func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	u := user.NewUser(s.Repository)
	userId := chi.URLParam(r, "id")

	user, err := u.GetUserByID(userId)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	u := user.NewUser(s.Repository)

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" {
		http.Error(w, "Username and email are required", http.StatusBadRequest)
		return
	}

	user, err := u.Register(req.Username, req.Email)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) QueryMediaHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	t := r.URL.Query().Get("type")
	if t == "" {
		http.Error(w, "Type parameter is required", http.StatusBadRequest)
		return
	}
	var (
		medias []entities.MediaEntity
		err    error
	)
	medias, err = s.Repository.QueryMedia(entities.MediaEntity{Title: query})
	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Printf("Failed to query media from database: %v\n", err)
		http.Error(w, "Failed to query media", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Queried media from database: %v\n", medias)
	// could update this to get both results and check if the length is equal?
	// or could do a more refined checked instead of doing things in batches
	// might also just allow user to manually insert media
	if len(medias) > 0 {
		fmt.Printf("Found media in database: %v\n", medias)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(medias)
		return
	}

	var mediaType string
	switch t {
	case "tv":
		mediaType = "tvdb"
	}

	if mediaType == "" {
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}

	source := config.Source{
		BaseURL: s.Config.TVSource.BaseURL,
		APIKey:  s.Config.TVSource.APIKey,
		PIN:     s.Config.TVSource.PIN,
	}

	sourcer := external.GetSource(&http.Client{}, source, "tvdb")
	resp, err := sourcer.Fetch(query)
	if err != nil {
		http.Error(w, "Failed to fetch media", http.StatusInternalServerError)
		return
	}

	var tvSearchResponse external.TVSearchResponse
	switch mediaType {
	case "tvdb":
		err = json.NewDecoder(resp.Body).Decode(&tvSearchResponse)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, "Failed to unmarshal response body", http.StatusInternalServerError)
			return
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
		http.Error(w, "Failed to convert response to media structs", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Inserting media into database: %v\n", medias)
	err = s.Repository.CreateMedia(medias)
	if err != nil {
		// this is a band aid for spelling changing between what the user entered and what TVDB returns
		// for exampel K Pop Demon Hunters returns as KPop Demon Hunters, which causes a duplicate key error when trying to insert into the database
		if strings.Contains(err.Error(), "duplicate key") {
			fmt.Printf("Media already exists in database: %v\n", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(medias)
			return
		}
		fmt.Printf("Failed to insert media into database: %v\n", err)
		http.Error(w, "Failed to insert media into database", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(medias)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(medias)
}
