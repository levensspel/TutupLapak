package httpServer

import (
	"github.com/gofiber/fiber/v2"
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

func (r *FileRepository) InsertURI(ctx *fiber.Ctx, fileUri string, thumbnailUri string) (*FileEntity, error) {
	entity := new(FileEntity)
	query := `
		INSERT INTO files (fileuri, thumbnailuri)
		VALUES($1, $2)
		RETURNING id, fileuri, thumbnailuri;
	`
	row := r.DB.QueryRow(ctx.Context(), query, fileUri, thumbnailUri)
	err := row.Scan(&entity.FileID, &entity.FileURI, &entity.ThumbnailURI)
	return entity, err
}
