package main

import (
	"fmt"
	"os"

	"sistema_libros/internal/storage"
	"sistema_libros/internal/ui"
)

func main() {
	// IMPORTANTE:
	// - En tu SSMS el cifrado está "Obligatorio" y marcas "Confiar en el certificado".
	// - Por eso agregamos encrypt=true;trustservercertificate=true
	//
	// Cambia TU_PASSWORD por la contraseña real del usuario sa.
	conn := "server=localhost,1433;user id=sa;password=TuNuevaClave123!;database=BibliotecaDB;encrypt=true;trustservercertificate=true"

	repo, err := storage.NewSQLServerStore(conn)
	if err != nil {
		fmt.Println("Error conectando a SQL Server:", err)
		os.Exit(1)
	}
	defer repo.Close()

	ui.Run(repo)
}
