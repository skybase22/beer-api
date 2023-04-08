package repositories

import (
	"context"
	"fmt"
	"math"
	"strings"

	"time"

	"beer-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository common repository
type Repository struct {
}

// NewRepository new repository
func NewRepository() Repository {
	return Repository{}
}

// DefaultContext default context
func (r *Repository) DefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*10)
}

// FindOneObjectByField find one
func (r *Repository) FindOneObjectByField(db *gorm.DB, field string, value interface{}, i interface{}) error {
	return db.Where(fmt.Sprintf("%s = ?", field), value).First(i).Error
}

// Create create
func (r *Repository) Create(db *gorm.DB, i interface{}) error {
	return db.Omit(clause.Associations).Create(i).Error
}

// Update update
func (r *Repository) Update(db *gorm.DB, i interface{}) error {
	return db.Omit(clause.Associations).Save(i).Error
}

// Delete update stamp deleted_at
func (r *Repository) Delete(db *gorm.DB, i interface{}) error {
	return db.Omit(clause.Associations).Delete(i).Error
}

// PageForm page info interface
type PageForm interface {
	GetPage() int
	GetSize() int
	GetQuery() string
	GetSort() string
	GetSorts() []string
	GetReverse() bool
	GetReverses() []bool
	GetOrderBy() string
}

const (
	// DefaultPage default page in page query
	DefaultPage int = 1
	// DefaultSize default size in page query
	DefaultSize int = 20
)

// FindAllAndPageInformation get page information
func (r *Repository) FindAllAndPageInformation(db *gorm.DB, pageForm PageForm, entities interface{}) (*models.PageInformation, error) {
	page := pageForm.GetPage()
	if pageForm.GetPage() < 1 {
		page = DefaultPage
	}

	limit := pageForm.GetSize()
	if pageForm.GetSize() == 0 {
		limit = DefaultSize
	}

	var count int64
	db.Model(entities).Count(&count)

	if pageForm.GetOrderBy() != "" {
		db = db.Order(pageForm.GetOrderBy())
	} else if pageForm.GetSort() != "" {
		order := pageForm.GetSort()
		if pageForm.GetReverse() {
			split := strings.Split(order, " ")
			if len(split) == 1 {
				order = order + " desc"
			}
		}
		db = db.Order(order)
	} else if len(pageForm.GetSorts()) > 0 {
		sorts := pageForm.GetSorts()
		for i, sort := range sorts {
			reverse := pageForm.GetReverses()
			if reverse[i] {
				split := strings.Split(sort, " ")
				if len(split) == 1 {
					sort = sort + " desc"
				}
			}
			db = db.Order(sort)
		}
	} else {
		db = db.Order("id DESC")
	}

	var offset int
	if page != 1 {
		offset = (page - 1) * limit
	}

	if err := db.
		Limit(limit).
		Offset(offset).
		Find(entities).Error; err != nil {
		return nil, err
	}

	return &models.PageInformation{
		Page:                  page,
		Size:                  limit,
		TotalNumberOfEntities: count,
		TotalNumberOfPages:    int(math.Ceil(float64(count) / float64(limit))),
	}, nil
}
