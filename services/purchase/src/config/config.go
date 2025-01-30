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
	return getEnv("PORT", "3000")
}

func GetMode() string {
	return getEnv("MODE", "DEBUG")
}

func GetUserGRPCHost() string {
	return getEnv("GPRC_USER_HOST", "localhost")
}

func GetUserGRPCPort() string {
	return getEnv("GPRC_USER_PORT", "50050")
}
