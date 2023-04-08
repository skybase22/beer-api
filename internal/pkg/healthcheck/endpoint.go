package healthcheck

import (
	"beer-api/internal/core/config"
	"beer-api/internal/handlers/render"
	"beer-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Endpoint endpoint interface
type Endpoint interface {
	HealthCheck(c *fiber.Ctx) error
}

type endpoint struct {
	config *config.Configs
	result *config.ReturnResult
}

// NewEndpoint new endpoint
func NewEndpoint() Endpoint {
	return &endpoint{
		config: config.CF,
		result: config.RR,
	}
}

// Get health check
func (ep *endpoint) HealthCheck(c *fiber.Ctx) error {
	return render.JSON(c, models.NewSuccessMessage())
}
