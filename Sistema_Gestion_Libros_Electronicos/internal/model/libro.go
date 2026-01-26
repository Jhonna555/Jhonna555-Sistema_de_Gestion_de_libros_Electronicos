package model

import "fmt"

type Libro struct {
	ID         int    `json:"id"`
	Titulo     string `json:"titulo"`
	Autor      string `json:"autor"`
	Categoria  string `json:"categoria"`
	Anio       int    `json:"anio"`
	Formato    string `json:"formato"`
	URLArchivo string `json:"urlArchivo"`
	Disponible bool   `json:"disponible"`
}

func (l Libro) String() string {
	estado := "No disponible"
	if l.Disponible {
		estado = "Disponible"
	}
	return fmt.Sprintf(
		"ID: %d | Título: %s | Autor: %s | Categoría: %s | Año: %d | Formato: %s | Archivo: %s | Estado: %s",
		l.ID, l.Titulo, l.Autor, l.Categoria, l.Anio, l.Formato, l.URLArchivo, estado,
	)
}
