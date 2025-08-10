package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"markdown-viewer/internal/filebrowser"

	"github.com/russross/blackfriday/v2"
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
	tmpl.Execute(w, nil)
}

func (s *Server) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := s.Templates["welcome.html"]
	if !ok {
		log.Println("template not found: welcome.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (s *Server) TreeViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := s.Templates["treeview.html"]
	if !ok {
		log.Println("template not found: treeview.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (s *Server) MarkdownViewHandler(w http.ResponseWriter, r *http.Request) {
	displayPath := strings.TrimPrefix(r.URL.Path, "/view")
	// Use s.Config.TargetDir as the root directory
	fullPath := filepath.Join(s.Config.TargetDir, displayPath)

	source, err := os.ReadFile(fullPath)
	if err != nil {
		s.RenderError(w, http.StatusNotFound)
		return
	}

	output := blackfriday.Run(source, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak))
	data := MarkdownData{
		Title:   filepath.Base(fullPath),
		Content: template.HTML(output),
	}

	tmpl, ok := s.Templates["markdown.html"]
	if !ok {
		log.Println("template not found: markdown.html")
		s.RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func (s *Server) ApiListHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	displayPath := filepath.Clean(pathParam)

	// Use s.Config.TargetDir as the root directory
	items, err := filebrowser.ListDirectory(s.Config.TargetDir, displayPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Directory not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
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
	tmpl.Execute(w, data)
}