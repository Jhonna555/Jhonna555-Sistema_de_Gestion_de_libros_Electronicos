package storage

import (
	"encoding/json"
	"errors"
	"os"

	"sistema_libros/internal/model"
)

func SaveBooks(path string, books []model.Libro) error {
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadBooks(path string) ([]model.Libro, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []model.Libro{}, nil
		}
		return nil, err
	}

	data, err := os.ReadFile(path)
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
