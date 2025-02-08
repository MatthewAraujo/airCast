package user

import (
	"log/slog"
	"net/http"

	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/gorilla/mux"
)

type Handler struct {
	db     *repository.Queries
	logger *slog.Logger
}

func NewHandler(db *repository.Queries, logger *slog.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/account", h.registerAccount).Methods(http.MethodPost)
}

func (h *Handler) registerAccount(w http.ResponseWriter, r *http.Request) {

}
