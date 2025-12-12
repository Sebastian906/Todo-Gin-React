package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Obtener todas las notas
func GetAllNotes(c *gin.Context) {
	collection := config.GetCollection("notes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var notes []models.Note

	// Opciones para ordenar por fecha de creación descendente
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		// Log del error en el servidor
		println("Error in getAllNotes controller:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &notes); err != nil {
		println("Error in getAllNotes controller:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Si no hay notas, devolver un array vacío
	if notes == nil {
		notes = []models.Note{}
	}

	c.JSON(http.StatusOK, notes)
}

// Obtener una nota por ID
func GetNoteById(c *gin.Context) {
	collection := config.GetCollection("notes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var note models.Note
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not not found"})
			return
		}
		println("Error in getNoteById controller:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// Crear una nueva nota
func CreateNote(c *gin.Context) {
	collection := config.GetCollection("notes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input models.CreateNoteInput

	// Validar el JSON de entrada 
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Crear la nueva nota
	newNote := models.Note{
		ID:        bson.NewObjectID(),
		Title:     input.Title,
		Content:   input.Content,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insertar en la base de datos
	_, err := collection.InsertOne(ctx, newNote)
	if err != nil {
		println("Error in CreateNote controller:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Note created successfully"})
}

// Actualizar una nota existente
func UpdateNote(c *gin.Context) {
	collection := config.GetCollection("notes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Actualizar en la base de datos
	update := bson.M{
		"$set": bson.M{
			"title":     input.Title,
			"content":   input.Content,
			"updatedAt": time.Now(),
		},
	}

	// Opciones para devolver el documento actualizado
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	
	var updatedNote models.Note
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		opts,
	).Decode(&updatedNote)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Note not found"})
			return
		}
		println("Error in updateNote controller:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, updatedNote)
}

// Eliminar una nota
func DeleteNote(c *gin.Context) {
	collection := config.GetCollection("notes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	// Buscar y eliminar 
	var deletedNote models.Note
	err = collection.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&deletedNote)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Note not found"})
			return
		}
		println("Error in deleteNote controller:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
