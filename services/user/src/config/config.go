package config

import "os"

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

func GetGRPCPort() string {
	return getEnv("GRPC_PORT", "50051")
}

var FILE_SERVICE_BASE_URL string
