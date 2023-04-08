package routes

import (
	ctx "context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"beer-api/internal/core/config"

	"beer-api/internal/handlers/middlewares"
	"beer-api/internal/pkg/beer"
	"beer-api/internal/pkg/healthcheck"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// NewRouter new router
func NewRouter() *fiber.App {
	app := fiber.New()

	app.Use(
		compress.New(),
		requestid.New(),
		cors.New(),
		middlewares.Logger(),
		middlewares.WrapError(),
		middlewares.TransactionPostgresql(func(c *fiber.Ctx) bool {
			return c.Method() == fiber.MethodGet
		}),
	)

	api := app.Group("/api")

	healthCheckEndpoint := healthcheck.NewEndpoint()
	helthCheck := api.Group("/health-check")
	helthCheck.Get("/", healthCheckEndpoint.HealthCheck)

	v1 := api.Group("/v1")
	v1.Use(middlewares.AcceptLanguage())

	v1.Static("/public", "./public")

	beerEndpoint := beer.NewEndpoint()
	beer := v1.Group("/beer")
	beer.Post("/", beerEndpoint.Create)
	beer.Get("/:id", beerEndpoint.GetOne)
	beer.Get("/", beerEndpoint.GetAll)
	beer.Put("/:id", beerEndpoint.Update)
	beer.Delete("/:id", beerEndpoint.Delete)

	api.Use(func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Not Found",
			"status":  404,
		})
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
		defer cancel()

		logrus.Info("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	logrus.Infof("Start server on port: %d ...", config.CF.App.Port)
	err := app.Listen(fmt.Sprintf(":%d", config.CF.App.Port))
	if err != nil {
		logrus.Panic(err)
	}

	return app
}
