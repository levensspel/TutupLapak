package config

import (
	"os"
	"strconv"
)

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

func GetRedisHost() string {
	return getEnv("REDIS_HOST", "127.0.0.1")
}

func GetRedisPort() string {
	return getEnv("REDIS_PORT", "6379")
}

func GetRedisPassword() string {
	return getEnv("REDIS_PASSWORD", "")
}

func GetRedisDbCount() int {
	count, err := strconv.Atoi(getEnv("REDIS_DB_COUNT", "0"))
	if err != nil {
		return 0
	}

	return count
}

func getFileServiceBaseURL() string {
	return getEnv("FILE_SERVICE_BASE_URL", "")
}

func getMode() string {
	return getEnv("MODE", "")
}
