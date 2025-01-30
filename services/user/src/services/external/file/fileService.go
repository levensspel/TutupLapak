package fileService

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/service"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/services/external/grpc/file"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileServiceInterface interface {
	// Call to external file service
	GetFile(ctx context.Context, fileId string) (response service.File, ok bool)
}

type fileService struct {
	GrpcClient file.FileServiceClient
	Logger     loggerZap.LoggerInterface
}

func NewFileService(grpcClient file.FileServiceClient, logger loggerZap.LoggerInterface) FileServiceInterface {
	return &fileService{
		GrpcClient: grpcClient,
		Logger:     logger,
	}
}

func NewFileServiceInject(i do.Injector) (FileServiceInterface, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	conn, err := grpc.NewClient(
		config.FILE_SERVICE_BASE_URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := file.NewFileServiceClient(conn)

	return NewFileService(client, _logger), nil
}

func (fs *fileService) GetFile(ctx context.Context, fileId string) (service.File, bool) {
	response, err := fs.GrpcClient.CheckExist(ctx, &file.FileRequest{FileId: fileId})
	if err != nil {
		fs.Logger.Error(err.Error(), functionCallerInfo.ExternalFileServiceGetFile)
		return service.File{}, false
	}

	return service.File{
		FileUri:          response.FileUri,
		FileThumbnailUri: response.ThumbnailUri,
	}, true
}
