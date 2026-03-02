package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"sistema_libros/internal/catalog"
	"sistema_libros/internal/model"
)

type Server struct {
	Repo model.BookRepository
	mu   sync.RWMutex // protege Load/Save contra escrituras concurrentes
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleListBooks(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (s *Server) handleGetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido", err)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	b, err := catalog.GetByID(books, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "libro no encontrado", err)
		return
	}

	writeJSON(w, http.StatusOK, b)
}

func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var in model.Libro
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "json inválido", err)
		return
	}

	// Por seguridad/consistencia: el ID lo asigna el sistema
	in.ID = 0

	s.mu.Lock()
	defer s.mu.Unlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	books, err = catalog.AddBook(books, in)
	if err != nil {
		writeError(w, http.StatusBadRequest, "datos inválidos", err)
		return
	}

	if err := s.Repo.Save(books); err != nil {
		writeError(w, http.StatusInternalServerError, "error guardando libros", err)
		return
	}

	// Devuelve el libro recién agregado (último del slice)
	created := books[len(books)-1]
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido", err)
		return
	}

	var cambios model.Libro
	if err := readJSON(r, &cambios); err != nil {
		writeError(w, http.StatusBadRequest, "json inválido", err)
		return
	}
	cambios.ID = id // asegura consistencia

	s.mu.Lock()
	defer s.mu.Unlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	books, err = catalog.UpdateBook(books, id, cambios)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, catalog.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, "no se pudo actualizar", err)
		return
	}

	if err := s.Repo.Save(books); err != nil {
		writeError(w, http.StatusInternalServerError, "error guardando libros", err)
		return
	}

	updated, _ := catalog.GetByID(books, id)
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	books, err = catalog.DeleteBook(books, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "libro no encontrado", err)
		return
	}

	if err := s.Repo.Save(books); err != nil {
		writeError(w, http.StatusInternalServerError, "error guardando libros", err)
		return
	}

	// 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleToggleAvailability(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	books, err = catalog.ToggleAvailability(books, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "libro no encontrado", err)
		return
	}

	if err := s.Repo.Save(books); err != nil {
		writeError(w, http.StatusInternalServerError, "error guardando libros", err)
		return
	}

	updated, _ := catalog.GetByID(books, id)
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleSearchBooks(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		writeError(w, http.StatusBadRequest, "falta parámetro q", errors.New("use ?q=texto"))
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	res := catalog.SearchBooks(books, q)
	writeJSON(w, http.StatusOK, res)
}

type Stats struct {
	Total          int            `json:"total"`
	Disponibles    int            `json:"disponibles"`
	NoDisponibles  int            `json:"noDisponibles"`
	PorCategoria   map[string]int `json:"porCategoria"`
	PorFormato     map[string]int `json:"porFormato"`
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books, err := s.Repo.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error cargando libros", err)
		return
	}

	st := Stats{
		Total:        len(books),
		PorCategoria: map[string]int{},
		PorFormato:   map[string]int{},
	}

	for _, b := range books {
		if b.Disponible {
			st.Disponibles++
		} else {
			st.NoDisponibles++
		}
		st.PorCategoria[b.Categoria]++
		st.PorFormato[b.Formato]++
	}

	writeJSON(w, http.StatusOK, st)
}

// ---- Helpers ----

func parseID(r *http.Request) (int, error) {
	raw := r.PathValue("id") // Go 1.22+
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, errors.New("id debe ser un entero positivo")
	}
	return id, nil
}

func readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string, err error) {
	writeJSON(w, status, map[string]any{
		"error":   message,
		"detail":  err.Error(),
		"status":  status,
		"ts":      time.Now().Format(time.RFC3339),
	})
}