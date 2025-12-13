package routes

import (
	"backend/controllers"
	"backend/middleware"
	"time"
	"github.com/gin-gonic/gin"
)

// Exportar las rutas de las notas
func SetupNoutesRoutes(router *gin.RouterGroup) {

	// GET - Obtener todas las notas
	router.GET("/", middleware.RateLimitWithConfig(20, 20*time.Second), controllers.GetAllNotes)

	// GET - Obtener una nota por ID
	router.GET("/:id", middleware.RateLimitWithConfig(10, 20*time.Second), controllers.GetNoteById)

	// POST - Crear una nueva nota
	router.POST("/", middleware.RateLimitWithConfig(5, 20*time.Second), controllers.CreateNote)

	// PUT - Actualizar una nota por ID
	router.PUT("/:id", middleware.RateLimitWithConfig(5, 20*time.Second), controllers.UpdateNote)

	// DELETE - Eliminar una nota por ID
	router.DELETE("/:id", middleware.RateLimitWithConfig(5, 20*time.Second), controllers.DeleteNote)
}