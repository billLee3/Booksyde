package main

import (
	"booksyde/internal/database"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("No DB Url is present. You must add a postgres connection string to your .env file with the var name of DB_URL")
	}
	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("unable to establish a connection with the database")
	}

	dbQueries := database.New(dbConn)
	log.Printf("Got the database queries pointer: %v", dbQueries)
}
