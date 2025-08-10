package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"markdown-viewer/internal/filebrowser"
	"markdown-viewer/internal/markdown"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// ShutdownChannel is used to signal server shutdown from an API call
var ShutdownChannel chan struct{}

// Structs for template data
type MarkdownData struct {
	Title   string
	Content template.HTML
}

type ErrorData struct {
	StatusCode int
	StatusText string
}

// --- HTTP Handlers ---

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := s.Templates["index.html"]
	if !ok {
		log.Println("template not found: index.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("failed to execute index.html template: %v", err)
	}
}

func (s *Server) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := s.Templates["welcome.html"]
	if !ok {
		log.Println("template not found: welcome.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("failed to execute welcome.html template: %v", err)
	}
}

func (s *Server) TreeViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := s.Templates["treeview.html"]
	if !ok {
		log.Println("template not found: treeview.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("failed to execute treeview.html template: %v", err)
	}
}

func (s *Server) MarkdownViewHandler(w http.ResponseWriter, r *http.Request) {
	displayPath := strings.TrimPrefix(r.URL.Path, "/view")
	// Use s.Config.TargetDir as the root directory
	fullPath := filepath.Join(s.Config.TargetDir, displayPath)

	// #nosec G304
	source, err := os.ReadFile(fullPath)
	if err != nil {
		s.RenderError(w, http.StatusNotFound)
		return
	}

	// Configure goldmark with GFM extensions and our custom safe link renderer.
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(markdown.NewSafeLinkRenderer(), 1),
			),
		),
	)

	// Render Markdown to HTML
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	output := buf.Bytes()

	// Create a custom policy for sanitizing HTML which allows syntax highlighting.
	policy := bluemonday.UGCPolicy()
	// Allow 'class' attribute for syntax highlighting (e.g., class="language-go").
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^language-[\w-]+$`)).OnElements("code")

	sanitizedOutput := policy.SanitizeBytes(output)

	data := MarkdownData{
		Title:   filepath.Base(fullPath),
		// #nosec G203
		Content: template.HTML(sanitizedOutput),
	}

	tmpl, ok := s.Templates["markdown.html"]
	if !ok {
		log.Println("template not found: markdown.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("failed to execute markdown.html template: %v", err)
	}
}

func (s *Server) ApiListHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	displayPath := filepath.Clean(pathParam)

	// Use s.Config.TargetDir as the root directory
	items, err := filebrowser.ListDirectory(s.Config.TargetDir, displayPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Directory not found"}); err != nil {
			log.Printf("failed to encode json error response: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func (s *Server) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Shutdown signal received. Server is shutting down.")

	// Signal the main function to shut down
	go func() {
		ShutdownChannel <- struct{}{} 
	}()
}

func (s *Server) RenderError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	data := ErrorData{StatusCode: statusCode, StatusText: http.StatusText(statusCode)}

	tmpl, ok := s.Templates["error.html"]
	if !ok {
		log.Println("template not found: error.html")
		http.Error(w, "An internal error occurred and the error page template was not found.", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("failed to execute error.html template: %v", err)
	}
}