package storage

import (
	"encoding/json"
	"errors"
	"os"

	"sistema_libros/internal/model"
)

// JSONStore implementa model.BookRepository usando un archivo JSON.
type JSONStore struct {
	Path string
}

// Save guarda los libros en el archivo JSON.
func (s JSONStore) Save(books []model.Libro) error {
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, data, 0644)
}

// Load carga los libros desde el archivo JSON.
// Si el archivo no existe o está vacío, retorna slice vacío sin error.
func (s JSONStore) Load() ([]model.Libro, error) {
	_, err := os.Stat(s.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []model.Libro{}, nil
		}
		return nil, err
	}

	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []model.Libro{}, nil
	}

	var books []model.Libro
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, err
	}
	return books, nil
}
