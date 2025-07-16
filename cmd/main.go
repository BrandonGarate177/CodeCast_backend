package main

import (
	"CodeCast_backend/db"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(`../.env`); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	if err := db.InitDB(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
