package beer

import (
	"beer-api/internal/models"
	"beer-api/internal/repositories"
	"fmt"

	"gorm.io/gorm"
)

// Repository repo interface
type Repository interface {
	FindOneObjectByField(db *gorm.DB, field string, value interface{}, i interface{}) error
	Create(db *gorm.DB, i interface{}) error
	Update(db *gorm.DB, i interface{}) error
	Delete(db *gorm.DB, i interface{}) error
	FindAll(db *gorm.DB, query *GetAllRequest) (*models.Page, error)
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

// FindAll find all
func (r *repository) FindAll(db *gorm.DB, query *GetAllRequest) (*models.Page, error) {
	entities := []*models.Beer{}
	pageInfo, err := r.
		FindAllAndPageInformation(
			r.query(db, query), &query.PageForm, &entities)
	if err != nil {
		return nil, err
	}

	return models.NewPage(pageInfo, entities), nil
}

func (r *repository) query(db *gorm.DB, query *GetAllRequest) *gorm.DB {
	if query.Query != "" {
		db = db.Where("name ILIKE ?",
			fmt.Sprintf("%%%s%%", query.Query))
	}

	return db
}
