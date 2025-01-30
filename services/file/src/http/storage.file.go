package httpServer

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type StorageClient interface {
	// PutFile puts a file to the storage.
	// The key is the filename or path in the storage.
	// The mimeType is the content type of the file.
	// The fileContent is the actual content of the file.
	// The isPublic flag determines whether the file is publicly accessible.
	// It returns the URL of the uploaded file on success, or an error on failure.
	PutFile(
		ctx context.Context,
		key string,
		mimeType string,
		fileContent []byte,
		isPublic bool,
	) (string, error)

	// GetFileContent retrieves the content of a file from the storage.
	// The key is the filename or path in the storage.
	// It returns the content of the file on success, or an error on failure.
	// GetFileContent(ctx context.Context, key string) ([]byte, error)

	// GetUrl generates the complete URL for a file in the storage.
	// The key is the filename or path in the storage.
	// It returns the file's URL as a string.
	GetUrl(key string) string
}

var (
	s3StorageClientOnce     sync.Once
	s3StorageClientInstance *S3StorageClient
	appConfig               *config.Configuration = config.GetConfig()
)

type S3StorageClient struct {
	s3 *s3.Client
}

func NewAWS() aws.Config {
	sdkConfig, _ := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(appConfig.AWSRegion),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			appConfig.AWSAccessKey,
			appConfig.AWSSecretAccessKey,
			"",
		)),
	)
	return sdkConfig
}

func NewS3StorageClient() StorageClient {
	s3StorageClientOnce.Do(func() {
		sdkConfig := NewAWS()
		_s3 := s3.NewFromConfig(sdkConfig)
		s3StorageClientInstance = &S3StorageClient{
			s3: _s3,
		}
	})
	return s3StorageClientInstance
}

func (sc *S3StorageClient) PutFile(ctx context.Context, key string, mimeType string, fileContent []byte, isPublic bool) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(appConfig.AWSBucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(fileContent),
		ContentLength: aws.Int64(int64(len(fileContent))),
		ContentType:   aws.String(mimeType),
		ACL: func() types.ObjectCannedACL {
			if isPublic {
				return types.ObjectCannedACLPublicRead
			}
			return types.ObjectCannedACLPrivate
		}(),
	}
	_, err := sc.s3.PutObject(ctx, input)
	if err != nil {
		return "", err
	}
	return sc.GetUrl(key), nil
}

func (sc *S3StorageClient) GetUrl(key string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		appConfig.AWSBucket,
		appConfig.AWSRegion,
		key,
	)
}
