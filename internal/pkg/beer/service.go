package beer

import (
	"beer-api/internal/pkg/file"
	"beer-api/internal/pkg/log"
	"fmt"
	"time"

	"beer-api/internal/core/config"
	"beer-api/internal/core/context"
	"beer-api/internal/core/mongodb"
	"beer-api/internal/models"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// Service service interface
type Service interface {
	Create(c *context.Context, request *createRequest) (*models.Beer, error)
	GetOne(c *context.Context, request *getOneRequest) (*models.Beer, error)
	GetAll(c *context.Context, request *GetAllRequest) (*models.Page, error)
	Update(c *context.Context, request *updateRequest) error
	Delete(c *context.Context, request *getOneRequest) error
}

type service struct {
	config         *config.Configs
	result         *config.ReturnResult
	repository     Repository
	fileRepository file.Repository
	logRepository  log.Repository
}

// NewService new service
func NewService() Service {
	return &service{
		config:         config.CF,
		result:         config.RR,
		repository:     NewRepository(),
		fileRepository: file.NewRepository(),
		logRepository:  log.NewRepository(),
	}
}

// Create create
func (s service) Create(c *context.Context, request *createRequest) (*models.Beer, error) {
	fileCreate, err := s.uploadFile(c)
	if err != nil {
		logrus.Errorf("upload file error: %s", err)
		return nil, err
	}

	beer := &models.Beer{}
	_ = copier.CopyWithOption(beer, request, copier.Option{IgnoreEmpty: true})
	if fileCreate != nil {
		beer.FileID = &fileCreate.ID
	}
	err = s.repository.Create(c.GetDatabase(), beer)
	if err != nil {
		logrus.Errorf("create beer error: %s", err)
		return nil, err
	}

	log := &models.Log{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("created beer name %s", beer.Name),
	}
	err = s.logRepository.Create(mongodb.MI.DB.Collection("logs"), log)
	if err != nil {
		logrus.Errorf("create log error: %s", err)
		return nil, err
	}

	return beer, nil
}

// GetOne get one
func (s service) GetOne(c *context.Context, request *getOneRequest) (*models.Beer, error) {
	beer := &models.Beer{}
	err := s.repository.FindOneObjectByField(
		c.GetDatabase().Preload(clause.Associations), "id", request.ID, beer)
	if err != nil {
		logrus.Errorf("find beer by id=%d error: %s", request.ID, err)
		return nil, s.result.Internal.DatabaseNotFound
	}

	return beer, nil
}

// GetAll get all
func (s service) GetAll(c *context.Context, request *GetAllRequest) (*models.Page, error) {
	page, err := s.repository.FindAll(c.GetDatabase().Preload(clause.Associations), request)
	if err != nil {
		logrus.Errorf("find all beer error: %s", err)
		return nil, err
	}

	return page, nil
}

// Update update
func (s service) Update(c *context.Context, request *updateRequest) error {
	beer := &models.Beer{}
	err := s.repository.FindOneObjectByField(
		c.GetDatabase().Preload(clause.Associations), "id", request.ID, beer)
	if err != nil {
		logrus.Errorf("find beer by id=%d error: %s", request.ID, err)
		return s.result.Internal.DatabaseNotFound
	}

	fileCreate, err := s.uploadFile(c)
	if err != nil {
		logrus.Errorf("upload file error: %s", err)
		return err
	}

	_ = copier.CopyWithOption(beer, request, copier.Option{IgnoreEmpty: true})
	if fileCreate != nil {
		beer.FileID = &fileCreate.ID
	}
	err = s.repository.Update(c.GetDatabase(), beer)
	if err != nil {
		logrus.Errorf("update beer error: %s", err)
		return err
	}

	log := &models.Log{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("updated beer id %d", beer.ID),
	}
	err = s.logRepository.Create(mongodb.MI.DB.Collection("logs"), log)
	if err != nil {
		logrus.Errorf("create log error: %s", err)
		return err
	}

	return nil
}

// Update update
func (s service) Delete(c *context.Context, request *getOneRequest) error {
	beer := &models.Beer{}
	err := s.repository.FindOneObjectByField(
		c.GetDatabase().Preload(clause.Associations), "id", request.ID, beer)
	if err != nil {
		logrus.Errorf("find beer by id=%d error: %s", request.ID, err)
		return s.result.Internal.DatabaseNotFound
	}

	err = s.repository.Delete(c.GetDatabase(), beer)
	if err != nil {
		logrus.Errorf("delete beer error: %s", err)
		return err
	}

	log := &models.Log{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("deleted beer id %d", beer.ID),
	}
	err = s.logRepository.Create(mongodb.MI.DB.Collection("logs"), log)
	if err != nil {
		logrus.Errorf("create log error: %s", err)
		return err
	}

	return nil
}
