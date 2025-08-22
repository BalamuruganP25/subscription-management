package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"subscription-management/pkg/handler"
	"syscall"
	"time"

	"subscription-management/pkg/repository"
)

func main() {

	var processConfig handler.ProcessConfig
	db, err := repository.SetUpDB()
	if err != nil {
		log.Fatal("failed to setup setup Database %w", err)
	}

	CurdRepo := repository.NewCurdRepo(db)
	processConfig.CurdRepo = CurdRepo
	processConfig.StripeKey = os.Getenv("STRIPE_SECRET_KEY")
	MigrateDBRepo(db)
	
	// Initialize the server (but don't start it yet)
	srv := initWebServer(processConfig)

	// Channel to listen for errors from server.ListenAndServe
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		log.Println("🚀 Server is starting on :8089...")
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for OS signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait for signal or server error
	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, shutting down gracefully...", sig)

		// Create a deadline to wait for current operations to finish
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to gracefully shutdown server: %v", err)
		}
	}

	log.Println("Server stopped")

}
