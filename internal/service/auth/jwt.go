package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	configs "github.com/MatthewAraujo/airCast/internal/config"
	"github.com/MatthewAraujo/airCast/internal/errors"
	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const USER_KEY = "userID"

func WithJWTAuth(handler http.HandlerFunc, store repository.Queries, ctx context.Context, logger slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil || !token.Valid {
			logger.Error("invalid token", "error", err)
			permissionDenied(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("failed to parse token claims")
			permissionDenied(w)
			return
		}

		userIDStr, ok := claims[USER_KEY].(string)
		if !ok {
			logger.Error("userID claim missing or invalid")
			permissionDenied(w)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.Error("invalid userID format", "userID", userIDStr, "error", err)
			permissionDenied(w)
			return
		}

		user, err := store.FindUserByID(ctx, userID)
		if err != nil {
			logger.Error("user not found", "userID", userID, "error", err)
			permissionDenied(w)
			return
		}

		ctx = context.WithValue(r.Context(), USER_KEY, user)
		handler(w, r.WithContext(ctx))
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth == "" {
		return ""
	}

	parts := strings.Split(tokenAuth, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}

	return ""
}

func CreateJWT(userID string) (string, error) {
	secret := []byte(configs.Envs.JWT.JWTSecret)
	expiration := time.Second * time.Duration(configs.Envs.JWT.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		USER_KEY:  userID,
		"expires": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWT.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, errors.AppError{
		Code:    errors.ERR_UNAUTHORIZED,
		Message: "permission denied",
	})
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(USER_KEY).(repository.User)
	if !ok {
		return uuid.Nil
	}

	return userID.ID
}
