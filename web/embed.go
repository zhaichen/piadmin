package web

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist/*
var assets embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(assets, "frontend/dist")
}
