package api

import (
	"net/http"

	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/service/session"
	"github.com/MatthewAraujo/airCast/internal/service/user"
	"github.com/MatthewAraujo/airCast/internal/service/video"
	"github.com/MatthewAraujo/airCast/internal/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *APIServer) loadRoutes() (http.Handler, error) {

	repo := repository.New(s.db)

	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	router.Use(c.Handler)

	// if the api changes in the future we can just change the version here, and the old version will still be available
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	userHandler := user.NewHandler(repo, s.logger)
	userHandler.RegisterRoutes(subrouter.PathPrefix("/user").Subrouter())
	s.logger.Info("user router up")

	sessionHandler := session.NewHandler(repo, s.logger)
	sessionHandler.RegisterRoutes(subrouter.PathPrefix("/video").Subrouter())

	subrouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, 200, map[string]string{"status": "api is healthy"})
	}).Methods(http.MethodGet)

	videoHandler := video.NewHandler(repo, s.logger)
	videoHandler.RegisterRoutes(subrouter)
	s.logger.Info("video router up")

	return subrouter, nil
}
