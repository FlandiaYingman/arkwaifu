package infra

import (
	"context"
	"go.uber.org/fx"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func ProvideFiber(lc fx.Lifecycle, config *Config) (*fiber.App, fiber.Router) {
	app := fiber.New(fiber.Config{
		AppName:      "Arkwaifu",
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
		BodyLimit:    16 * 1024 * 1024,
		UnescapePath: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	})).Use(logger.New())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := app.Listen(config.Address)
				if err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})

	return app, app.Group("api/v1")
}
