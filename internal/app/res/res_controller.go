package res

import (
	"arkwaifu/internal/app/server"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return Controller{service}
}

func (c *Controller) GetImages(ctx *fiber.Ctx) error {
	images, err := c.service.GetImages(ctx.Context())
	if err != nil {
		return err
	}
	var imageNames []string
	linq.From(images).Select(func(i interface{}) interface{} {
		return i.(Resource).Name
	}).ToSlice(&imageNames)
	return ctx.JSON(imageNames)
}

func (c *Controller) GetImageByName(ctx *fiber.Ctx, name string, resTypeStr string) error {
	resType, err := ResTypeFromString(resTypeStr)
	if err != nil {
		return err
	}
	image, err := c.service.GetImageByName(ctx.Context(), name, resType)
	if err != nil {
		return err
	}
	if image != nil {
		return ctx.SendFile(image.Path)
	} else {
		return ctx.
			Status(404).
			SendString(fmt.Sprintf("image by name %s cannot be found", name))
	}
}

func (c *Controller) GetBackgrounds(ctx *fiber.Ctx) error {
	background, err := c.service.GetBackgrounds(ctx.Context())
	if err != nil {
		return err
	}
	var imageNames []string
	linq.From(background).Select(func(i interface{}) interface{} {
		return i.(Resource).Name
	}).ToSlice(&imageNames)
	return ctx.JSON(imageNames)
}

func (c *Controller) GetBackgroundByName(ctx *fiber.Ctx, name string, resTypeStr string) error {
	resType, err := ResTypeFromString(resTypeStr)
	if err != nil {
		return err
	}
	background, err := c.service.GetBackgroundByName(ctx.Context(), name, resType)
	if err != nil {
		return err
	}
	if background != nil {
		return ctx.SendFile(background.Path)
	} else {
		return ctx.
			Status(404).
			SendString(fmt.Sprintf("background by name %s cannot be found", name))
	}
}

func RegisterController(v0 *server.V0, c Controller) {
	v0.Get("resources/images", c.GetImages)
	v0.Get("resources/images/:imageName", func(ctx *fiber.Ctx) error {
		return c.GetImageByName(ctx, ctx.Params("imageName"), ctx.Query("resType"))
	})
	v0.Get("resources/backgrounds", c.GetBackgrounds)
	v0.Get("resources/backgrounds/:backgroundName", func(ctx *fiber.Ctx) error {
		return c.GetBackgroundByName(ctx, ctx.Params("backgroundName"), ctx.Query("resType"))
	})
}
