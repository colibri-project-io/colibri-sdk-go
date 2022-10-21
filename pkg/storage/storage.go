package storage

import (
	"context"
	"mime/multipart"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
)

const (
	storage_transaction = "Storage"
	connection_error    = "An error occurred when trying to connect to the storage provider. Error: %s"
)

type storage interface {
	uploadFile(ctx context.Context, bucket, key string, file *multipart.File) (string, error)
	deleteFile(ctx context.Context, bucket, key string) error
}

var instance storage

func Initialize() {
	switch config.CLOUD {
	case config.CLOUD_AWS:
		instance = newAwsStorage()
	case config.CLOUD_GCP, config.CLOUD_FIREBASE:
		instance = newGcpStorage()
	}

	logging.Info("Storage provider connected")
}

func UploadFile(ctx context.Context, bucket, key string, file *multipart.File) (string, error) {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, storage_transaction, map[string]interface{}{
			"method": "Upload",
			"bucket": bucket,
			"key":    key,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	return instance.uploadFile(ctx, bucket, key, file)
}

func DeleteFile(ctx context.Context, bucket, key string) error {
	txn := monitoring.GetTransactionInContext(ctx)
	if txn != nil {
		segment := monitoring.StartTransactionSegment(txn, storage_transaction, map[string]interface{}{
			"method": "Delete",
			"bucket": bucket,
			"key":    key,
		})
		defer monitoring.EndTransactionSegment(segment)
	}

	return instance.deleteFile(ctx, bucket, key)
}
