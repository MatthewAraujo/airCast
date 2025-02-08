package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MatthewAraujo/airCast/internal/errors"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err errors.AppError) {
	WriteJSON(w, status, map[string]any{
		"status":  "error",
		"message": err.Message,
		"enum":    err.EnumName(),
	})
}

func WriteSuccess(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, map[string]any{
		"status": "success",
		"data":   data,
	})
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10) // 10 é a base decimal
}
