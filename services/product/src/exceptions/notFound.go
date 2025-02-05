package exceptions

type NotFoundError struct {
	Message    string
	StatusCode int
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message, StatusCode: 404}
}

func (e *NotFoundError) Error() string {
	return e.Message
}
