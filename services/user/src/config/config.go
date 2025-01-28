package config

import "os"

var FILE_SERVICE_BASE_URL string
var MODE string

const (
	MODE_DEBUG      string = "DEBUG"
	MODE_PRODUCTION        = "PRODUCTION"
)

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetPort() string {
	return getEnv("PORT", "8080")
}

func GetFileServiceBaseURL() string {
	return getEnv("FILE_SERVICE_BASE_URL", "")
}

func GetMode() string {
	return getEnv("MODE", "")
}
