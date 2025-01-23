package exceptions

type ErrorResponse struct {
	StatusCode int16
	Message    string
}

func NewErrorResponse(statusCode int16, message string) ErrorResponse {
	return ErrorResponse{StatusCode: statusCode, Message: message}
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func ErrBadRequest(message string) error {
	return NewErrorResponse(400, message)
}

func ErrUnauthorized(message string) error {
	return NewErrorResponse(401, message)
}

func ErrNotFound(message string) error {
	return NewErrorResponse(404, message)
}

func ErrConflict(message string) error {
	return NewErrorResponse(409, message)
}

func ErrServer(message string) error {
	return NewErrorResponse(500, message)
}
