package errors

type AppErrorType int

type AppError struct {
	Code    AppErrorType
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

var (
	errorMessages = map[AppErrorType]string{}
	errorNames    = map[AppErrorType]string{}
)

func NewError(code AppErrorType, name string, message string) AppError {
	errorMessages[code] = message
	errorNames[code] = name
	return AppError{Code: code, Message: message}
}

func (e AppError) EnumName() string {
	if name, ok := errorNames[e.Code]; ok {
		return name
	}
	return "UnknownError"
}
