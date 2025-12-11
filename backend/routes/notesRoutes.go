package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Exportar las rutas de las notas
func SetupNoutesRoutes(router *gin.RouterGroup) {

	// GET - Obtener todas las notas
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "You just fetched the notes",
		})
	})

	// POST - Crear una nueva nota
	router.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Note created successfully",
		})
	})

	// PUT - Actualizar una nota existente
	router.PUT("/:id", func(c *gin.Context) {
		// Obtener el ID de la URL
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "Note updated successfully",
			"id":      id,
		})
	})

	// DELETE - Eliminar una nota existente
	router.DELETE("/:id", func(c *gin.Context) {
		// Obtener el ID de la URL
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "note deleted successfully",
			"id":      id,
		})
	})
}