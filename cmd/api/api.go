package api

import (
	"log"
	"net/http"

	"github.com/MatthewAraujo/airCast/service/video"
	"github.com/MatthewAraujo/airCast/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
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

	videoHandler := video.NewHandler()
	videoHandler.RegisterRoutes(subrouter)

	subrouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, 200, map[string]string{"status": "api is healthy"})
	})

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
