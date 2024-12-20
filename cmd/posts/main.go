package main

import (
	"context"
	"fmt"
	"hornet/api/posts"
	"hornet/api/posts/repository"
	"hornet/api/posts/service"
	config "hornet/config/posts"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create a parent context for the application with cancellation support
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gracefully handle shutdown signals (e.g., Ctrl+C)
	go handleShutdown(cancel)

	// Set up MongoDB client and defer disconnect
	client, db := setupMongoClient(ctx, cfg.MongoURI, cfg.DBName)
	defer client.Disconnect(ctx)

	// Initialize repository and service layers
	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)

	// Set up router with service
	r := posts.Router(postService)

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

// setupMongoClient initializes and returns a MongoDB client and the database.
func setupMongoClient(ctx context.Context, uri, dbName string) (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, client.Database(dbName)
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
