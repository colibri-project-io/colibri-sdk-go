package storage

import (
	"context"
	"io"
	"mime/multipart"

	gcp_storage "cloud.google.com/go/storage"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type gcpStorage struct {
	client *gcp_storage.Client
}

func newGcpStorage() *gcpStorage {
	client, err := gcp_storage.NewClient(context.Background())
	if err != nil {
		logging.Fatal(connection_error, err)
	}

	return &gcpStorage{client}

}

func (s *gcpStorage) uploadFile(ctx context.Context, bucket, key string, file *multipart.File) (string, error) {
	wc := s.client.Bucket(bucket).Object(key).NewWriter(ctx)
	if _, err := io.Copy(wc, *file); err != nil {
		return "", err
	}
	fileLink := wc.MediaLink
	if err := wc.Close(); err != nil {
		return "", err
	}

	return fileLink, nil
}

func (s *gcpStorage) deleteFile(ctx context.Context, bucket, key string) error {
	return s.client.Bucket(bucket).Object(key).Delete(ctx)
}
