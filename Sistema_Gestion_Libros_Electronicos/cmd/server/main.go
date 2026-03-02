package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"sistema_libros/internal/api"
	"sistema_libros/internal/storage"
)

func main() {
	// DB
	conn := os.Getenv("DB_CONN")
	if conn == "" {
		conn = "server=localhost,1433;user id=sa;password=TuNuevaClave123!;database=BibliotecaDB;encrypt=true;trustservercertificate=true"
	}

	repo, err := storage.NewSQLServerStore(conn)
	if err != nil {
		log.Fatalf("Error conectando a SQL Server: %v", err)
	}
	defer repo.Close()

	// API mux (tus rutas /api y /health)
	s := &api.Server{Repo: repo}
	apiMux := api.NewMux(s)

	// Mux principal
	rootMux := http.NewServeMux()

	// 1) API
	rootMux.Handle("/api/", apiMux)
	rootMux.Handle("/health", apiMux)

	// 2) Frontend: servir carpeta web/
	//    IMPORTANTE: "web" debe estar en la raíz del proyecto (como en tu captura)
	rootMux.Handle("/", http.FileServer(http.Dir("web")))

	// Logging
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rootMux.ServeHTTP(w, r)
		fmt.Printf("%s %s (%s)\n", r.Method, r.URL.Path, time.Since(start))
	})

	addr := ":8080"
	fmt.Println("=====================================")
	fmt.Println("Servidor:  http://localhost" + addr + "/")
	fmt.Println("API:       http://localhost" + addr + "/api/libros")
	fmt.Println("Health:    http://localhost" + addr + "/health")
	fmt.Println("=====================================")

	log.Fatal(http.ListenAndServe(addr, handler))
}
