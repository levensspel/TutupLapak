package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port                string
	GRPCPort            string
	DBConnection        string
	DBConnectionMigrate string
	AutoMigrate         bool
	MigrateFileLocation string
	IsProduction        bool
	AWSAccessKey        string
	AWSSecretAccessKey  string
	AWSRegion           string
	AWSBucket           string
}

var (
	instance *Configuration
	once     sync.Once
)

func GetConfig() *Configuration {
	once.Do(func() {
		godotenv.Load()
		instance = &Configuration{
			Port:                getPort(),
			GRPCPort:            getEnv("GRPC_PORT", "5000"),
			DBConnection:        getDBConnection(),
			DBConnectionMigrate: getDBConnectionMigrate(),
			AutoMigrate:         getAutoMigrate(),
			MigrateFileLocation: getLocationMigrate(),
			IsProduction:        isProduction(),
			AWSAccessKey:        getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
			AWSRegion:           getEnv("AWS_REGION", ""),
			AWSBucket:           getEnv("AWS_BUCKET", ""),
		}
	})
	return instance
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getPort() string {
	return getEnv("PORT", "3000")
}

func getDBConnection() string {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "postgres")
	databaseurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	return databaseurl
}

func getDBConnectionMigrate() string {
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "postgres")
	databaseurl := fmt.Sprintf("pgx://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	return databaseurl
}

func getAutoMigrate() bool {
	return strings.ToUpper(getEnv("ENABLE_AUTO_MIGRATE", "FALSE")) == "TRUE"
}

func getLocationMigrate() string {
	return getEnv("MIGRATION_FILE_PATH", "file://src/database/migrations")
}

func isProduction() bool {
	return strings.ToUpper(getEnv("MODE", "DEBUG")) == "PRODUCTION"
}
