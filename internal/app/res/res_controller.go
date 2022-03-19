package res

import (
	"github.com/flandiayingman/arkwaifu/internal/app/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"time"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return Controller{service}
}

func (c *Controller) GetAssets(ctx *fiber.Ctx, variantStr string, kindStr string) error {
	kind, err := ParseKind(kindStr)
	if err != nil {
		return err
	}
	variant, err := ParseVariant(variantStr)
	if err != nil {
		return err
	}
	assets, err := c.service.GetAssets(variant, kind)
	if err != nil {
		return err
	}
	return ctx.JSON(assets)
}

func (c *Controller) GetAssetsByID(ctx *fiber.Ctx, variantStr string, kindStr string, id string) error {
	kind, err := ParseKind(kindStr)
	if err != nil {
		return err
	}
	variant, err := ParseVariant(variantStr)
	if err != nil {
		return err
	}
	image, err := c.service.GetAssetByID(id, variant, kind)
	if err != nil {
		return err
	}
	return ctx.SendFile(image)
}

func RegisterController(v0 *server.V0, c Controller) {
	router := v0.
		Group("assets").
		Use(cache.New(cache.Config{
			Expiration:   24 * time.Hour,
			CacheControl: true,
		})).
		Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))
	router.Get(":variant/:kind", func(ctx *fiber.Ctx) error {
		return c.GetAssets(ctx, ctx.Params("variant"), ctx.Params("kind"))
	})
	router.Get(":variant/:kind/:id", func(ctx *fiber.Ctx) error {
		return c.GetAssetsByID(ctx, ctx.Params("variant"), ctx.Params("kind"), ctx.Params("id"))
	})
}
