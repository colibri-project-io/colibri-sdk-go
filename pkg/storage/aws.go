package storage

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/cloud"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type awsStorage struct {
	s3Service *s3.S3
	uploader  *s3manager.Uploader
}

func newAwsStorage() *awsStorage {
	var s awsStorage
	s.s3Service = s3.New(cloud.GetAwsSession())
	if _, err := s.s3Service.ListBuckets(nil); err != nil {
		logging.Fatal("An error occurred when trying to connect to the storage provider. Error: %s", err)
	}

	s.uploader = s3manager.NewUploader(cloud.GetAwsSession())
	return &s
}

func (s awsStorage) uploadFile(ctx context.Context, bucket, key string, file *multipart.File) (string, error) {
	result, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   *file,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (s awsStorage) deleteFile(ctx context.Context, bucket, key string) error {
	_, err := s.s3Service.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
