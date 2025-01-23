package httpServer

import "github.com/jackc/pgx/v5/pgxpool"

type FileRepository struct {
	DB *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) FileRepository {
	return FileRepository{
		DB: db,
	}
}
