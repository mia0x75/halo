// Package web serves the compiled frontend from the server binary.
package web

import (
	"embed"
	"io/fs"
)

// The build copies frontend/dist into this directory before compiling; the
// repository only tracks .gitkeep so `go build` still works on a fresh clone.
//
//go:embed all:dist
var assets embed.FS

// Dist returns the embedded frontend tree.
func Dist() (fs.FS, error) {
	return fs.Sub(assets, "dist")
}
