package connection

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"primo_test_11/internal/config"
	"primo_test_11/internal/database"
	"syscall"
	"time"
)

type ConnectionHandler struct {
	config    *config.Config
	dbHandler *database.DBHandler
}

// construct a connection handler for storing config and db handler
// this struct will make router function handler able to access database
// and manipulate data to database with limitation of db handler
func NewConnectionHandler(config *config.Config, dbHandler *database.DBHandler) *ConnectionHandler {
	slog.Info("Create connection")
	return &ConnectionHandler{
		config:    config,
		dbHandler: dbHandler,
	}
}

// run server to lister as port which registered
// after end will send signal through given channel
func runServer(server *http.Server, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	// run serve
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Stop service with error")
	}
}

// run server and handle to shutdown server gracefully
// such as got signal interupt or terminate
func (c *ConnectionHandler) RunServe() {
	// get address to serve from config
	addr := c.config.GetAddr()
	slog.Info(fmt.Sprintf("Run serve address %s", addr))

	// register channel for trap signal to this process
	// if got sigint or sigterm, will handle to close service gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// get context for cancel it when background is done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// server info
	server := &http.Server{
		Addr: addr,
	}

	// channel for handle server done
	serveDone := make(chan struct{})
	go runServer(server, serveDone)

	// wait all channel
	select {
	// case got signal to process, shutdown server
	case <-sigChan:
		server.Shutdown(ctx)
	// case server process is done even error
	case <-serveDone:
		slog.Info("Server is shutting down")
	}
}
