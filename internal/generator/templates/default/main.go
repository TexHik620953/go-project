package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"@APPNAME@/internal/appconfig"
	"@APPNAME@/internal/application"
)

func main() {
	appContext, cancelFunc := context.WithCancel(context.Background())

	appConfig, err := appconfig.LoadAppConfig()
	if err != nil {
		log.Fatalf("failed to load app config: %s", err.Error())
	}

	app, err := application.New(appContext, appConfig)
	if err != nil {
		log.Fatalf("failed to create application: %s", err.Error())
	}

	err = app.Start()
	if err != nil {
		log.Fatalf("failed to start application: %s", err.Error())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancelFunc()
}
