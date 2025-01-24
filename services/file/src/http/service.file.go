package httpServer

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"strings"

	"github.com/TimDebug/TutupLapak/File/src/logger"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

type FileEntity struct {
	FileID       string `json:"fileId"`
	FileURI      string `json:"fileUri"`
	ThumbnailURI string `json:"fileThumbnailUri"`
}

type FileService struct {
	Repo          FileRepository
	StorageClient StorageClient
}

func NewFileService(repo FileRepository, storageClient StorageClient) FileService {
	return FileService{
		Repo:          repo,
		StorageClient: storageClient,
	}
}

func (fs *FileService) UploadFile(
	ctx *fiber.Ctx,
	originalFilename string,
	targetFilename string,
	file multipart.File,
	mimetype string,
) (*FileEntity, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}
	mainUri, err := fs.StorageClient.PutFile(ctx.Context(), targetFilename, mimetype, fileContent, true)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}
	// compress file into thumbnail
	fileBuf, err := fs.compressImage(fileContent)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: fmt.Sprintf("unable to compress file: %+v", err)}
	}
	// upload the thumbnail to S3
	thumbFileName := fmt.Sprintf("%s-%s", "thumbnail", targetFilename)
	thumbnailUri, err := fs.StorageClient.PutFile(ctx.Context(), thumbFileName, mimetype, fileBuf, true)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}
	// store the record
	return fs.Repo.InsertURI(ctx, mainUri, thumbnailUri)
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
