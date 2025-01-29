package models

type FileEntity struct {
	FileID       string `json:"fileId"`
	FileURI      string `json:"fileUri"`
	ThumbnailURI string `json:"fileThumbnailUri"`
}
