package exceptions

type NotFoundError struct {
	Message    string
	StatusCode int16
}

func NewNotFoundError(message string, statusCode int16) *NotFoundError {
	return &NotFoundError{Message: message, StatusCode: statusCode}
}

func (e *NotFoundError) Error() string {
	return e.Message
}
