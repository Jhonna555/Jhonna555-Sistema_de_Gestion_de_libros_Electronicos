package model

// BookRepository define un contrato para persistir libros.
// La UI y la lógica del sistema deben depender de esta interfaz,
// no de una implementación concreta (JSON, BD, etc.).
type BookRepository interface {
	Load() ([]Libro, error)
	Save(books []Libro) error
}
