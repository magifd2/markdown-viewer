package assets

import (
	"embed"
	"io/fs"
)

//go:embed embed_assets/**/*
var EmbeddedFiles embed.FS

// GetStaticFS returns the embedded static filesystem.
func GetStaticFS() (fs.FS, error) {
	return fs.Sub(EmbeddedFiles, "embed_assets/static")
}

// GetTemplatesFS returns the embedded templates filesystem.
func GetTemplatesFS() (fs.FS, error) {
	return fs.Sub(EmbeddedFiles, "embed_assets/templates")
}
