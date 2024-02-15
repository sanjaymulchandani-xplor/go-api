package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
	Status    int    `json:"status"`
}

var db *sql.DB
var err error

func main() {
	// Replace "your_username:your_password@/test_db" with your actual database credentials
	db, err = sql.Open("mysql", "root:poT*Hl449)s<@tcp(localhost)/test_users")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/user_details", getUsers)
	r.GET("/user_details/:user_id", getUser)
	r.POST("/user_details", createUser)
	r.PUT("/user_details/:user_id", updateUser)
	r.DELETE("/user_details/:user_id", deleteUser)

	r.Run(":8080")
}

func getUsers(c *gin.Context) {
	var users []User
	rows, err := db.Query("SELECT user_id, username, first_name, last_name, gender, password, status FROM user_details")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.Password, &user.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	userID := c.Param("user_id")
	var user User
	err := db.QueryRow("SELECT user_id, username, first_name, last_name, gender, password, status FROM user_details WHERE user_id = ?", userID).Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Gender, &user.Password, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec("INSERT INTO user_details (username, first_name, last_name, gender, password, status) VALUES (?, ?, ?, ?, ?, ?)", user.Username, user.FirstName, user.LastName, user.Gender, user.Password, user.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	userID := c.Param("user_id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec("UPDATE user_details SET username = ?, first_name = ?, last_name = ?, gender = ?, password = ?, status = ? WHERE user_id = ?", user.Username, user.FirstName, user.LastName, user.Gender, user.Password, user.Status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func deleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	_, err := db.Exec("DELETE FROM user_details WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
