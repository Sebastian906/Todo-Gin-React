package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

// Exportar las rutas de las notas
func SetupNoutesRoutes(router *gin.RouterGroup) {

	// GET - Obtener todas las notas
	router.GET("/", controllers.GetAllNotes)

	// GET - Obtener una nota por ID
	router.GET("/:id", controllers.GetNoteById)

	// POST - Crear una nueva nota
	router.POST("/", controllers.CreateNote)

	// PUT - Actualizar una nota por ID
	router.PUT("/:id", controllers.UpdateNote)

	// DELETE - Eliminar una nota por ID
	router.DELETE("/:id", controllers.DeleteNote)
}