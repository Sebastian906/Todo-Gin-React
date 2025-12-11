package main

import (
	"backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Crear una instancia de Gin
	router := gin.Default()

	// Configurar grupo de rutas para /api/notes
	notesGroup := router.Group("/api/notes")
	routes.SetupNoutesRoutes(notesGroup)

	// Iniciar el servidor en el puerto 5001
	router.Run(":5001")
}