package storage

import (
	"errors"
)

type Client interface {
	Get(key string) ([]byte, error)
	GetObjectSize(key string) (int64, error)
	Put(key string, object string) error
	Delete(object string)
}

const (
	BUCKET_NAME string = "graphless-graph-file-bucket"
	BUCKET_KEY         = "graphFileKey"
)

//https://graphless-graph-file-bucket.s3.us-east-2.amazonaws.com/graphFileKey-dota-league-properties
//GetStorageClient returns an implementation of the StorageClient, used to access remote storage
func GetStorageClient(client ClientType) (Client, error) {
	switch client {
	case S3:
		return newS3StorageClient()
	case GoogleCloudStorage:
		return newGoogleCloudStorageClient()
	case LocalFileSystem:
		return newLocalFileSystemStorageClient()
	}

	return nil, errors.New("Unsupported storage client")
}
