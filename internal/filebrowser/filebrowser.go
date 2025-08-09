package filebrowser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListItem represents a file or directory for the API response.
type ListItem struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

// ListDirectory scans a directory and returns a list of items.
func ListDirectory(rootPath, displayPath string) ([]ListItem, error) {
	fullPath := filepath.Join(rootPath, displayPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var items []ListItem
	for _, entry := range entries {
		if entry.IsDir() {
			items = append(items, ListItem{
				Name:  entry.Name(),
				Path:  filepath.ToSlash(filepath.Join(displayPath, entry.Name())),
				IsDir: true,
			})
		} else {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".md" || ext == ".markdown" {
				items = append(items, ListItem{
					Name:  entry.Name(),
					Path:  filepath.ToSlash(filepath.Join(displayPath, entry.Name())),
					IsDir: false,
				})
			}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir // Dirs first
		}
		return items[i].Name < items[j].Name
	})

	return items, nil
}
