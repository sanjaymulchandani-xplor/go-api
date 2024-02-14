package main

import (
	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var employees = make(map[string]Employee)

func main() {
	r := gin.Default()

	// CRUD operations
	r.POST("/employees", func(c *gin.Context) {
		var newEmployee Employee
		if err := c.ShouldBindJSON(&newEmployee); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		employees[newEmployee.ID] = newEmployee
		c.JSON(201, newEmployee)
	})

	r.GET("/employees", func(c *gin.Context) {
		c.JSON(200, employees)
	})

	r.GET("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")
		if employee, ok := employees[id]; ok {
			c.JSON(200, employee)
			return
		}
		c.JSON(404, gin.H{"error": "Employee not found"})
	})

	r.PUT("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")
		if _, ok := employees[id]; !ok {
			c.JSON(404, gin.H{"error": "Employee not found"})
			return
		}
		var updatedEmployee Employee
		if err := c.ShouldBindJSON(&updatedEmployee); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		employees[id] = updatedEmployee
		c.JSON(200, updatedEmployee)
	})

	r.DELETE("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")
		if _, ok := employees[id]; !ok {
			c.JSON(404, gin.H{"error": "Employee not found"})
			return
		}
		delete(employees, id)
		c.JSON(204, gin.H{})
	})

	r.Run() // go to 0.0.0.0:8080
}
