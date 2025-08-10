package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
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
func openBrowser(rawURL string) error {
	// First, validate the URL to ensure it's a well-formed http/https URL.
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL scheme must be http or https")
	}

	var cmdName string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmdName = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmdName = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmdName = "xdg-open"
	}
	args = append(args, u.String()) // Use the parsed and validated URL

	// Find the absolute path of the command to be executed.
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("command not found: %s", cmdName)
	}

	// #nosec G204
	// The command and URL have been validated, so we can safely proceed.
	cmd := exec.Command(cmdPath, args...)
	return cmd.Start()
}