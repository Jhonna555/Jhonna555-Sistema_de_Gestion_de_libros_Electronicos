package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sistema_libros/internal/catalog"
	"sistema_libros/internal/model"
)

// Run ejecuta la UI recibiendo un repositorio abstracto (interfaz).
// Esto desacopla la UI del almacenamiento concreto (JSON, BD, etc.).
func Run(repo model.BookRepository) {
	reader := bufio.NewReader(os.Stdin)

	books, err := repo.Load()
	if err != nil {
		fmt.Println("Error cargando datos:", err)
		books = []model.Libro{}
	}

	for {
		printMenu()

		op := readInt(reader, "Seleccione una opción: ")
		fmt.Println()

		switch op {
		case 1:
			books = uiAddBook(reader, books)
		case 2:
			uiListBooks(books)
		case 3:
			uiSearchBooks(reader, books)
		case 4:
			books = uiEditBook(reader, books)
		case 5:
			books = uiDeleteBook(reader, books)
		case 6:
			books = uiToggleAvailability(reader, books)

		case 7:
			if err := repo.Save(books); err != nil {
				fmt.Println("Error al guardar:", err)
			} else {
				fmt.Println("✅ Datos guardados.")
			}

		case 8:
			loaded, err := repo.Load()
			if err != nil {
				fmt.Println("Error al cargar:", err)
			} else {
				books = loaded
				fmt.Println("✅ Datos cargados.")
			}

		case 0:
			// Guardado automático al salir (suma puntos)
			_ = repo.Save(books)
			fmt.Println("Saliendo... ✅")
			return
		default:
			fmt.Println("Opción inválida.")
		}

		fmt.Println()
		pause(reader)
	}
}

func printMenu() {
	fmt.Println("======================================")
	fmt.Println("  SISTEMA DE GESTIÓN DE E-BOOKS (GO)  ")
	fmt.Println("======================================")
	fmt.Println("1) Agregar libro")
	fmt.Println("2) Listar libros")
	fmt.Println("3) Buscar libro (título/autor/categoría)")
	fmt.Println("4) Editar libro")
	fmt.Println("5) Eliminar libro")
	fmt.Println("6) Cambiar disponibilidad")
	fmt.Println("7) Guardar")
	fmt.Println("8) Cargar")
	fmt.Println("0) Salir")
	fmt.Println("======================================")
}

func uiAddBook(r *bufio.Reader, books []model.Libro) []model.Libro {
	fmt.Println("== Agregar libro ==")

	titulo := readLine(r, "Título: ")
	autor := readLine(r, "Autor: ")
	categoria := readLine(r, "Categoría: ")
	anio := readInt(r, "Año (0 si no aplica): ")
	formato := readLine(r, "Formato (PDF/EPUB/MOBI): ")
	url := readLine(r, "URL o ruta del archivo: ")
	disponible := readBool(r, "¿Disponible? (s/n): ")

	nuevo := model.Libro{
		ID:         0, // se calcula automático en catalog.AddBook
		Titulo:     titulo,
		Autor:      autor,
		Categoria:  categoria,
		Anio:       anio,
		Formato:    formato,
		URLArchivo: url,
		Disponible: disponible,
	}

	updated, err := catalog.AddBook(books, nuevo)
	if err != nil {
		fmt.Println("Error:", err)
		return books
	}

	fmt.Println("✅ Libro agregado.")
	return updated
}

func uiListBooks(books []model.Libro) {
	fmt.Println("== Listado de libros ==")
	list := catalog.ListBooks(books)
	if len(list) == 0 {
		fmt.Println("No hay libros registrados.")
		return
	}
	for _, b := range list {
		fmt.Println(b.String())
	}
}

func uiSearchBooks(r *bufio.Reader, books []model.Libro) {
	fmt.Println("== Buscar libros ==")
	q := readLine(r, "Buscar: ")
	res := catalog.SearchBooks(books, q)
	if len(res) == 0 {
		fmt.Println("No se encontraron resultados.")
		return
	}
	for _, b := range res {
		fmt.Println(b.String())
	}
}

func uiEditBook(r *bufio.Reader, books []model.Libro) []model.Libro {
	fmt.Println("== Editar libro ==")
	id := readInt(r, "ID del libro a editar: ")

	original, err := catalog.GetByID(books, id)
	if err != nil {
		fmt.Println("Error:", err)
		return books
	}

	fmt.Println("Libro actual:")
	fmt.Println(original.String())
	fmt.Println("Deja vacío un campo para mantener el valor actual.")

	titulo := readLine(r, "Nuevo título: ")
	autor := readLine(r, "Nuevo autor: ")
	categoria := readLine(r, "Nueva categoría: ")
	anioStr := readLine(r, "Nuevo año (enter para mantener): ")
	formato := readLine(r, "Nuevo formato: ")
	url := readLine(r, "Nueva URL/ruta: ")

	cambios := model.Libro{
		Titulo:     titulo,
		Autor:      autor,
		Categoria:  categoria,
		Formato:    formato,
		URLArchivo: url,
		Anio:       0,
	}

	if strings.TrimSpace(anioStr) != "" {
		if v, err := strconv.Atoi(strings.TrimSpace(anioStr)); err == nil {
			cambios.Anio = v
		}
	}

	updated, err := catalog.UpdateBook(books, id, cambios)
	if err != nil {
		fmt.Println("Error:", err)
		return books
	}

	fmt.Println("✅ Libro actualizado.")
	return updated
}

func uiDeleteBook(r *bufio.Reader, books []model.Libro) []model.Libro {
	fmt.Println("== Eliminar libro ==")
	id := readInt(r, "ID del libro a eliminar: ")

	updated, err := catalog.DeleteBook(books, id)
	if err != nil {
		fmt.Println("Error:", err)
		return books
	}

	fmt.Println("✅ Libro eliminado.")
	return updated
}

func uiToggleAvailability(r *bufio.Reader, books []model.Libro) []model.Libro {
	fmt.Println("== Cambiar disponibilidad ==")
	id := readInt(r, "ID del libro: ")

	updated, err := catalog.ToggleAvailability(books, id)
	if err != nil {
		fmt.Println("Error:", err)
		return books
	}

	fmt.Println("✅ Disponibilidad actualizada.")
	return updated
}

func readLine(r *bufio.Reader, label string) string {
	fmt.Print(label)
	text, _ := r.ReadString('\n')
	return strings.TrimSpace(text)
}

func readInt(r *bufio.Reader, label string) int {
	for {
		s := readLine(r, label)
		if s == "" {
			return 0
		}
		v, err := strconv.Atoi(s)
		if err == nil {
			return v
		}
		fmt.Println("Ingrese un número válido.")
	}
}

func readBool(r *bufio.Reader, label string) bool {
	for {
		s := strings.ToLower(readLine(r, label))
		if s == "s" || s == "si" || s == "sí" || s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
		fmt.Println("Responda con s/n.")
	}
}

func pause(r *bufio.Reader) {
	fmt.Print("Presione ENTER para continuar...")
	_, _ = r.ReadString('\n')
}
