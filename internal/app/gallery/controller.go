package gallery

import (
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	service *Service
}

func newController(service *Service) *controller {
	return &controller{service}
}
func registerController(c *controller, router fiber.Router) {
	if router == nil {
		return
	}
	router.Get(":server/galleries", c.ListGalleries)
	router.Get(":server/galleries/:id", c.GetGalleryByID)
	router.Get(":server/gallery-arts", c.ListGalleryArt)
	router.Get(":server/gallery-arts/:id", c.GetGalleryArtByID)
}

func (c *controller) ListGalleries(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	galleries, err := c.service.ListGalleries(server)
	if err != nil {
		return err
	}

	return ctx.JSON(galleries)
}
func (c *controller) GetGalleryByID(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	gallery, err := c.service.GetGalleryByID(server, id)

	return ctx.JSON(gallery)
}

func (c *controller) ListGalleryArt(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	galleryArts, err := c.service.ListGalleryArts(server)
	if err != nil {
		return err
	}

	return ctx.JSON(galleryArts)
}
func (c *controller) GetGalleryArtByID(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	galleryArt, err := c.service.GetGalleryArtByID(server, id)
	if err != nil {
		return err
	}

	return ctx.JSON(galleryArt)
}
