package main

import (
	"CodeCast_backend/db"
	"CodeCast_backend/modules/anon_sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(`../.env`); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	if err := db.InitDB(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	r := gin.Default()

	// Looks Good
	r.POST("/api/v1/sessions/anon", anon_sessions.CreateAnonSession)
	r.POST("/api/v1/sessions/anon/:code/join", anon_sessions.JoinAnonSession)

	// Still needs work
	r.POST("/api/v1/sessions/anon/:code/snippets", anon_sessions.PushAnonSnippet)
	r.GET("/api/v1/sessions/anon/:code/snippets", anon_sessions.GetAnonSnippets)
	r.POST("/api/v1/sessions/anon/:code/end", anon_sessions.EndAnonSession)

	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
