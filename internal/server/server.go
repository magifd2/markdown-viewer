package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"markdown-viewer/internal/config"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	httpServer *http.Server
	Templates  map[string]*template.Template
	Config     config.Config
	staticFS   http.Handler
}

// NewServer creates a new Server instance.
func NewServer(cfg config.Config) (*Server, error) {
	s := &Server{
		Config:   cfg,
		staticFS: http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	}

	// Load templates
	if err := s.LoadTemplates("templates"); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", cfg.Port),
		Handler:           s, // The server itself is the handler
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return s, nil
}

// ServeHTTP is the single entry point for all HTTP requests.
// It validates the path for security and then routes to the appropriate handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// --- 1. Path Validation (Security Check) ---
	// Use the raw RequestURI to catch traversal attempts before any cleaning.
	parsedURL, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		log.Printf("[PathValidation] Bad Request: Malformed URI: %q", r.RequestURI)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if the raw path contains any ".." segments. This is the most reliable way
	// to catch directory traversal attempts, including encoded ones, because the
	// path has been unescaped by url.ParseRequestURI.
	if containsDotDot(parsedURL.Path) {
		log.Printf("[PathValidation] Forbidden: Directory traversal attempt detected in URI: %q", r.RequestURI)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// --- 2. Routing ---
	// Now that we've blocked malicious paths, we can safely use the standard, cleaned
	// r.URL.Path for routing.
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/static/"):
		s.staticFS.ServeHTTP(w, r)
	case strings.HasPrefix(path, "/view/"):
		s.MarkdownViewHandler(w, r)
	case path == "/api/list":
		s.ApiListHandler(w, r)
	case path == "/api/shutdown":
		s.ShutdownHandler(w, r)
	case path == "/files/" || path == "/files":
		s.TreeViewHandler(w, r)
	case path == "/welcome":
		s.WelcomeHandler(w, r)
	case path == "/":
		s.RootHandler(w, r)
	default:
		s.RenderError(w, http.StatusNotFound)
	}
}

// containsDotDot checks if a path contains ".." segments. It's a simple but effective
// way to prevent directory traversal attacks.
func containsDotDot(path string) bool {
	for _, segment := range strings.Split(path, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
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
