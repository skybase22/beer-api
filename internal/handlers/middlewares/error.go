package middlewares

import (
	"beer-api/internal/handlers/render"

	"github.com/gofiber/fiber/v2"
)

// WrapError wrap error
func WrapError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return render.Error(c, err)
		}
		return nil
	}
}
