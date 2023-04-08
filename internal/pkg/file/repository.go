package file

import (
	"beer-api/internal/repositories"

	"gorm.io/gorm"
)

// Repository repo interface
type Repository interface {
	Create(db *gorm.DB, i interface{}) error
	FindOneObjectByField(db *gorm.DB, field string, value interface{}, i interface{}) error
}

type repository struct {
	repositories.Repository
}

// NewRepository new repository
func NewRepository() Repository {
	return &repository{
		repositories.NewRepository(),
	}
}
