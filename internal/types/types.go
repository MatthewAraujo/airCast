package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type ValidationErrorResponse struct {
	Field      string `json:"field"`
	Validation string `json:"validation"`
	Value      string `json:"value,omitempty"`
	Message    string `json:"message"`
}

type Session struct {
	ID           string   `json:"id"`
	HostID       string   `json:"host_id"`
	VideoId      string   `json:"video_id"`
	Participants []string `json:"participants"`
}

type CreateSessionPayload struct {
	VideoId string `json:"video_id" validate:"required"`
}

type JoinSessionPayload struct {
	SessionId string `json:"session_id" validate:"required"`
}
