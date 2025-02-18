package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	internal_error "github.com/MatthewAraujo/airCast/internal/errors"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func WriteSuccess(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, map[string]any{
		"status": "success",
		"data":   data,
	})
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10) // 10 é a base decimal
}

func TranslateValidationErrors(errs validator.ValidationErrors) []string {
	var messages []string
	fmt.Printf("errs: %v\n", errs)
	for _, err := range errs {
		message := fmt.Sprintf("The field '%s' failed on validation: %s", err.Field(), err.Tag())
		if err.Tag() == "required" {
			message = fmt.Sprintf("The field '%s' is required.", err.Field())
		} else if err.Tag() == "oneof" {
			message = fmt.Sprintf("The field '%s' must be one of these values: %s.", err.Field(), err.Param())
		}
		messages = append(messages, message)
	}
	return messages
}

func WriteError(w http.ResponseWriter, status int, err internal_error.AppError) {
	response := map[string]any{
		"status":  "error",
		"message": err.Message,
		"enum":    err.EnumName(),
	}

	if len(err.Messages) > 0 {
		response["messages"] = err.Messages
	}

	WriteJSON(w, status, response)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	return json.NewDecoder(r.Body).Decode(payload)

}
