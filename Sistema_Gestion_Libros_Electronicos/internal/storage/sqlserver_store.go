package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"sistema_libros/internal/model"
)

// SQLServerStore implementa model.BookRepository usando SQL Server.
// Mantiene tu arquitectura actual: Load() trae todo a memoria y Save() guarda todo.
type SQLServerStore struct {
	db *sql.DB
}

// NewSQLServerStore crea el store y valida conexión.
func NewSQLServerStore(connString string) (*SQLServerStore, error) {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SQLServerStore{db: db}, nil
}

func (s *SQLServerStore) Close() error {
	return s.db.Close()
}

// Load trae todos los libros desde dbo.Libros.
func (s *SQLServerStore) Load() ([]model.Libro, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `
		SELECT Id, Titulo, Autor, Categoria, Anio, Formato, UrlArchivo, Disponible
		FROM dbo.Libros
		ORDER BY Id;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Libro
	for rows.Next() {
		var b model.Libro
		if err := rows.Scan(&b.ID, &b.Titulo, &b.Autor, &b.Categoria, &b.Anio, &b.Formato, &b.URLArchivo, &b.Disponible); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// Save guarda toda la lista. Para mantener tu diseño actual:
//// 1) transacción
//// 2) DELETE de la tabla
//// 3) IDENTITY_INSERT ON
//// 4) inserta todos los libros respetando b.ID
//// 5) IDENTITY_INSERT OFF
func (s *SQLServerStore) Save(books []model.Libro) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1) vacía la tabla
	if _, err := tx.ExecContext(ctx, `DELETE FROM dbo.Libros;`); err != nil {
		return err
	}

	// 2) permitir insertar Id explícito
	if _, err := tx.ExecContext(ctx, `SET IDENTITY_INSERT dbo.Libros ON;`); err != nil {
		return err
	}

	// 3) insertar todo
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO dbo.Libros (Id, Titulo, Autor, Categoria, Anio, Formato, UrlArchivo, Disponible)
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, b := range books {
		if _, err := stmt.ExecContext(ctx,
			b.ID, b.Titulo, b.Autor, b.Categoria, b.Anio, b.Formato, b.URLArchivo, b.Disponible,
		); err != nil {
			return fmt.Errorf("insert falló (id=%d): %w", b.ID, err)
		}
	}

	// 4) apagar identity_insert
	if _, err := tx.ExecContext(ctx, `SET IDENTITY_INSERT dbo.Libros OFF;`); err != nil {
		return err
	}

	// 5) commit
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}