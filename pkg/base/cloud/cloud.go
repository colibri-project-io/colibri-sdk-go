package cloud

import (
	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type Cloud struct {
	aws      *session.Session
	firebase *firebase.App
}

var instance *Cloud

// Initialize loads the cloud settings according to the configured environment.
func Initialize() {
	instance = &Cloud{}

	switch config.CLOUD {
	case config.CLOUD_AWS:
		instance.aws = newAwsSession()
	case config.CLOUD_FIREBASE:
		instance.firebase = newFirebaseSession()
	case config.CLOUD_AZURE, config.CLOUD_GCP:
		logging.Fatal("Not implemented yet")
	}

	logging.Info("Cloud provider connected")
}

func GetAwsSession() *session.Session {
	return instance.aws
}

func GetFirebaseSession() *firebase.App {
	return instance.firebase
}
