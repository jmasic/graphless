package storage

import (
	"io/ioutil"
	"os"
)

const (
	localFsPath   = "../simple-graphs/graphalytics/s3/"
	localFsPrefix = "graphFileKey"
)

type localFileSystemStorageClient struct {
}

func newLocalFileSystemStorageClient() (Client, error) {
	client := &localFileSystemStorageClient{}
	return client, nil
}

func (storage *localFileSystemStorageClient) Get(key string) ([]byte, error) {
	filePath := localFsPath + localFsPrefix + key
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer closeFile(file)
	return ioutil.ReadAll(file)
}

func (storage *localFileSystemStorageClient) Put(key string, object string) error {
	filePath := localFsPath + localFsPrefix + key
	println("The file path is: ", filePath)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(object)
	defer closeFile(file)
	return err
}

func (storage *localFileSystemStorageClient) Delete(key string) {
	panic("Delete not implemented")
}

func (storage *localFileSystemStorageClient) GetObjectSize(key string) (int64, error) {
	panic("GetObjectSize not implemented")
}

func closeFile(f *os.File) {
	func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
}
