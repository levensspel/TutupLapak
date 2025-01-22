package config

import (
	"fmt"
	"os"
	"strings"
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

func GetDBConnection() string {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "postgres")

	databaseurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	return databaseurl

}

func GetDBConnectionMigrate() string {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "postgres")

	databaseurl := fmt.Sprintf("pgx://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	return databaseurl
}

func GetAutoMigrate() bool {
	return strings.ToUpper(getEnv("ENABLE_AUTO_MIGRATE", "FALSE")) == "TRUE"
}

func GetLocationMigrate() string {
	return getEnv("MIGRATION_FILE_PATH", "file://src/database/migrations")
}
