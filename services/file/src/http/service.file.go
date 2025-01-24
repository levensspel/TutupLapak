package httpServer

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
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

func (fs *FileService) UploadFile(ctx *fiber.Ctx, fileName string, file multipart.File, mimetype string) (*FileEntity, error) {
	var entity *FileEntity
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}
	mainUri, err := fs.StorageClient.PutFile(ctx.Context(), fileName, mimetype, fileContent, true)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}
	// compress file into thumbnail
	fileBuf, err := fs.compressImage(file)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to compress the file"}
	}
	// upload the thumbnail to S3
	thumbFileName := fmt.Sprintf("%s-%s", "thumbnail", fileName)
	thumbnailUri, err := fs.StorageClient.PutFile(ctx.Context(), thumbFileName, mimetype, fileBuf.Bytes(), true)
	if err != nil {
		return nil, &fiber.Error{Code: 500, Message: "unable to read file content"}
	}
	entity.FileURI = mainUri
	entity.ThumbnailURI = thumbnailUri
	// TODO: store the record to database
	return entity, nil
}

func (fs *FileService) compressImage(inputFile multipart.File) (*bytes.Buffer, error) {
	// Create a buffer to store the compressed output
	var outputBuffer bytes.Buffer

	// Create an input buffer to store the file content
	inputBuffer := &bytes.Buffer{}
	_, err := io.Copy(inputBuffer, inputFile)
	if err != nil {
		return nil, err
	}

	// Use FFmpeg to compress the image
	err = ffmpeg_go.
		Input("pipe:0"). // Specify that input comes from a buffer
		Output("pipe:1", ffmpeg_go.KwArgs{
			"q:v": "5", // Set quality; adjust as needed
		}).
		WithInput(inputBuffer).    // Provide the input buffer
		WithOutput(&outputBuffer). // Write output to the output buffer
		OverWriteOutput().
		Run()

	if err != nil {
		return nil, err
	}

	// Return the compressed image as a buffer
	return &outputBuffer, nil
}
