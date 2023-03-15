package colibri

import (
	"fmt"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/cloud"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/config"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/monitoring"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/observer"
	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/validator"
)

const banner = `
      .   _            _ _ _          _ 
     { \/'o;===       | (_) |        (_)
.----'-/'-/  ___  ___ | |_| |__  _ __ _ 
 '-..-| /   / __ / _ \| | | '_ \| '__| |
    /\/\   | (__| (_) | | | |_) | |  | |
    '--'    \___ \___/|_|_|_.__/|_|  |_|
            project - v%s
`

var Version string = "v.0.0.0"

func InitializeApp() {
	if err := config.Load(); err != nil {
		logging.Fatal("Occurred a error on try load configs: %v", err)
	}

	printBanner()
	printApplicationName()

	validator.Initialize()
	observer.Initialize()
	monitoring.Initialize()
	cloud.Initialize()
}

func printBanner() {
	fmt.Printf(banner, Version)
}

func printApplicationName() {
	fmt.Printf("\n# %s #\n\n", config.APP_NAME)
}
