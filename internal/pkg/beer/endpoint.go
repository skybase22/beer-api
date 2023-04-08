package beer

import (
	"beer-api/internal/core/config"
	"beer-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// Endpoint endpoint interface
type Endpoint interface {
	Create(c *fiber.Ctx) error
	GetOne(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type endpoint struct {
	config  *config.Configs
	result  *config.ReturnResult
	service Service
}

// NewEndpoint new endpoint
func NewEndpoint() Endpoint {
	return &endpoint{
		config:  config.CF,
		result:  config.RR,
		service: NewService(),
	}
}

// Create create beer
func (ep *endpoint) Create(c *fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.Create, &createRequest{})
}

// Get get one
func (ep *endpoint) GetOne(c *fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.GetOne, &getOneRequest{})
}

// GetAll get all
func (ep *endpoint) GetAll(c *fiber.Ctx) error {
	return handlers.ResponseObject(c, ep.service.GetAll, &GetAllRequest{})
}

// Update update
func (ep *endpoint) Update(c *fiber.Ctx) error {
	return handlers.ResponseSuccess(c, ep.service.Update, &updateRequest{})
}

// Delete delete
func (ep *endpoint) Delete(c *fiber.Ctx) error {
	return handlers.ResponseSuccess(c, ep.service.Delete, &getOneRequest{})
}
