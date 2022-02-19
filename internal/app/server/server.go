package server

import "github.com/gofiber/fiber/v2"

type V0 struct {
	fiber.Router
}

func ProvideV0(app *fiber.App) *V0 {
	r := app.Group("/api/v0")
	return &V0{
		Router: r,
	}
}
