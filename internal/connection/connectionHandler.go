package connection

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"primo_test_11/internal/config"
	"primo_test_11/internal/handler"
	"syscall"
	"time"

	_ "primo_test_11/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ConnectionHandler struct {
	config  *config.Config
	product *handler.ProductHandler
}

// construct a connection handler for storing config and db handler
// this struct will make router function handler able to access database
// and manipulate data to database with limitation of db handler
func NewConnectionHandler(config *config.Config, product *handler.ProductHandler) *ConnectionHandler {
	fmt.Println("Create connection")

	routePaths := product.GetRouteInfo()
	for _, routePath := range routePaths {
		http.Handle(routePath.Path, handler.MakeHandler(routePath))
		fmt.Println(fmt.Sprintf("Found path for register %s", routePath.Path))
	}

	return &ConnectionHandler{
		config:  config,
		product: product,
	}
}

// run server to lister as port which registered
// after end will send signal through given channel
func runServer(server *http.Server, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	// run serve
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Stop service with error")
	}
}

// run server and handle to shutdown server gracefully
// such as got signal interupt or terminate
func (c *ConnectionHandler) RunServe() {
	http.Handle("/api-docs/", httpSwagger.WrapHandler)

	// get address to serve from config
	addr := c.config.GetAddr()
	fmt.Println(fmt.Sprintf("Run serve address %s", addr))

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
		fmt.Println("Server is shutting down")
	}
}
