package storage

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Client struct {
	bucketName      string
	bucketKey       string
	uploadManager   *s3manager.Uploader
	downloadManager *s3manager.Downloader
	*s3.S3
}

func newS3StorageClient() (Client, error) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-2")}))
	awsS3Client := s3.New(sess)
	return &s3Client{
		bucketName: BUCKET_NAME,
		bucketKey:  BUCKET_KEY,
		uploadManager: s3manager.NewUploaderWithClient(awsS3Client, func(d *s3manager.Uploader) {
			d.Concurrency = 2
		}),
		downloadManager: s3manager.NewDownloaderWithClient(awsS3Client, func(d *s3manager.Downloader) {
			d.Concurrency = 2
		}),
		S3: awsS3Client,
	}, nil
}

func (storage *s3Client) Get(key string) (resultBuf []byte, err error) {
	var byteBuf []byte
	buf := aws.NewWriteAtBuffer(byteBuf)
	log.Info("Getting file from S3: (bucket='", storage.bucketName, "',key='", storage.bucketKey+key, "'")
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

func (storage *s3Client) Put(key string, object string) error {
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

func (storage *s3Client) Delete(object string) {
	fmt.Println("Deleting object ", object)
	//to implement
}

func (storage *s3Client) GetObjectSize(key string) (int64, error) {
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
