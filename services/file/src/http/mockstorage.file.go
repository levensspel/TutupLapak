package httpServer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

type MockS3StorageClient struct{}

func NewMockS3StorageClient() StorageClient {
	return &MockS3StorageClient{}
}

func (m *MockS3StorageClient) PutFile(ctx context.Context, key string, mimeType string, fileContent []byte, isPublic bool) (string, error) {
	if strings.Contains(key, "mock_failed") {
		return "", errors.New("failed to put file")
	}

	// Tentukan lokasi penyimpanan file
	savePath := fmt.Sprintf("./.uploads/%s", key)

	// Buat direktori jika belum ada
	if err := os.MkdirAll("./.uploads", os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Simpan file ke direktori
	if err := os.WriteFile(savePath, fileContent, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Kembalikan URL file yang disimpan
	return m.GetUrl(key), nil
}

func (m *MockS3StorageClient) GetUrl(key string) string {
	// Buat URI untuk file
	baseURL := "/uploads"
	return fmt.Sprintf("%s/%s", baseURL, key)
}
