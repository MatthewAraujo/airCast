package session

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/MatthewAraujo/airCast/internal/errors"
	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/service/auth"
	"github.com/MatthewAraujo/airCast/internal/types"
	"github.com/MatthewAraujo/airCast/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	db     *repository.Queries
	logger *slog.Logger

	sessions map[string]*types.Session
	mu       sync.RWMutex
}

func NewHandler(db *repository.Queries, logger *slog.Logger) *Handler {
	return &Handler{
		db:       db,
		logger:   logger,
		sessions: make(map[string]*types.Session),
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		auth.WithJWTAuth(h.createSession, *h.db, r.Context(), *h.logger)(w, r)
	}).Methods(http.MethodPost)

	router.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		auth.WithJWTAuth(h.joinSession, *h.db, r.Context(), *h.logger)(w, r)
	}).Methods(http.MethodPatch)

	router.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		auth.WithJWTAuth(h.getSession, *h.db, r.Context(), *h.logger)(w, r)
	}).Methods(http.MethodGet)

}

func (h *Handler) getSession(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	utils.WriteJSON(
		w, http.StatusAccepted, h.sessions)
}

func (h *Handler) joinSession(w http.ResponseWriter, r *http.Request) {

	h.mu.Lock()
	defer h.mu.Unlock()

	var payload types.JoinSessionPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := utils.TranslateValidationErrors(validationErrors)
			utils.WriteError(w, http.StatusBadRequest, errors.NewValidationError(errorMessages))
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, errors.NewError(errors.ERR_VALIDATION_FAILED, "ERR_UNKNOWN_VALIDATION", "Unexpected validation error"))
		return
	}

	session, exists := h.sessions[payload.SessionId]
	if !exists {
		utils.WriteError(w, http.StatusNotFound, errors.SessionNotFound)
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())

	if userID == uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.Unauthorized)
	}

	session.Participants = append(session.Participants, userID.String())

	utils.WriteJSON(w, http.StatusOK, nil)
}
func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var payload types.CreateSessionPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := utils.TranslateValidationErrors(validationErrors)
			utils.WriteError(w, http.StatusBadRequest, errors.NewValidationError(errorMessages))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, errors.NewError(errors.ERR_VALIDATION_FAILED, "ERR_UNKNOWN_VALIDATION", "Unexpected validation error"))
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())

	if userID == uuid.Nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.Unauthorized)
		return
	}

	session := &types.Session{
		ID:           h.generateSessionID(),
		HostID:       userID.String(),
		VideoId:      payload.VideoId,
		Participants: []string{userID.String()},
	}

	h.sessions[session.ID] = session

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) generateSessionID() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		h.logger.Error("error creating random UUID", "error", err)
	}
	return "sess-" + uuid.String()
}
