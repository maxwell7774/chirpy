package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/maxwell7774/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
    port := ":8080"
    filepathRoot := http.Dir(".")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL must be set")
    }

    platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatal("PLATFORM must be set")
    }
    
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldn't open a database connection")
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	cfg := &apiConfig{
		db:       dbQueries,
		platform: platform,
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(filepathRoot))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

    mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGet)
    mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
    mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpGet)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	server := http.Server{
		Handler: mux,
		Addr:    port,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
