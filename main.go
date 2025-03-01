package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// Get database connection string
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/globetrotter"
	}

	pendingMigration := os.Getenv("PENDING_MIGRATION")

	// Connect to database
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database connection string: %v\n", err)
	}

	db, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	if pendingMigration == "true" {
		MigrateDestinations()
	}

	// Setup Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := router.Group("/api")
	{
		// User routes
		api.POST("/users/register", registerUser)
		api.POST("/users/login", loginUser)

		// Game routes
		api.GET("/game/question", getRandomQuestion)
		api.POST("/game/answer", submitAnswer)
		api.POST("/game/reveal-clue", revealClue)

		// Challenge routes
		api.POST("/challenges/create", createChallenge)
		api.GET("/challenges/:id", getChallenge)
		api.POST("/challenges/end", endChallenge)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	router.Run(fmt.Sprintf(":%s", port))
}
