package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	configs "github.com/MatthewAraujo/airCast/internal/config"
	"github.com/MatthewAraujo/airCast/internal/errors"
	"github.com/MatthewAraujo/airCast/internal/repository"
	"github.com/MatthewAraujo/airCast/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handleFunc http.HandlerFunc, store repository.Queries, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("token is invalid")
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userID, err := uuid.Parse(str)
		if err != nil {
			log.Printf("error parsing userID: %v", err)
			permissionDenied(w)
			return
		}

		u, err := store.FindUserByID(ctx, userID)
		if err != nil {
			log.Printf("error fetching user: %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handleFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth == "" {
		return tokenAuth
	}
	return ""
}

func CreateJWT(userID string) (string, error) {
	secret := []byte(configs.Envs.JWT.JWTSecret)
	expiration := time.Second * time.Duration(configs.Envs.JWT.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":  userID,
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
	utils.WriteError(w, http.StatusForbidden, errors.AppError{})
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userID
}
