package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/axschech/rockbot-backend/external"
	"github.com/axschech/rockbot-backend/internal/config"
	"github.com/axschech/rockbot-backend/internal/database/repository"
	"github.com/axschech/rockbot-backend/internal/entities"
	"github.com/axschech/rockbot-backend/internal/media"
	"github.com/axschech/rockbot-backend/internal/media_user"
	"github.com/axschech/rockbot-backend/internal/routing"
	"github.com/axschech/rockbot-backend/internal/user"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Config     config.Config
	Repository repository.Repository
	Router     routing.Router
	HTTPClient *http.Client
}

func NewService(
	cfg config.Config,
	r repository.Repository,
	router routing.Router,
	httpClient *http.Client,
) *Service {
	return &Service{
		Config:     cfg,
		Repository: r,
		Router:     router,
		HTTPClient: httpClient,
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func (s *Service) Run() error {
	s.Router.R.Get("/ping", Ping)
	s.Router.R.Route("/api", func(r chi.Router) {
		r.Get("/user", s.GetUserHandler)
		r.Post("/user", s.PostUserHandler)
		r.Get("/media/user/{id}", s.GetMediaUsersWithUserIDHandler)
		r.Post("/media/user", s.PostMediaUserHandler)
		r.Put("/media/user", s.PutMediaUserHandler)
		r.Get("/search/media", s.QueryMediaHandler)
	})

	return s.Router.Listen()
}

// handlers should probably be their own structs, with an interface called Handlerer
func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	u := user.NewUser(s.Repository)
	userId := r.URL.Query().Get("id")
	email := r.URL.Query().Get("email")

	if userId == "" && email == "" {
		http.Error(w, "user ID or email is required", http.StatusBadRequest)
		return
	}

	if email != "" {
		user, err := u.GetUserByEmail(email)
		if err != nil {
			fmt.Printf("error getting user by email: %+v\n", err)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Printf("error converting user ID to int: %+v\n", err)
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := u.GetUserByID(id)

	if err != nil {
		fmt.Printf("error getting user by ID: %+v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	u := user.NewUser(s.Repository)

	var req PostUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %+v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request to create user: %+v\n", req)

	if req.Username == "" || req.Email == "" {
		http.Error(w, "username and email are required", http.StatusBadRequest)
		return
	}

	user, err := u.Register(req.Username, req.Email)
	if err != nil {
		fmt.Printf("Error creating user: %+v\n", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) QueryMediaHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	t := r.URL.Query().Get("type")
	if t == "" {
		http.Error(w, "type parameter is required", http.StatusBadRequest)
		return
	}

	var mediaType string
	switch t {
	case "tv":
		mediaType = "tvdb"
	}

	if mediaType == "" {
		http.Error(w, "invalid media type", http.StatusBadRequest)
		return
	}

	source := config.Source{
		BaseURL: s.Config.TVSource.BaseURL,
		APIKey:  s.Config.TVSource.APIKey,
		PIN:     s.Config.TVSource.PIN,
	}

	sourcer := external.GetSource(s.HTTPClient, source, mediaType)

	nm := media.NewMedia(s.Repository, sourcer)

	medias, err := nm.GetOrSaveMedia(query, mediaType)
	if err != nil {
		fmt.Printf("Failed to get or save media: %+v\n", err)
		http.Error(w, "failed to get or save media", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(medias)
}

func (s *Service) GetMediaUsersWithUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")

	id, err := strconv.Atoi(userId)
	if err != nil {
		fmt.Printf("error converting user ID to int: %+v\n", err)
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user := user.NewUser(s.Repository)
	// not sure if these checks should be in the handler or business logic
	_, err = user.GetUserByID(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	withMedia := r.URL.Query().Get("with_media") == "true"
	fmt.Printf("Received request to get media users for user ID %d with media: %t\n", id, withMedia)
	um := media_user.NewMediaUser(s.Repository)

	mediaUsersWithMedia, err := um.QueryMediaUsersWithUserID(id, withMedia)
	if err != nil {
		fmt.Printf("Failed to query media users: %+v\n", err)
		http.Error(w, "failed to query media users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mediaUsersWithMedia)
}

func (s *Service) PostMediaUserHandler(w http.ResponseWriter, r *http.Request) {
	var req PostMediaUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %+v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	mediaUser := entities.MediaUserEntity{
		UserID:  req.UserID,
		MediaID: req.MediaID,
		Status:  req.Status,
	}

	user := user.NewUser(s.Repository)

	_, err := user.GetUserByID(mediaUser.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	um := media_user.NewMediaUser(s.Repository)

	createdMediaUser, err := um.CreateMediaUser(mediaUser)
	if err != nil {
		http.Error(w, "failed to create media user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMediaUser)
}

func (s *Service) PutMediaUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to update media user")
	// might not need this, could allow continous insert and just get the last one, using this for now
	var req PutMediaUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %+v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.MediaUserID == 0 {
		http.Error(w, "media user ID is required", http.StatusBadRequest)
		return
	}

	allowedStatuses := []string{"not watched", "will watch", "watching", "watched", "wont watch"}
	status := req.Status
	if status == "" || !slices.Contains(allowedStatuses, status) {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	mediaUser := entities.MediaUserEntity{
		ID:     req.MediaUserID,
		Status: req.Status,
	}

	um := media_user.NewMediaUser(s.Repository)

	updatedMediaUser, err := um.UpdateMediaUser(mediaUser)
	if err != nil {
		fmt.Printf("Failed to update media user: %+v\n", err)
		http.Error(w, "failed to update media user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMediaUser)
}
