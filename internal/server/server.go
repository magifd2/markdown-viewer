package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"markdown-viewer/internal/config"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	httpServer *http.Server
	Templates  map[string]*template.Template
	Config     config.Config
}

// NewServer creates a new Server instance.
func NewServer(cfg config.Config) (*Server, error) {
	s := &Server{
		Config: cfg,
	}

	// Load templates
	if err := s.LoadTemplates("templates"); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	// Setup HTTP handlers
	mux := http.NewServeMux()

	// Static file server for /static/ directory
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register application handlers
	mux.HandleFunc("/", s.RootHandler)
	mux.HandleFunc("/welcome", s.WelcomeHandler)
	mux.HandleFunc("/files/", s.TreeViewHandler)
	mux.HandleFunc("/view/", s.MarkdownViewHandler)
	mux.HandleFunc("/api/list", s.ApiListHandler)
	mux.HandleFunc("/api/shutdown", s.ShutdownHandler)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.Port),
		Handler: mux,
	}

	return s, nil
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	log.Printf("Server listening on http://127.0.0.1:%d", s.Config.Port)
	// TODO: Implement browser auto-open here

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// LoadTemplates parses all html files from the templates directory
func (s *Server) LoadTemplates(dir string) error {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return err
		}
		cache[name] = tmpl
	}

	s.Templates = cache
	return nil
}