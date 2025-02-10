package user

import (
	"log/slog"
	"net/http"

	"github.com/MatthewAraujo/airCast/internal/errors"
	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/types"
	"github.com/MatthewAraujo/airCast/internal/utils"
	"github.com/go-playground/validator/v10"
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
	router.HandleFunc("/login", h.login).Methods(http.MethodPost)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) registerAccount(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := utils.TranslateValidationErrors(validationErrors)
			utils.WriteError(w, http.StatusBadRequest, errors.NewValidationError(errorMessages))
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, errors.NewError(errors.ERR_VALIDATION_FAILED, "ERR_UNKNOWN_VALIDATION", "Unexpected validation error"))
		return
	}
}
