package asset

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/flandiayingman/arkwaifu/internal/app/server"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/samber/lo"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service}
}
func RegisterController(v0 *server.V0, c *Controller) {
	router := v0.
		Group("asset").
		Use(newCacheMiddleware()).
		Use(newETagMiddleware())

	router.Get("/assets", c.GetAssets)
	router.Get("/assets/:kind", c.GetAssets)
	router.Get("/assets/:kind/:name", c.GetAsset)

	router.Get("/variants", c.GetVariants)
	router.Get("/variants/:kind/:name", c.GetVariants)
	router.Get("/variants/:kind/:name/:variant", c.GetVariant)
	router.Get("/variants/:kind/:name/:variant/file", c.GetVariantFile)

	router.Get("/kind-names", c.GetKindNames)
	router.Get("/variant-names", c.GetVariantNames)

	router.Post("/variants/:kind/:name/:variant", c.PostVariant)
}

func newCacheMiddleware() fiber.Handler {
	return cache.New(cache.Config{
		Next: doSkipCacheMiddleware,
	})
}

func newETagMiddleware() fiber.Handler {
	return etag.New(etag.Config{
		Next: doSkipETagMiddleware,
	})
}

// doSkipCacheMiddleware skips (returns true) Cache middleware when any of the following is true:
//
// 1. the client is requesting a file.
// 2. the request contains header `Cache-Control: no-cache`
func doSkipCacheMiddleware(ctx *fiber.Ctx) bool {
	return strings.HasSuffix(ctx.Path(), "/file") ||
		strings.Contains(ctx.GetReqHeaders()["Cache-Control"], "no-cache")
}

// doSkipETagMiddleware skips (returns true) E-Tag middleware when any of the following is true:
//
// 1. the client isn't requesting a file.
func doSkipETagMiddleware(ctx *fiber.Ctx) bool {
	return !strings.HasSuffix(ctx.Path(), "/file")
}

func (c *Controller) GetAsset(ctx *fiber.Ctx) error {
	kind := ctx.Params("kind")
	name := ctx.Params("name")
	if kind == "" || name == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("kind and name are required")
	}

	asset, err := c.service.GetAsset(ctx.Context(), kind, name)
	if err != nil {
		return err
	}
	if asset != nil {
		return ctx.JSON(asset)
	} else {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
}
func (c *Controller) GetAssets(ctx *fiber.Ctx) error {
	kind := ctx.Params("kind")
	assets, err := c.service.GetAssets(ctx.Context(), kind)
	if err != nil {
		return err
	}
	return ctx.JSON(assets)
}

func (c *Controller) GetVariant(ctx *fiber.Ctx) error {
	kind := ctx.Params("kind")
	name := ctx.Params("name")
	variant := ctx.Params("variant")
	if kind == "" || name == "" || variant == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("kind, name and variant are required")
	}

	v, err := c.service.GetVariant(ctx.Context(), kind, name, variant)
	if err != nil {
		return err
	}
	if v != nil {
		return ctx.JSON(v)
	} else {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
}
func (c *Controller) GetVariants(ctx *fiber.Ctx) error {
	kind := ctx.Params("kind")
	name := ctx.Params("name")
	if kind == "" || name == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("kind and name are required")
	}

	v, err := c.service.GetVariants(ctx.Context(), kind, name)
	if err != nil {
		return err
	}
	return ctx.JSON(v)
}
func (c *Controller) GetVariantFile(ctx *fiber.Ctx) error {
	kind := ctx.Params("kind")
	name := ctx.Params("name")
	variant := ctx.Params("variant")
	if kind == "" || name == "" || variant == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("kind, name and variant are required")
	}

	v, err := c.service.GetVariant(ctx.Context(), kind, name, variant)
	vFile, err := c.service.GetVariantFile(ctx.Context(), kind, name, variant)
	if err != nil {
		return err
	}
	if vFile != nil {
		ctx.Attachment(v.Filename)
		return ctx.SendFile(escapePath(*vFile))
	} else {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
}

func (c *Controller) GetKindNames(ctx *fiber.Ctx) error {
	kindNames, err := c.service.GetKindNames(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(kindNames)
}
func (c *Controller) GetVariantNames(ctx *fiber.Ctx) error {
	variantNames, err := c.service.GetVariantNames(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(variantNames)
}

func (c *Controller) PostVariant(ctx *fiber.Ctx) error {
	kind := ctx.Params("kind")
	name := ctx.Params("name")
	variant := ctx.Params("variant")
	if kind == "" || name == "" || variant == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("kind, name and variant are required")
	}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	file, err := formFile.Open()
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	v := Variant{
		Variant:  variant,
		Filename: formFile.Filename,
		Asset: &Asset{
			Kind:     kind,
			Name:     name,
			Variants: nil,
		},
	}
	err = c.service.PostVariant(ctx.Context(), kind, name, v, file)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func escapePath(path string) string {
	splits := pathutil.Splits(path)
	return filepath.Join(lo.Map(splits, func(s string, i int) string {
		if i == 0 && strings.Contains(s, "/") {
			return s
		}
		return url.PathEscape(s)
	})...)
}
