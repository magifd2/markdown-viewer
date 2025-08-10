package main

import (
	"markdown-viewer/cmd"
)

// Version is set at build time
var version = "dev"

func main() {
	cmd.Execute(version)
}
