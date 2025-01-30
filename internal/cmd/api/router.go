package api

import (
	"net/http"

	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/service/video"
	"github.com/MatthewAraujo/airCast/internal/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *APIServer) loadRoutes() (http.Handler, error) {

	repo := repository.New(s.db)

	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	router.Use(c.Handler)
	// if the api changes in the future we can just change the version here, and the old version will still be available
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	videoHandler := video.NewHandler(repo)
	videoHandler.RegisterRoutes(subrouter)

	subrouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, 200, map[string]string{"status": "api is healthy"})
	})

	return subrouter, nil
}
