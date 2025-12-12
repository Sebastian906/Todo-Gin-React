package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Conectar a MongoDB
	config.ConnectDB()

	// Crear una instancia de Gin
	router := gin.Default()

	// Middleware para parsear JSON
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configurar el grupo de rutas para /api/notes
	notesGroup := router.Group("/api/notes")
	routes.SetupNoutesRoutes(notesGroup)

	// Obtener el puerto desde variables de entorno o usar 5001 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	// Iniciar el servidor
	router.Run(":" + port)
}