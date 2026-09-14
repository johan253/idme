package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johan253/idme/internal/config"
	"github.com/johan253/idme/internal/db"
	"github.com/johan253/idme/internal/keys"
	"github.com/johan253/idme/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Create a database connection pool
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create database connection pool: %v", err)
		return
	}
	queries := db.New(pool)

	// Create a new Cipher object for decrypting JWKs
	cipher, err := keys.NewAESGCMCipher([]byte(cfg.Kek))
	if err != nil {
		log.Fatalf("Failed to create cipher object: %v", err)
		return
	}

	// Create a new Manager object for managing JWKs in this process
	manager := keys.NewManager(queries, cipher, cfg.JwkRefreshInterval)
	if err := manager.Start(context.Background()); err != nil {
		log.Fatalf("Failed to begin key manager process: %v", err)
		return
	}

	server := server.NewServer(cfg, queries, manager)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
