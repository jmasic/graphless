package clients

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	bucketName      string
	bucketKey       string
	uploadManager   *s3manager.Uploader
	downloadManager *s3manager.Downloader
	*s3.S3
}

func (storage *S3Client) Get(key string) (resultBuf []byte, err error) {
	var byteBuf []byte
	buf := aws.NewWriteAtBuffer(byteBuf)
	_, err = storage.downloadManager.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(storage.bucketName),
		Key:    aws.String(storage.bucketKey + key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return nil, err
	}

	return buf.Bytes(), nil
}

func (storage *S3Client) Put(key string, object string) error {
	//fmt.Println("Putting object ", object)
	_, err := storage.uploadManager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(storage.bucketName),
		Key:    aws.String(storage.bucketKey + key),
		Body:   aws.ReadSeekCloser(strings.NewReader(object)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return err
}

func (storage *S3Client) Delete(object string) {
	fmt.Println("Deleting object ", object)
	//to implement
}

func (storage *S3Client) GetObjectSize(key string) (int64, error) {
	resp, err := storage.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(storage.bucketName)})
	if err != nil {
		fmt.Errorf("Unable to list items in bucket %q, %v", storage.bucketName, err)
		return 0, err
	}

	for _, item := range resp.Contents {
		if *item.Key == storage.bucketKey+key {
			return *item.Size, nil
		}

	}

	return 0, errors.New(fmt.Sprintf("Could not find object with key %s", key))
}
