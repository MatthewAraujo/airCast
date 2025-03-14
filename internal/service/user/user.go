package user

import (
	"database/sql"
	internal_error "errors"
	"log/slog"
	"net/http"

	"github.com/MatthewAraujo/airCast/internal/errors"
	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/service/auth"
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
	router.HandleFunc("/register", h.registerAccount).Methods(http.MethodPost)
	router.HandleFunc("/login", h.login).Methods(http.MethodPost)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

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

	ctx := r.Context()

	u, err := h.db.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.WriteError(w, http.StatusBadRequest, errors.InternalServerError)
			return
		}
	}

	if u.Email != payload.Email {
		utils.WriteError(w, http.StatusBadRequest, errors.UserNotFound)
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, errors.InvalidCredentials)
		return
	}

	token, err := auth.CreateJWT(u.ID.String())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.InternalServerError)
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, map[string]string{
		"token": token,
	})
}

func (h *Handler) registerAccount(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Debug("validate payload")
	if err := utils.Validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := utils.TranslateValidationErrors(validationErrors)
			utils.WriteError(w, http.StatusBadRequest, errors.NewValidationError(errorMessages))
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, errors.NewError(errors.ERR_VALIDATION_FAILED, "ERR_UNKNOWN_VALIDATION", "Unexpected validation error"))
		return
	}
	ctx := r.Context()

	h.logger.Debug("checking if email exists")
	emailAlreadyExists, err := h.db.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		if internal_error.Is(err, sql.ErrNoRows) {
			h.logger.Debug("email does not exist", "email", payload.Email)
		} else {
			h.logger.Error("database error", "error", err)
			utils.WriteError(w, http.StatusInternalServerError, errors.InternalServerError)
			return
		}
	}

	if emailAlreadyExists.Email != "" {
		h.logger.Debug("email exists")
		utils.WriteError(w, http.StatusBadRequest, errors.EmailAlreadyExists)
		return
	}

	h.logger.Debug("hashing password")
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errors.NewError(errors.ERR_VALIDATION_FAILED, "ERR_UNKNOWN_VALIDATION", "Unexpected validation error"))
		return
	}

	h.logger.Debug("inserting user")

	_, err = h.db.InsertUsers(ctx,
		repository.InsertUsersParams{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: hashedPassword,
		})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.EmailAlreadyExists)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, nil)
}
