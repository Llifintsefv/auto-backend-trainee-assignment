package main

import (
	"auto-backend-trainee-assignment/internal/handlers"
	"auto-backend-trainee-assignment/internal/repository"
	"auto-backend-trainee-assignment/internal/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func InitDB(connStr string) (*sql.DB, error) {
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY,LongUrl TEXT, ShortUrl TEXT)")
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return db, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == ""{
		 log.Fatal("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE must be set")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	

	db, err := InitDB(connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handlers.NewHandler(service)

	r := mux.NewRouter()

	r.HandleFunc("/shorter",handler.ShortenHandler).Methods("POST")
	r.HandleFunc("/{shortURL}",handler.RedirectHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))

}