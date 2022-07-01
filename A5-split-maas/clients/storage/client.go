package storage

import (
	"errors"
	"github.com/devLucian93/thesis-go/domain"
)

type Client interface {
	Get(key string) ([]byte, error)
	GetObjectSize(key string) (int64, error)
	Put(key string, object string) error
	Delete(object string)
}

//https://graphless-graph-file-bucket.s3.us-east-2.amazonaws.com/graphFileKey-dota-league-properties
//GetStorageClient returns an implementation of the StorageClient, used to access remote storage
func GetStorageClient(client ClientType, storageConfig domain.StorageConfig) (Client, error) {
	switch client {
	case S3:
		return newS3StorageClient(storageConfig)
	case GoogleCloudStorage:
		return newGoogleCloudStorageClient(storageConfig)
	case LocalFileSystem:
		return newLocalFileSystemStorageClient(storageConfig)
	}

	return nil, errors.New("Unsupported storage client")
}
