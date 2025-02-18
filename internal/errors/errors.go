package errors

import "fmt"

type AppErrorType int

type AppError struct {
	Code     AppErrorType
	Message  string
	Messages []string
}

var (
	errorNames    = make(map[AppErrorType]string)
	errorMessages = make(map[AppErrorType]string)
)

func (e AppError) Error() string {
	if len(e.Messages) > 0 {
		return fmt.Sprintf("%s: %v", e.Message, e.Messages)
	}
	return e.Message
}

func NewError(code AppErrorType, name string, message string) AppError {
	errorNames[code] = name
	errorMessages[code] = message
	return AppError{Code: code, Message: message}
}

func NewValidationError(messages []string) AppError {
	return AppError{
		Code:     ERR_VALIDATION_FAILED,
		Message:  "Validation failed",
		Messages: messages,
	}
}

func (e AppError) EnumName() string {
	if name, ok := errorNames[e.Code]; ok {
		return name
	}
	return "UNKNOWN_ERROR"
}

var (
	InternalServerError = NewError(ERR_INTERNAL_SERVER_ERROR, "ERR_INTERNAL_SERVER_ERROR", "internal server error")
	ValidationError     = NewError(ERR_VALIDATION_FAILED, "ERR_VALIDATION_FAILED", "validation failed")
)
