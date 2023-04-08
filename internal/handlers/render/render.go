package render

import (
	"beer-api/internal/core/config"

	"github.com/gofiber/fiber/v2"
)

// JSON render json to client
func JSON(c *fiber.Ctx, response interface{}) error {
	return c.
		Status(config.RR.Internal.Success.HTTPStatusCode()).
		JSON(response)
}

// Byte render byte to client
func Byte(c *fiber.Ctx, bytes []byte) error {
	_, err := c.Status(config.RR.Internal.Success.HTTPStatusCode()).
		Write(bytes)

	return err
}

// Error render error to client
func Error(c *fiber.Ctx, err error) error {
	errMsg := config.RR.Internal.ConnectionError
	if locErr, ok := err.(config.Result); ok {
		errMsg = locErr
	}

	return c.
		Status(errMsg.HTTPStatusCode()).
		JSON(errMsg.WithLocale(c))
}
