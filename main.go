package main

import (
	"log/slog"
	"primo_test_11/internal/config"
	"primo_test_11/internal/connection"
	"primo_test_11/internal/database"
	"primo_test_11/internal/handler"
	"primo_test_11/internal/repository"
)

func main() {

	// load config of database
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		slog.Error("Got error when load config")
		slog.Error(err.Error())
		return
	}
	// construct db handler with executor
	pgExecutor := database.NewPGExecutor(dbConfig)
	dbHandler := database.NewDBHandler(pgExecutor)

	if dbHandler == nil {
		return
	}

	// construct repo and handler of repo
	productRepo := repository.NewProductRepo(dbHandler)
	productHandler := handler.NewProductHandler(productRepo)

	// construct http connection handler
	serverConfig := config.LoadConfig()
	handler := connection.NewConnectionHandler(serverConfig, productHandler)

	// run forever loop
	handler.RunServe()
}
