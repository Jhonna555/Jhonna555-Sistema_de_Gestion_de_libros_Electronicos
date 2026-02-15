package main

import (
	"sistema_libros/internal/storage"
	"sistema_libros/internal/ui"
)

func main() {
	repo := storage.JSONStore{Path: "data/libros.json"}
	ui.Run(repo)
}
