package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// Esquema de la nota de MongoDB
type Note struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title" binding:"required"`
	Content   string        `json:"content" bson:"content" binding:"required"`
	Completed bool          `json:"completed" bson:"completed"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// Validar datos ingresados del usuario
type CreateNoteInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// Actualizar datos de la nota
type UpdateNoteInput struct {
	Title     *string `json:"title"`
	Content   *string `json:"content"`
	Completed *bool   `json:"completed"`
}