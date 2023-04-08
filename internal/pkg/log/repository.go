package log

import (
	"beer-api/internal/mongorepo"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository repo interface
type Repository interface {
	Create(collection *mongo.Collection, i interface{}) error
}

type repository struct {
	mongorepo.Repository
}

// NewRepository new repository
func NewRepository() Repository {
	return &repository{
		mongorepo.NewRepository(),
	}
}
