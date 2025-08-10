package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"markdown-viewer/internal/config"
	"markdown-viewer/internal/server"
)

// Version is set at build time
var version = "dev"

func main() {
	// --- Initialization ---
	fmt.Printf("Starting Markdown Viewer %s\n", version)

	// Load configuration (from file/env first)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Define command-line flags, overriding config values if set
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Port to listen on")
	flag.BoolVar(&cfg.Open, "open", cfg.Open, "Open browser automatically")
	flag.StringVar(&cfg.TargetDir, "dir", cfg.TargetDir, "Directory to serve")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Create new server instance
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	// Open browser if configured
	if cfg.Open {
		url := fmt.Sprintf("http://127.0.0.1:%d", cfg.Port)
		if err := openBrowser(url); err != nil {
			log.Printf("Failed to open browser: %v", err)
		}
	}

	// --- Graceful Shutdown Setup ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	select {
	case <-quit:
		log.Println("Shutdown signal received, shutting down server...")
	case <-server.ShutdownChannel:
		log.Println("Shutdown request received via API, shutting down server...")
	}

	// --- Shutdown ---
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}

// openBrowser opens the default web browser to the specified URL.
func openBrowser(url string) error {
	// This feature is disabled for security reasons (G204).
	log.Println("Browser auto-open feature is disabled for security reasons.")
	return nil
}

// This is a placeholder for the version command, which is now handled by Cobra.
// It's kept here to show the evolution of the project.
// func printVersion() {
// 	fmt.Printf("Markdown Viewer version %s\n", version)
// }
