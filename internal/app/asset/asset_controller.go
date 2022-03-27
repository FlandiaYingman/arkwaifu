package asset

import (
	"github.com/flandiayingman/arkwaifu/internal/app/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"strings"
	"time"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return Controller{service}
}
func RegisterController(v0 *server.V0, c Controller) {
	router := v0.
		Group("assets").
		Use(newCompress()).
		Use(newETag()).
		Use(newCache())

	router.Get("/", c.GetAssets)

	router.Get("/kinds", c.GetKinds)
	router.Get("/kinds/:kind", c.GetAssets)

	router.Get("/kinds/:kind/names", c.GetNames)
	router.Get("/kinds/:kind/names/:name", c.GetAssets)

	router.Get("/kinds/:kind/names/:name/variants", c.GetVariants)
	router.Get("/kinds/:kind/names/:name/variants/:variant", c.GetAsset)

	router.Get("/kinds/:kind/names/:name/variants/:variant/file", c.GetAssetFile)
}
func newCompress() fiber.Handler {
	return compress.New(compress.Config{
		Next: func(ctx *fiber.Ctx) bool {
			return strings.HasSuffix(ctx.Path(), "/file")
		},
		Level: compress.LevelBestSpeed,
	})
}
func newETag() fiber.Handler {
	return etag.New(etag.Config{
		Next: func(ctx *fiber.Ctx) bool {
			// skip ETag if it isn't a file request
			return !strings.HasSuffix(ctx.Path(), "/file")
		},
	})
}
func newCache() fiber.Handler {
	return cache.New(cache.Config{
		CacheControl: true,
		Next: func(ctx *fiber.Ctx) bool {
			switch {
			case strings.HasSuffix(ctx.Path(), "/file"):
				return ctx.Response().StatusCode() != fiber.StatusOK
			default:
				return false
			}
		},
		ExpirationGenerator: func(ctx *fiber.Ctx, config *cache.Config) time.Duration {
			switch {
			case strings.HasSuffix(ctx.Path(), "/file"):
				return 24 * time.Hour
			default:
				return 1 * time.Minute
			}
		},
	})
}

func (c *Controller) GetAssets(ctx *fiber.Ctx) error {
	kind, name, variant := parseParams(ctx)
	assets, err := c.service.GetAssets(ctx.Context(), kind, name, variant)
	if err != nil {
		return err
	}
	return ctx.JSON(assets)
}
func (c *Controller) GetAsset(ctx *fiber.Ctx) error {
	kind, name, variant := parseParams(ctx)
	asset, err := c.service.GetAsset(ctx.Context(), kind, name, variant)
	if err != nil {
		return err
	}
	return ctx.JSON(asset)
}
func (c *Controller) GetKinds(ctx *fiber.Ctx) error {
	kinds, err := c.service.GetKinds(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(kinds)
}
func (c *Controller) GetNames(ctx *fiber.Ctx) error {
	kind, _, _ := parseParams(ctx)
	names, err := c.service.GetNames(ctx.Context(), kind)
	if err != nil {
		return err
	}
	return ctx.JSON(names)
}
func (c *Controller) GetVariants(ctx *fiber.Ctx) error {
	kind, name, _ := parseParams(ctx)
	variants, err := c.service.GetVariants(ctx.Context(), kind, name)
	if err != nil {
		return err
	}
	return ctx.JSON(variants)
}
func (c *Controller) GetAssetFile(ctx *fiber.Ctx) error {
	kind, name, variant := parseParams(ctx)
	assetFilePath, err := c.service.GetAssetFilePath(ctx.Context(), kind, name, variant)
	if err != nil {
		return err
	}
	return ctx.SendFile(assetFilePath)
}

func parseParams(ctx *fiber.Ctx) (string, string, string) {
	kind := ctx.Params("kind")
	name := ctx.Params("name")
	variant := ctx.Params("variant")
	return kind, name, variant
}
