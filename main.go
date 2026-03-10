package main

import (
	"fmt"
	"primo_test_11/internal/config"
	"primo_test_11/internal/connection"
	"primo_test_11/internal/database"
	"primo_test_11/internal/handler"
	"primo_test_11/internal/repository"
)

func main() {

	// load config of database
	dbConfig := config.LoadDatabaseConfig()
	// construct db handler with executor
	pgExecutor := database.NewPGExecutor(dbConfig)
	if pgExecutor == nil {
		return
	}
	err := pgExecutor.Connect()
	if err != nil {
		fmt.Println("Error when connect to database", err)
		return
	}

	// construct repo and handler of repo
	productRepo := repository.NewProductRepo(pgExecutor)
	productHandler := handler.NewProductHandler(productRepo)

	// construct http connection handler
	serverConfig := config.LoadConfig()
	if serverConfig == nil {
		fmt.Println("Cannot load config")
		return
	}
	handler := connection.NewConnectionHandler(serverConfig, productHandler)

	// run forever loop
	handler.RunServe()
}
