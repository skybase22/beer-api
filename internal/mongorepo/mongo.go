package mongorepo

import (
	"beer-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository common repository
type Repository struct {
}

// NewRepository new repository
func NewRepository() Repository {
	return Repository{}
}

// Create create
func (r *Repository) Create(collection *mongo.Collection, i interface{}) error {
	if m, ok := i.(models.MongoModel); ok {
		m.ID = primitive.NewObjectID()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	i, err := collection.InsertOne(ctx, i)
	if err != nil {
		panic(err)
	}

	return nil
}
