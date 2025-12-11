package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Crear una instancia de Gin
	// Para desarrollo usa gin.Default() que incluye logger y recovery middleware
	router := gin.Default()

	// Ruta GET para obtener notas
	router.GET("/api/notes", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "You got 5 notes",
		})
	})

	// Iniciar el servidor en el puerto 5001
	router.Run(":5001")
}