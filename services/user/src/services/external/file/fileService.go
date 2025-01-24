package fileService

import (
	"encoding/json"
	"fmt"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/config"
	functionCallerInfo "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/helper"
	loggerZap "github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/logger/zap"
	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/service"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
)

type FileServiceInterface interface {
	// Call to external file service
	GetFile(ctx *fiber.Ctx, fileId string) (file service.File, statusCode int)
}

type fileService struct {
	Logger loggerZap.LoggerInterface
}

func NewFileService(logger loggerZap.LoggerInterface) FileServiceInterface {
	return &fileService{Logger: logger}
}

func NewFileServiceInject(i do.Injector) (FileServiceInterface, error) {
	_logger := do.MustInvoke[loggerZap.LoggerInterface](i)

	return NewFileService(_logger), nil
}

func (fs *fileService) GetFile(ctx *fiber.Ctx, fileId string) (file service.File, statusCode int) {
	url := fmt.Sprintf("%s/v1/file/%s", config.FILE_SERVICE_BASE_URL, fileId)
	agent := fiber.Get(url)
	statusCode, body, errs := agent.Bytes()
	errMsg := fmt.Sprintf("status code [%v] ", statusCode)
	if len(errs) > 0 {
		for i, err := range errs {
			errMsg += fmt.Sprintf("Errs[%v]: %s || \n", i, err.Error())
		}
		fs.Logger.Error(errMsg, functionCallerInfo.ExternalFileServiceGetFile)
		return service.File{}, statusCode
	}

	err := json.Unmarshal(body, &file)
	if err != nil {
		fs.Logger.Error(err.Error(), functionCallerInfo.ExternalFileServiceGetFile, statusCode, errMsg)
		return service.File{}, fiber.StatusInternalServerError
	}

	return file, statusCode
}
