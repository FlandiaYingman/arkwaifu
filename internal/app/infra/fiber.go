package infra

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"time"
)

func ProvideFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Ark Waifu",
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
	})
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	return app
}
