package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//inisialiasai Gin
	router := gin.Default()

	//membuat route dengan method GET
	router.GET("/", func(c *gin.Context) {

		//return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World321",
		})
	})

	// Route POST
	router.POST("/post", func(c *gin.Context) {
		// Struct untuk binding JSON input
		type Input struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		var input Input
		// Parse JSON input dan validasi
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Return response JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "Data received successfully",
			"name":    input.Name,
			"email":   input.Email,
		})
	})

	//mulai server dengan port 3000
	router.Run(":3000")
}