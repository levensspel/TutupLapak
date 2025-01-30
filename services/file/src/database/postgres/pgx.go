package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxConnect() (*pgxpool.Pool, error) {
	DbString := config.GetConfig().DBConnection
	fmt.Printf("Database connection string: %s\n", DbString)

	// Create connection pool
	db, err := pgxpool.New(context.Background(), DbString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Get max_connections (optional)
	var maxConStr string
	err = db.QueryRow(context.Background(), "SHOW max_connections").Scan(&maxConStr)
	if err != nil {
		log.Printf("Warning: Failed to retrieve max_connections: %v", err)
	} else {
		maxConn, err := strconv.Atoi(maxConStr)
		if err == nil {
			db.Config().MaxConns = int32(float64(maxConn) * 0.8)
			fmt.Printf("Setting max connections to: %d\n", db.Config().MaxConns)
		}
	}

	return db, nil
}
