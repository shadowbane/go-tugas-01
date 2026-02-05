package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/shadowbane/go-tugas-01/handlers"
	"github.com/shadowbane/go-tugas-01/pkg/config"
	"github.com/shadowbane/go-tugas-01/pkg/database"
	"github.com/shadowbane/go-tugas-01/pkg/exithandler"
	"github.com/shadowbane/go-tugas-01/pkg/server"
	"github.com/shadowbane/go-tugas-01/repositories"
	"github.com/shadowbane/go-tugas-01/router"
	"github.com/shadowbane/go-tugas-01/services"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Get()

	// Setup database
	db, err := database.InitDB(cfg.GetPSQLConnectionString())
	if err != nil {
		zap.S().Errorf("Failed to initialize database: %s", err.Error())
		os.Exit(1)
	}

	// Dependency Injection
	// Note to mentor: IDK if you read this,
	//  but I hate this pattern. lol. v(^_^)
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	globalHandler := &handlers.Handler{
		CategoryHandler:    handlers.NewCategoryHandler(categoryService),
		ProductHandler:     handlers.NewProductHandler(productService),
		TransactionHandler: handlers.NewTransactionHandler(transactionService),
		ReportHandler:      handlers.NewReportHandler(transactionService),
	}

	srv := server.
		Get().
		WithAddr(cfg.GetAddr()).
		WithRouter(router.Api(globalHandler)).
		WithErrLogger(zap.S())

	// start the api server
	go func() {
		zap.S().Info("starting api server at ", cfg.GetAddr())

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

		if err := db.Close(); err != nil {
			zap.S().Error("Error closing database: ", err)
		}

		zap.S().Info("Application Closed")
	})

	zap.S().Info("Bye!")
}
