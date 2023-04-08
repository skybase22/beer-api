package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

const (
	// Location
	Location string = "location"
	// LocationUTC
	LocationUTC string = "UTC"
)

// Model base model
type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id" copier:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Modeler model interface
type Modeler interface {
	GetID() uint
}

// GetID get id
func (m Model) GetID() uint {
	return m.ID
}

// Model base model
type MongoModel struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}
