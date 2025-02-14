package main

import (
	"context"
	"fmt"
	"hornet/api/followers"
	"hornet/api/followers/repository"
	"hornet/api/followers/service"
	config "hornet/config/followers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create a parent context for the application with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gracefully handle shutdown signals (e.g., Ctrl+C)
	go handleShutdown(cancel)

	// Set up Neo4j driver and defer disconnect
	driver, err := SetupNeo4jDriver(ctx, cfg.Neo4jURI, "neo4j", "neo4j")
	if err != nil {
		log.Fatalf("Failed to set up Neo4j client: %v", err)
	}
	defer driver.Close(ctx)

	// Initialize repository and service layers
	followersRepository := repository.NewFollowersRepository(driver)
	followersService := service.NewFollowersService(followersRepository)

	// Set up router with service
	r := followers.Router(followersService)

	// Start the Gin server
	server := startServer(r, cfg.ServerPort)

	// Wait for shutdown signal (context cancellation)
	<-ctx.Done()

	// Gracefully shut down the server
	gracefulShutdown(server)
}

// handleShutdown listens for interrupt signals to initiate a graceful shutdown.
func handleShutdown(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Received shutdown signal...")
	cancel() // Cancel the context to initiate shutdown
}

// SetupNeo4jDriver creates and returns a Neo4j driver instance.
func SetupNeo4jDriver(ctx context.Context, uri, username, password string) (neo4j.DriverWithContext, error) {
	// Create a new driver with context.
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	// Verify connectivity using the provided context.
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, err
	}
	return driver, nil
}

// startServer starts the HTTP server in a goroutine.
func startServer(r http.Handler, port string) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Run the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started successfully. Listening on port", port)
	return server
}

// gracefulShutdown shuts down the server gracefully, allowing ongoing requests to complete.
func gracefulShutdown(server *http.Server) {
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout
	defer cancelShutdown()

	// Shut down the server gracefully
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	} else {
		log.Println("Server shutdown successfully.")
	}
}
