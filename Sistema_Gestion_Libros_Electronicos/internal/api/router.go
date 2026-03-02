package api

import "net/http"

// NewMux configura las rutas (8+ servicios) usando el ServeMux de Go 1.22+.
// Patrones con método + path y variables: "GET /api/libros/{id}".
func NewMux(s *Server) *http.ServeMux {
	mux := http.NewServeMux()

	// Salud (útil para demo)
	mux.HandleFunc("GET /health", s.handleHealth)

	// 1) GET /api/libros
	mux.HandleFunc("GET /api/libros", s.handleListBooks)

	// 7) GET /api/libros/buscar?q=texto
	mux.HandleFunc("GET /api/libros/buscar", s.handleSearchBooks)

	// 8) GET /api/estadisticas
	mux.HandleFunc("GET /api/estadisticas", s.handleStats)

	// 3) POST /api/libros
	mux.HandleFunc("POST /api/libros", s.handleCreateBook)

	// 2) GET /api/libros/{id}
	mux.HandleFunc("GET /api/libros/{id}", s.handleGetBookByID)

	// 4) PUT /api/libros/{id}
	mux.HandleFunc("PUT /api/libros/{id}", s.handleUpdateBook)

	// 5) DELETE /api/libros/{id}
	mux.HandleFunc("DELETE /api/libros/{id}", s.handleDeleteBook)

	// 6) PATCH /api/libros/{id}/disponible
	mux.HandleFunc("PATCH /api/libros/{id}/disponible", s.handleToggleAvailability)

	return mux
}