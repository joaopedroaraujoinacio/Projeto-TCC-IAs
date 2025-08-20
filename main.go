package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


type User struct {
	ID 		int			`json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

var db *sql.DB

func initDB() {
	var err error
//db conn
	connStr := "postgresql://localhost:5432/rag_database_test?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL!")
}

//!!!c
func createUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
//!! gin.h
	c.JSON(http.StatusCreated, gin.H{
	"message": "User created successfully",
	"user": user,
	})
}


func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var user User
		//!!!
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
			return			
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}


func main() {
	initDB()
	defer db.Close()
	//!!!

	r := gin.Default()

	r.POST("/users", createUser)
	r.GET("/users,", getUsers)


	r.Run(":8080")
}
