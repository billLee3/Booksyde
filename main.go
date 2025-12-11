package main

import (
	"booksyde/internal/database"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type CreateUserRequest struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"passwordhash" binding:"required"`
	Subscribed   bool   `json:"subscribed" binding:"required"`
	BirthMonth   string `json:"birth_month" binding:"required"`
	BirthYear    int32  `json:"birth_year" binding:"required"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Subscribed bool      `json:"subscribed"`
}

type GetUserResponse struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordhash"`
	Subscribed   bool      `json:"subscribed"`
	BirthMonth   string    `json:"birth_month"`
	BirthYear    int32     `json:"birth_year"`
}

var users = []GetUserResponse{
	{
		ID:           uuid.New(),
		FirstName:    "Alexandra",
		LastName:     "Chen",
		Email:        "alexandra.c@example.com",
		PasswordHash: "$2a$10$AbCdEfGhIjKlMnOpQrStUvWxYz0123456789./248fHwV3n",
		Subscribed:   true,
		BirthMonth:   time.Month(11).String(),
		BirthYear:    1995,
	},
	{
		ID:        uuid.New(), // Generates a new unique UUID
		FirstName: "Alex",
		LastName:  "Chang",
		Email:     "alex.chang@example.com",
		// A common bcrypt hash prefix and structure
		PasswordHash: "$3a$10$AbCdEfGhIjKlMnOpQrStUvWxYz0123456789./248fHwV3n",
		Subscribed:   true,
		BirthMonth:   time.Month(11).String(), // November
		BirthYear:    1995,
	},
	{
		ID:        uuid.New(), // Generates a new unique UUID
		FirstName: "Lex",
		LastName:  "Chen",
		Email:     "lex.c@example.com",
		// A common bcrypt hash prefix and structure
		PasswordHash: "$3a$10$AbCdEfGhIjKlMnOpQrStUvWxYz0123456789./248fHwV3n",
		Subscribed:   false,
		BirthMonth:   time.Month(11).String(), // November
		BirthYear:    1995,
	},
}

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

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/users", getUsers)

	router.POST("/user", func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Fatalf("unable to bind: %v", err.Error())
		}

		dbUser, err := dbQueries.CreateUser(c, database.CreateUserParams{
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Email:        req.Email,
			PasswordHash: req.PasswordHash,
			Subscribed:   req.Subscribed,
			BirthMonth:   req.BirthMonth,
			BirthYear:    req.BirthYear,
		})

		if err != nil {
			log.Fatalf("unable to create user on database: %v", err)
		}

		res := UserResponse{
			ID:         dbUser.ID,
			Email:      dbUser.Email,
			Subscribed: dbUser.Subscribed,
		}

		mes := fmt.Sprintf("ID: %v |Email: %v |Subscribed: %v", res.ID, res.Email, res.Subscribed)

		c.JSON(200, gin.H{
			"message": mes,
		})
	})

	router.Run()
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}
