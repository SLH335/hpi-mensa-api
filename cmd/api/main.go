package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hpi-mensa/internal/mealdata/common"
	"hpi-mensa/internal/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Info().Msg("Shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Info().Msg("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Initializing meal data providers
	err := mealdata.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize meal data providers")
	}
	log.Info().Int("count", mealdata.ProviderCount()).Msg("Initialized meal data providers")

	// Starting HTTP server
	server := server.NewServer()
	log.Info().Str("port", server.Addr).Msg("Started HTTP server")

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panic().Err(err).Msg("HTTP server error")
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Info().Msg("Graceful shutdown complete.")
}
