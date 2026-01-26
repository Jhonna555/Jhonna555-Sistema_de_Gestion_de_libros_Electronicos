package catalog

import (
	"errors"
	"strings"

	"sistema_libros/internal/model"
)

var (
	ErrNotFound    = errors.New("libro no encontrado")
	ErrInvalidData = errors.New("datos inválidos")
)

func NextID(books []model.Libro) int {
	maxID := 0
	for _, b := range books {
		if b.ID > maxID {
			maxID = b.ID
		}
	}
	return maxID + 1
}

func AddBook(books []model.Libro, nuevo model.Libro) ([]model.Libro, error) {
	nuevo.Titulo = strings.TrimSpace(nuevo.Titulo)
	nuevo.Autor = strings.TrimSpace(nuevo.Autor)
	nuevo.Categoria = strings.TrimSpace(nuevo.Categoria)
	nuevo.Formato = strings.TrimSpace(nuevo.Formato)
	nuevo.URLArchivo = strings.TrimSpace(nuevo.URLArchivo)

	if nuevo.Titulo == "" || nuevo.Autor == "" {
		return books, ErrInvalidData
	}
	if nuevo.ID <= 0 {
		nuevo.ID = NextID(books)
	}
	books = append(books, nuevo)
	return books, nil
}

func ListBooks(books []model.Libro) []model.Libro {
	// Enfoque funcional: devolvemos una copia (evita modificar el slice original por accidente)
	out := make([]model.Libro, len(books))
	copy(out, books)
	return out
}

func SearchBooks(books []model.Libro, query string) []model.Libro {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}

	var res []model.Libro
	for _, b := range books {
		if strings.Contains(strings.ToLower(b.Titulo), q) ||
			strings.Contains(strings.ToLower(b.Autor), q) ||
			strings.Contains(strings.ToLower(b.Categoria), q) {
			res = append(res, b)
		}
	}
	return res
}

func UpdateBook(books []model.Libro, id int, cambios model.Libro) ([]model.Libro, error) {
	idx := indexByID(books, id)
	if idx == -1 {
		return books, ErrNotFound
	}

	// Solo actualiza campos si vienen con valor (para facilitar edición parcial)
	if strings.TrimSpace(cambios.Titulo) != "" {
		books[idx].Titulo = strings.TrimSpace(cambios.Titulo)
	}
	if strings.TrimSpace(cambios.Autor) != "" {
		books[idx].Autor = strings.TrimSpace(cambios.Autor)
	}
	if strings.TrimSpace(cambios.Categoria) != "" {
		books[idx].Categoria = strings.TrimSpace(cambios.Categoria)
	}
	if cambios.Anio != 0 {
		books[idx].Anio = cambios.Anio
	}
	if strings.TrimSpace(cambios.Formato) != "" {
		books[idx].Formato = strings.TrimSpace(cambios.Formato)
	}
	if strings.TrimSpace(cambios.URLArchivo) != "" {
		books[idx].URLArchivo = strings.TrimSpace(cambios.URLArchivo)
	}

	// Disponible se actualiza solo si el usuario quiere (el menú puede usar ToggleAvailability)
	return books, nil
}

func DeleteBook(books []model.Libro, id int) ([]model.Libro, error) {
	idx := indexByID(books, id)
	if idx == -1 {
		return books, ErrNotFound
	}
	// elimina sin dejar huecos
	books = append(books[:idx], books[idx+1:]...)
	return books, nil
}

func ToggleAvailability(books []model.Libro, id int) ([]model.Libro, error) {
	idx := indexByID(books, id)
	if idx == -1 {
		return books, ErrNotFound
	}
	books[idx].Disponible = !books[idx].Disponible
	return books, nil
}

func GetByID(books []model.Libro, id int) (model.Libro, error) {
	idx := indexByID(books, id)
	if idx == -1 {
		return model.Libro{}, ErrNotFound
	}
	return books[idx], nil
}

func indexByID(books []model.Libro, id int) int {
	for i, b := range books {
		if b.ID == id {
			return i
		}
	}
	return -1
}
