package clients

import (
	"context"
	"errors"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageClient interface {
	Get(key string) ([]byte, error)
	GetObjectSize(key string) (int64, error)
	Put(key string, object string) error
	Delete(object string)
}

const (
	BUCKET_NAME       string = "graphless-graph-file-bucket"
	BUCKET_KEY               = "graphFileKey"
	GOOGLE_PROJECT_ID        = "granular-graph-processing"
)

//GetStorageClient returns an implementation of the StorageClient, used to access remote storage
func GetStorageClient(client StorageClientType) (StorageClient, error) {
	switch client {
	case S3:
		sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-2")}))
		s3Client := s3.New(sess)
		return &S3Client{
			bucketName: BUCKET_NAME,
			bucketKey:  BUCKET_KEY,
			uploadManager: s3manager.NewUploaderWithClient(s3Client, func(d *s3manager.Uploader) {
				d.Concurrency = 2
			}),
			downloadManager: s3manager.NewDownloaderWithClient(s3Client, func(d *s3manager.Downloader) {
				d.Concurrency = 2
			}),
			S3: s3Client,
		}, nil
	case GOOGLE_CLOUD_STORAGE:
		ctx := context.Background()
		cloudStorageClient, err := storage.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		bucket := cloudStorageClient.Bucket(BUCKET_NAME)
		bucket.If(storage.BucketConditions{MetagenerationMatch: int64(0)}).Create(ctx, GOOGLE_PROJECT_ID, nil)
		if err != nil {
			return nil, err
		}

		return &GoogleCloudStorageClient{
			bucketName: BUCKET_NAME,
			bucketKey:  BUCKET_KEY,
			client:     cloudStorageClient,
			context:    ctx,
		}, nil

	}

	return nil, errors.New("Unsupported storage client")
}
