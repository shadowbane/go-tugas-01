package main

import (
	"shadowbane/go-tugas-01/app/exithandler"

	"github.com/shadowbane/go-logger"
	"go.uber.org/zap"

	"shadowbane/go-tugas-01/app/server"
	"shadowbane/go-tugas-01/router"
)

var ApiPort = "0.0.0.0:8080"

func main() {
	logger.Init(logger.LoadEnvForLogger())

	srv := server.
		Get().
		WithAddr(ApiPort).
		WithRouter(router.Api()).
		WithErrLogger(zap.S())

	// start the api server
	go func() {
		zap.S().Info("starting api server at ", ApiPort)

		if err := srv.Start(); err != nil {
			zap.S().Warn(err.Error())
		}
	}()

	exithandler.Init(func() {
		zap.S().Info("Closing Application")

		if err := srv.Close(); err != nil {
			zap.S().Error(err.Error())
		}

		zap.S().Info("Application Closed")
	})

	zap.S().Info("Bye!")
}
