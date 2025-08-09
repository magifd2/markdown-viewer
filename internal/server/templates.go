package server

import (
	"html/template"
	"path/filepath"
)

// TemplateCache holds the parsed templates
var TemplateCache map[string]*template.Template

// LoadTemplates parses all html files from the templates directory
func LoadTemplates(dir string) error {
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

	TemplateCache = cache
	return nil
}
