package repo

import (
	"context"

	"github.com/TimDebug/TutupLapak/File/src/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	DB *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) FileRepository {
	return FileRepository{
		DB: db,
	}
}

func (r *FileRepository) InsertURI(ctx context.Context, fileUri string, thumbnailUri string) (*models.FileEntity, error) {
	entity := new(models.FileEntity)
	query := `
		INSERT INTO files (fileuri, thumbnailuri)
		VALUES($1, $2)
		RETURNING id, fileuri, thumbnailuri;
	`
	row := r.DB.QueryRow(ctx, query, fileUri, thumbnailUri)
	err := row.Scan(&entity.FileID, &entity.FileURI, &entity.ThumbnailURI)
	return entity, err
}

func (r *FileRepository) GetRecordsById(ctx context.Context, fileId string) (*models.FileEntity, error) {
	query := `
		SELECT id, fileuri, thumbnailuri
		FROM files
		WHERE id = $1;
	`
	row := r.DB.QueryRow(ctx, query, fileId)
	entity := new(models.FileEntity)
	err := row.Scan(&entity.FileID, &entity.FileURI, &entity.ThumbnailURI)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
