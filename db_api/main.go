// STILL IN PROGRESS

package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Item struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Colour   string `json:"colour"`
}

func initDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func createItem(c *gin.Context) {
	var newItem Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `INSERT INTO items (name) VALUES ($1) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, newItem.Name).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insertion failed"})
		return
	}
	newItem.ID = id

	c.JSON(http.StatusCreated, newItem)
}

func getItems(c *gin.Context) {
	var items []Item

	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error"})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func getItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item Item

	sqlStatement := `SELECT id, name FROM items WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	if err := row.Scan(&item.ID, &item.Name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func updateItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `UPDATE items SET name=$1 WHERE id=$2`
	_, err := db.Exec(sqlStatement, item.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

func deleteItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	sqlStatement := `DELETE FROM items WHERE id=$1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func main() {
	databaseURL := "postgresql://postgres:dFbce*E-a3CgfFa3f-4GBE-25e-*d1FE@viaduct.proxy.rlwy.net:52559/railway"
	initDB(databaseURL)

	router := gin.Default()

	router.POST("/items", createItem)
	router.GET("/items", getItems)
	router.GET("/items/:id", getItem)
	router.PUT("/items/:id", updateItem)
	router.DELETE("/items/:id", deleteItem)

	router.Run(":8080")
}
