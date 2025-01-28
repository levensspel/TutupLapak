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
	return getEnv("PORT", "3000")
}

func GetUserGRPCHost() string {
	return getEnv("GPRC_USER_HOST", "localhost")
}

func GetUserGRPCPort() string {
	return getEnv("GPRC_USER_PORT", "50050")
}
