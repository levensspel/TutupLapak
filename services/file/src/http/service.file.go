package httpServer

import "github.com/jackc/pgx/v5/pgxpool"

type FileService struct {
	Repo          FileRepository
	DB            *pgxpool.Pool
	StorageClient StorageClient
}

func NewFileService(repo FileRepository, db *pgxpool.Pool, storageClient StorageClient) FileService {
	return FileService{
		Repo:          repo,
		DB:            db,
		StorageClient: storageClient,
	}
}
