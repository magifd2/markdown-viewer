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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := TemplateCache["index.html"]
	if !ok {
		log.Println("template not found: index.html")
		RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := TemplateCache["welcome.html"]
	if !ok {
		log.Println("template not found: welcome.html")
		RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func TreeViewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := TemplateCache["treeview.html"]
	if !ok {
		log.Println("template not found: treeview.html")
		RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func MarkdownViewHandler(w http.ResponseWriter, r *http.Request) {
	displayPath := strings.TrimPrefix(r.URL.Path, "/view")
	rootDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
		RenderError(w, http.StatusInternalServerError)
		return
	}
	fullPath := filepath.Join(rootDir, displayPath)

	source, err := os.ReadFile(fullPath)
	if err != nil {
		RenderError(w, http.StatusNotFound)
		return
	}

	output := blackfriday.Run(source, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak))
	data := MarkdownData{
		Title:   filepath.Base(fullPath),
		Content: template.HTML(output),
	}

	tmpl, ok := TemplateCache["markdown.html"]
	if !ok {
		log.Println("template not found: markdown.html")
		RenderError(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func ApiListHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	displayPath := filepath.Clean(pathParam)

	rootDir, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot get working directory"})
		return
	}

	items, err := filebrowser.ListDirectory(rootDir, displayPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Directory not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
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

func RenderError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	data := ErrorData{StatusCode: statusCode, StatusText: http.StatusText(statusCode)}

	tmpl, ok := TemplateCache["error.html"]
	if !ok {
		log.Println("template not found: error.html")
		http.Error(w, "An internal error occurred and the error page template was not found.", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
