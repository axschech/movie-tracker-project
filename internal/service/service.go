package service

import (
	"encoding/json"
	"net/http"

	"github.com/axschech/rockbot-backend/external"
	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/media"
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

	sourcer := external.GetSource(&http.Client{}, source, mediaType)

	nm := media.NewMedia(s.Repository, sourcer)

	medias, err := nm.GetOrSaveMedia(query, mediaType)
	if err != nil {
		http.Error(w, "Failed to get or save media", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(medias)
}
