package exceptions

type UnauthorizedError struct {
	Message    string
	StatusCode int
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message, StatusCode: 401}
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
