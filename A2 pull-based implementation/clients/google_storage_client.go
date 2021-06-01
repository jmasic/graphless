package clients

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

type GoogleCloudStorageClient struct {
	bucketName string
	bucketKey  string
	client     *storage.Client
	context    context.Context
}

func (storage *GoogleCloudStorageClient) Get(key string) ([]byte, error) {
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

func (storage *GoogleCloudStorageClient) Put(key string, object string) error {
	objHandle := storage.client.Bucket(storage.bucketName).Object(key)
	writer := objHandle.NewWriter(storage.context)
	defer writer.Close()

	_, err := writer.Write([]byte(object))

	return err
}

func (storage *GoogleCloudStorageClient) Delete(key string) {
	panic("Delete not implemented")
}

func (storage *GoogleCloudStorageClient) GetObjectSize(key string) (int64, error) {
	panic("GetObjectSize not implemented")
}
