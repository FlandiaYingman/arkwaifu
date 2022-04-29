package infra

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func ProvideFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Ark Waifu",
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
		BodyLimit:    16 * 1024 * 1024,
		UnescapePath: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	})).Use(logger.New())
	return app
}
