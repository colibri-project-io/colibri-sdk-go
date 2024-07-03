package storage

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
)

const (
	storage_transaction = "Storage"
	connection_error    = "An error occurred when trying to connect to the storage provider. Error: %s"
)

type storage interface {
	downloadFile(ctx context.Context, bucket, key string) (*os.File, error)
	uploadFile(ctx context.Context, bucket, key string, file *multipart.File) (string, error)
	deleteFile(ctx context.Context, bucket, key string) error
}

var instance storage

// Initialize initializes the storage provider based on the configured cloud.
//
// No parameters.
// No return values.
func Initialize() {
	switch config.CLOUD {
	case config.CLOUD_AWS:
		instance = newAwsStorage()
	case config.CLOUD_GCP, config.CLOUD_FIREBASE:
		instance = newGcpStorage()
	}

	logging.Info("Storage provider connected")
}

// DownloadFile downloads a file from the storage provider.
//
// ctx: the context for the operation.
// bucket: the storage bucket from which the file is downloaded.
// key: the key or identifier of the file to be downloaded.
// Returns a file pointer and an error.
func DownloadFile(ctx context.Context, bucket, key string) (*os.File, error) {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(ctx, storage_transaction, map[string]string{
			"method": "Download",
			"bucket": bucket,
			"key":    key,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	return instance.downloadFile(ctx, bucket, key)
}

// UploadFile uploads a file to the storage provider.
//
// ctx: the context for the operation.
// bucket: the storage bucket to upload the file to.
// key: the key or identifier of the file to be uploaded.
// file: the file to be uploaded.
// Returns the location of the uploaded file and an error, if any.
func UploadFile(ctx context.Context, bucket, key string, file *multipart.File) (string, error) {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(ctx, storage_transaction, map[string]string{
			"method": "Upload",
			"bucket": bucket,
			"key":    key,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	return instance.uploadFile(ctx, bucket, key, file)
}

// DeleteFile deletes a file from the storage provider.
//
// ctx: the context for the operation.
// bucket: the storage bucket from which the file is deleted.
// key: the key or identifier of the file to be deleted.
// Returns an error.
func DeleteFile(ctx context.Context, bucket, key string) error {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(ctx, storage_transaction, map[string]string{
			"method": "Delete",
			"bucket": bucket,
			"key":    key,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	return instance.deleteFile(ctx, bucket, key)
}
