package grpc

import (
	"context"
	"fmt"

	"github.com/TimDebug/TutupLapak/File/src/grpc/proto/model/file"
	"github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/TimDebug/TutupLapak/File/src/repo"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fileService struct {
	file.UnimplementedFileServiceServer
	repo *repo.FileRepository
}

func NewFileService(repo *repo.FileRepository) file.FileServiceServer {
	return &fileService{repo: repo}
}

func (f *fileService) CheckExist(ctx context.Context, req *file.FileRequest) (*file.FileResponse, error) {
	fileId := req.FileId
	if fileId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid fileId")
	}
	entity, err := f.repo.GetRecordsById(ctx, fileId)
	if err != nil {
		logger.Logger.Error().Err(err).Msg(fmt.Sprintf("%+v", err))
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "records not found")
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return nil, status.Errorf(codes.InvalidArgument, "invalid fileId")
			}
		}
		return nil, status.Errorf(codes.Internal, "server error")
	}
	grpcEntity := file.FileResponse{
		FileId:       entity.FileID,
		FileUri:      entity.FileURI,
		ThumbnailUri: entity.ThumbnailURI,
	}
	return &grpcEntity, nil
}
