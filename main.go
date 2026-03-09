package main

import (
	"log/slog"
	"primo_test_11/internal/config"
	"primo_test_11/internal/connection"
	"primo_test_11/internal/database"
)

func main() {

	// load config
	serverConfig := config.LoadConfig()
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		slog.Error("Got error when load config")
		slog.Error(err.Error())
		return
	}
	pgExecutor := database.NewPGExecutor(dbConfig)
	dbHandler := database.NewDBHandler(pgExecutor)

	if dbHandler == nil {
		return
	}

	// construct http connection handler
	handler := connection.NewConnectionHandler(serverConfig, dbHandler)

	// run forever loop
	handler.RunServe()

}
