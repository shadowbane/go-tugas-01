package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/shadowbane/go-logger"
	"github.com/shadowbane/go-tugas-01/app/exithandler"
	"github.com/shadowbane/go-tugas-01/app/server"
	"github.com/shadowbane/go-tugas-01/router"
	"go.uber.org/zap"
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

		// Note to self: this should fix the issue on dockerized app,
		// where the server cannot start (when not using traefik)
		// because of port conflict. ToDo: update your other app on homeserver!
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Errorf("Server Error: %s", err.Error())

			// Close after detecting server error
			// this will trigger the container to be recreated.
			// In case of using supervisor, this could trigger restart
			os.Exit(1)
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
