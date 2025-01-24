package exceptions

type UnauthorizedError struct {
	Message    string
	StatusCode int16
}

func NewUnauthorizedError(message string, statusCode int16) *UnauthorizedError {
	return &UnauthorizedError{Message: message, StatusCode: statusCode}
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
