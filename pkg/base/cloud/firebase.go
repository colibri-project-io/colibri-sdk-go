package cloud

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"google.golang.org/api/option"
)

func newFirebaseSession() *firebase.App {
	opt := option.WithCredentialsJSON([]byte(config.CLOUD_SECRET))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logging.Fatal("Firebase initialization error")
	}

	return app
}
