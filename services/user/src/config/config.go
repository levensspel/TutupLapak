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

// Percentage use of max database connections.
// Default to 0.8 (80%).
// Range: 0.0 < x <= 10.
func GetDbMaxConnPercentage() float64 {
	percentage, err := strconv.ParseFloat(getEnv("DB_MAX_CONN_PERCENTAGE", "0.8"), 64)
	if err != nil || percentage <= 0.0 || percentage > 1.0 {
		panic("Invalid DB_MAX_CONN_PERCENTAGE: 0.0 < value <= 1.0")
	}

	return percentage
}

func getFileServiceBaseURL() string {
	return getEnv("FILE_SERVICE_BASE_URL", "")
}

func getMode() string {
	return getEnv("MODE", "")
}
