package cloud

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"google.golang.org/api/option"
)

func newFirebaseSession() *firebase.App {
	serviceAccountKeyFilePath, err := filepath.Abs(".serviceAccountKey.json")
	if err != nil {
		logging.Fatal("Unable to load service account key")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logging.Fatal("Firebase initialization error")
	}

	return app
}
