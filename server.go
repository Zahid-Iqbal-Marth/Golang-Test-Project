package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	routes "github.com/Zahid-Iqbal-Marth/Golang-Test-Project/routes"
	"github.com/Zahid-Iqbal-Marth/Golang-Test-Project/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Initializing webhook receiver")

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded: batch_size=%d, post_endpoint=%s",
		config.BatchSize, config.PostEndpoint)

	processor := utils.NewBatchProcessor(config.BatchSize, config.PostEndpoint)

	router := routes.SetupRouter(processor)

	server := &http.Server{
		Addr:    ":" + config.ServerPort,
		Handler: router,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("Starting server on port %s", config.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block main and wait for either server errors or shutdown signal
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	case <-shutdown:
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server exited properly")
	}
}
