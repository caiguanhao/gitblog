// +build frontend

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist/*
var frontend embed.FS

func init() {
	f, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		panic(err)
	}
	FrontendFS = http.FS(f)
}
