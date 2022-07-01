package storage

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

const (
	GOOGLE_PROJECT_ID = "granular-graph-processing"
)

type googleCloudStorageClient struct {
	bucketName string
	bucketKey  string
	client     *storage.Client
	context    context.Context
}

func newGoogleCloudStorageClient() (Client, error) {
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

	return &googleCloudStorageClient{
		bucketName: BUCKET_NAME,
		bucketKey:  BUCKET_KEY,
		client:     cloudStorageClient,
		context:    ctx,
	}, nil
}

func (storage *googleCloudStorageClient) Get(key string) ([]byte, error) {
	objHandle := storage.client.Bucket(storage.bucketName).Object(key)

	reader, err := objHandle.NewReader(storage.context)
	if err != nil {
		return nil, err
	}
	var obj []byte
	if obj, err = ioutil.ReadAll(reader); err != nil {
		return nil, err
	}

	return obj, nil
}

func (storage *googleCloudStorageClient) Put(key string, object string) error {
	objHandle := storage.client.Bucket(storage.bucketName).Object(key)
	writer := objHandle.NewWriter(storage.context)
	defer writer.Close()

	_, err := writer.Write([]byte(object))

	return err
}

func (storage *googleCloudStorageClient) Delete(key string) {
	panic("Delete not implemented")
}

func (storage *googleCloudStorageClient) GetObjectSize(key string) (int64, error) {
	panic("GetObjectSize not implemented")
}
