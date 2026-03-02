package catalog

import (
	"testing"

	"sistema_libros/internal/model"
)

func sampleBooks() []model.Libro {
	return []model.Libro{
		{ID: 1, Titulo: "El Principito", Autor: "Antoine de Saint-Exupéry", Categoria: "Literatura", Anio: 1943, Formato: "PDF", URLArchivo: "el_principito.pdf", Disponible: true},
		{ID: 2, Titulo: "Cien años de soledad", Autor: "Gabriel García Márquez", Categoria: "Novela", Anio: 1967, Formato: "EPUB", URLArchivo: "cien_anos_de_soledad.epub", Disponible: true},
		{ID: 3, Titulo: "1984", Autor: "George Orwell", Categoria: "Ciencia ficción", Anio: 1949, Formato: "EPUB", URLArchivo: "1984.epub", Disponible: true},
	}
}

// ---------- GetByID ----------

func TestGetByID_Found(t *testing.T) {
	books := sampleBooks()

	got, err := GetByID(books, 2)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got.ID != 2 {
		t.Fatalf("expected ID=2, got=%d", got.ID)
	}
	if got.Titulo != "Cien años de soledad" {
		t.Fatalf("expected titulo 'Cien años de soledad', got=%q", got.Titulo)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	books := sampleBooks()

	_, err := GetByID(books, 999)
	if err == nil {
		t.Fatalf("expected error for missing id, got nil")
	}
}

// ---------- SearchBooks ----------

func TestSearchBooks_ByTitle(t *testing.T) {
	books := sampleBooks()

	res := SearchBooks(books, "principito")
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got=%d", len(res))
	}
	if res[0].ID != 1 {
		t.Fatalf("expected ID=1, got=%d", res[0].ID)
	}
}

func TestSearchBooks_ByAuthor(t *testing.T) {
	books := sampleBooks()

	res := SearchBooks(books, "orwell")
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got=%d", len(res))
	}
	if res[0].ID != 3 {
		t.Fatalf("expected ID=3, got=%d", res[0].ID)
	}
}

func TestSearchBooks_ByCategory(t *testing.T) {
	books := sampleBooks()

	res := SearchBooks(books, "novela")
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got=%d", len(res))
	}
	if res[0].ID != 2 {
		t.Fatalf("expected ID=2, got=%d", res[0].ID)
	}
}

// ---------- AddBook ----------

func TestAddBook_AssignsNewID(t *testing.T) {
	books := sampleBooks()

	newBook := model.Libro{
		// ID debería asignarlo el sistema
		Titulo:     "La Odisea",
		Autor:      "Homero",
		Categoria:  "Épico",
		Anio:       -700,
		Formato:    "PDF",
		URLArchivo: "la_odisea.pdf",
		Disponible: false,
	}

	updated, err := AddBook(books, newBook)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if len(updated) != len(books)+1 {
		t.Fatalf("expected len=%d, got=%d", len(books)+1, len(updated))
	}

	last := updated[len(updated)-1]
	if last.ID == 0 {
		t.Fatalf("expected new book to have ID assigned, got ID=0")
	}
	if last.ID != 4 {
		// Asumiendo IDs correlativos (max+1). Si tu lógica es distinta, ajustamos este assert.
		t.Fatalf("expected new ID=4, got=%d", last.ID)
	}
	if last.Titulo != "La Odisea" {
		t.Fatalf("expected titulo 'La Odisea', got=%q", last.Titulo)
	}
}

// ---------- UpdateBook ----------

func TestUpdateBook_UpdatesFields(t *testing.T) {
	books := sampleBooks()

	changes := model.Libro{
		Titulo:     "1984 (Editado)",
		Autor:      "George Orwell",
		Categoria:  "Ciencia ficción",
		Anio:       1949,
		Formato:    "EPUB",
		URLArchivo: "1984_v2.epub",
		Disponible: false,
	}

	updated, err := UpdateBook(books, 3, changes)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	got, err := GetByID(updated, 3)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if got.Titulo != "1984 (Editado)" {
		t.Fatalf("expected updated titulo, got=%q", got.Titulo)
	}
	if got.URLArchivo != "1984_v2.epub" {
		t.Fatalf("expected updated urlArchivo, got=%q", got.URLArchivo)
	}
	if got.Disponible != true {
		t.Fatalf("expected disponible=true (sin cambios), got=%v", got.Disponible)
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	books := sampleBooks()

	_, err := UpdateBook(books, 999, model.Libro{Titulo: "X"})
	if err == nil {
		t.Fatalf("expected error when updating missing id, got nil")
	}
}

// ---------- DeleteBook ----------

func TestDeleteBook_RemovesItem(t *testing.T) {
	books := sampleBooks()

	updated, err := DeleteBook(books, 2)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("expected len=2, got=%d", len(updated))
	}

	_, err = GetByID(updated, 2)
	if err == nil {
		t.Fatalf("expected error after delete, book should not exist")
	}
}

func TestDeleteBook_NotFound(t *testing.T) {
	books := sampleBooks()

	_, err := DeleteBook(books, 999)
	if err == nil {
		t.Fatalf("expected error when deleting missing id, got nil")
	}
}

// ---------- ToggleAvailability ----------

func TestToggleAvailability_Toggles(t *testing.T) {
	books := sampleBooks()

	// ID=1 inicia true
	b, _ := GetByID(books, 1)
	if b.Disponible != true {
		t.Fatalf("expected initial disponible=true, got=%v", b.Disponible)
	}

	updated, err := ToggleAvailability(books, 1)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	b2, err := GetByID(updated, 1)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if b2.Disponible != false {
		t.Fatalf("expected disponible=false after toggle, got=%v", b2.Disponible)
	}
}

func TestToggleAvailability_NotFound(t *testing.T) {
	books := sampleBooks()

	_, err := ToggleAvailability(books, 999)
	if err == nil {
		t.Fatalf("expected error when toggling missing id, got nil")
	}
}
