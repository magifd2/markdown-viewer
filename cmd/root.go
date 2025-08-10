package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"markdown-viewer/internal/browser"
	"markdown-viewer/internal/config"
	"markdown-viewer/internal/server"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "mdv",
	Short: "mdv is a simple Markdown viewer with a built-in directory tree navigator.",
	Long: `mdv is a lightweight, single-binary Markdown viewer that turns any directory of Markdown files into a browsable, elegant documentation site.

It provides a 2-pane UI to navigate and render a directory of Markdown files.`, 
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Bind flags to viper
		if err := viper.BindPFlag("port", cmd.PersistentFlags().Lookup("port")); err != nil {
			return fmt.Errorf("failed to bind 'port' flag: %w", err)
		}
		if err := viper.BindPFlag("open", cmd.PersistentFlags().Lookup("open")); err != nil {
			return fmt.Errorf("failed to bind 'open' flag: %w", err)
		}
		if err := viper.BindPFlag("target_dir", cmd.PersistentFlags().Lookup("dir")); err != nil {
			return fmt.Errorf("failed to bind 'dir' flag: %w", err)
		}

		// Unmarshal final config
		if err := viper.Unmarshal(&cfg); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// Change working directory if target_dir is specified
		if cfg.TargetDir != "." && cfg.TargetDir != "" {
			absPath, err := filepath.Abs(cfg.TargetDir)
			if err != nil {
				return fmt.Errorf("invalid target directory: %w", err)
			}
			if err := os.Chdir(absPath); err != nil {
				return fmt.Errorf("failed to change directory to %s: %w", absPath, err)
			}
			log.Printf("Changed working directory to: %s", absPath)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize the shutdown channel
		server.ShutdownChannel = make(chan struct{})

		// Create a new server instance
		srv, err := server.NewServer(cfg)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		// Start the server in a goroutine
		go func() {
			log.Printf("Server listening on http://127.0.0.1:%d", cfg.Port)
			if cfg.Open {
				url := fmt.Sprintf("http://127.0.0.1:%d", cfg.Port)
				if err := browser.Open(url); err != nil {
					log.Printf("Failed to open browser: %v", err)
				}
			}
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

		// ---
		// Shutdown
		// ---
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		log.Println("Server exited gracefully")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	rootCmd.PersistentFlags().IntP("port", "p", 8080, "Port to listen on")
	rootCmd.PersistentFlags().BoolP("open", "o", false, "Open browser automatically")
	rootCmd.PersistentFlags().StringP("dir", "d", ".", "Directory to serve Markdown files from")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
