package helper

import (
	"strings"

	"github.com/jackc/pgx/v5"
)

func MapPgxError(err error) (statusCode int16, message string) {
	switch {
	case err == pgx.ErrNoRows:
		return 404, "Resource not found"
	}

	var sqlState string
	parts := strings.Split(err.Error(), "SQLSTATE ")
	if len(parts) > 1 {
		sqlState = parts[1][:5]
	}

	switch sqlState {
	case "23505": // Unique violation
		return 409, "Duplicate entry"
	case "23503": // Foreign key violation
		return 400, "Invalid reference"
	case "23514": // Check constraint violation
		return 400, "Constraint violation"
	default:
		return 500, "Internal server error"
	}
}
