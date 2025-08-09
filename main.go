package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"markdown-viewer/internal/server"
)

// Version is set at build time
var version = "dev"

func main() {
	// --- Initialization ---
	fmt.Printf("Starting Markdown Viewer %s\n", version)

	if err := server.LoadTemplates("templates"); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// --- Server and Handler Setup ---
	server.ShutdownChannel = make(chan struct{})

	mux := http.NewServeMux()

	// Static file server for /static/ directory
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register application handlers
	mux.HandleFunc("/", server.RootHandler)
	mux.HandleFunc("/welcome", server.WelcomeHandler)
	mux.HandleFunc("/files/", server.TreeViewHandler)
	mux.HandleFunc("/view/", server.MarkdownViewHandler)
	mux.HandleFunc("/api/list", server.ApiListHandler)
	mux.HandleFunc("/api/shutdown", server.ShutdownHandler)

	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	// --- Graceful Shutdown Setup ---
	go func() {
		fmt.Println("Server listening on http://127.0.0.1:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

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