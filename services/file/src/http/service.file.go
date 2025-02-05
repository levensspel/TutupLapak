package httpServer

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"runtime"
	"strings"

	conf "github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/TimDebug/TutupLapak/File/src/models"
	"github.com/TimDebug/TutupLapak/File/src/repo"
	"github.com/disintegration/imaging"
	wpool "github.com/gammazero/workerpool"
	"github.com/gofiber/fiber/v2"
)

var (
	c *conf.Configuration = conf.GetConfig()
)

type FileService struct {
	Repo          repo.FileRepository
	StorageClient StorageClient
	wp            *wpool.WorkerPool
}

func NewFileService(repo repo.FileRepository, storageClient StorageClient) FileService {
	return FileService{
		Repo:          repo,
		StorageClient: storageClient,
		wp:            wpool.New(runtime.NumCPU() * 3),
	}
}

func (fs *FileService) UploadFile(
	ctx *fiber.Ctx,
	targetFilename string,
	file multipart.File,
	mimetype string,
) (*models.FileEntity, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}

	thumbFileName := fmt.Sprintf("thumbnail-%s", targetFilename)

	fs.wp.Submit(func() {
		fs.StorageClient.PutFile(ctx.Context(), targetFilename, mimetype, fileContent, true)
	})

	fs.wp.Submit(func() {
		fileBuf, err := fs.compressImage(fileContent)
		if err == nil {
			fs.StorageClient.PutFile(ctx.Context(), thumbFileName, mimetype, fileBuf, true)
		}
	})

	mainUri := fs.StorageClient.GetUrl(targetFilename)
	thumbnailUri := fs.StorageClient.GetUrl(thumbFileName)
	entity, err := fs.Repo.InsertURI(ctx.Context(), mainUri, thumbnailUri)
	if err != nil {
		return nil, &fiber.Error{Code: 400, Message: fmt.Sprintf("database insert failed: %w", err)}
	}
	return entity, nil
}

func (fs *FileService) compressImage(content []byte) ([]byte, error) {
	img, format, err := fs.decodeImage(content)
	logger.Logger.Info().Str("format", format).Msg("cek doang")
	if err != nil {
		return nil, err
	}
	resizedWidth := int(float64(img.Bounds().Dx()) * 0.1)
	resizedImg := imaging.Resize(img, resizedWidth, 0, imaging.Lanczos)
	return fs.imageToBytes(resizedImg, format)
}

func (fs *FileService) imageToBytes(img image.Image, fileExt string) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	logger.Logger.Info().Str("file_ext", fileExt).Msg("cek file eks")
	if strings.ToLower(fileExt) == "png" {
		err = imaging.Encode(&buf, img, imaging.PNG)
	} else {
		err = imaging.Encode(&buf, img, imaging.JPEG)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (fs *FileService) decodeImage(data []byte) (image.Image, string, error) {
	reader := bytes.NewReader(data)
	return image.Decode(reader)
}

func (fs *FileService) Shutdown() {
	fs.wp.StopWait()
}
