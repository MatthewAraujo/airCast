package errors

var (
	UserNotFound       = NewError(ERR_USER_NOT_FOUND, "ERR_USER_NOT_FOUND", "User not found")
	EmailAlreadyExists = NewError(ERR_EMAIL_ALREADY_EXISTS, "ERR_EMAIL_ALREADY_EXISTS", "Email already exists")
	InvalidCredentials = NewError(ERR_INVALID_CREDENTIALS, "ERR_INVALID_CREDENTIALS", "Invalid email or password")
)
