package config

import "os"

var FILE_SERVICE_BASE_URL string
var MODE string

const (
	MODE_DEBUG      string = "DEBUG"
	MODE_PRODUCTION        = "PRODUCTION"
)

func SetupReusableEnv() {
	FILE_SERVICE_BASE_URL = getFileServiceBaseURL()
	if FILE_SERVICE_BASE_URL == "" {
		panic("FILE_SERVICE_BASE_URL value requires to be set")
	}

	MODE = getMode()
	if MODE == "" {
		panic("MODE value requires to be set")
	}
}

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

func GetGRPCPort() string {
	return getEnv("GRPC_PORT", "50051")
}

func GetMetricGRPCPort() string {
	return getEnv("METRIC_GRPC_PORT", "9090")
}

func getFileServiceBaseURL() string {
	return getEnv("FILE_SERVICE_BASE_URL", "")
}

func getMode() string {
	return getEnv("MODE", "")
}
