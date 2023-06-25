package art

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	service *Service
}

func newController(s *Service) *controller {
	c := controller{s}
	return &c
}
func registerController(c *controller, router fiber.Router) {
	if router == nil {
		return
	}

	router.Get("arts", c.GetArts)
	router.Get("arts/:id", c.GetArt)
	router.Put("arts/:id", c.PutArt)

	router.Put("arts/:id/variants/:variation", c.PutVariant)

	router.Get("arts/:id/variants/:variation/content", c.GetVariantContent)
	router.Put("arts/:id/variants/:variation/content", c.PutVariantContent)
}

func (c *controller) GetArts(ctx *fiber.Ctx) error {
	filter := artQueryFilter{}
	err := ctx.QueryParser(&filter)
	if err != nil {
		return err
	}

	var arts []*Art
	if variation := ctx.Query("absent-variation", ""); variation != "" {
		arts, err = c.service.SelectArtsWhoseVariantAbsent(variation)
	} else {
		arts, err = c.service.SelectArts(filter.category)
	}
	if err != nil {
		return err
	}
	return ctx.JSON(arts)
}
func (c *controller) GetArt(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	art, err := c.service.SelectArt(id)
	if err != nil {
		return err
	}

	return ctx.JSON(art)
}
func (c *controller) GetVariants(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	variants, err := c.service.SelectVariants(id)
	if err != nil {
		return err
	}

	return ctx.JSON(variants)
}
func (c *controller) GetVariant(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	kind := ctx.Params("variation")

	variants, err := c.service.SelectVariant(id, kind)
	if err != nil {
		return err
	}

	return ctx.JSON(variants)
}

func (c *controller) PutArt(ctx *fiber.Ctx) error {
	// TODO: Authorization

	id := ctx.Params("id")

	art := new(Art)
	err := ctx.BodyParser(&art)
	if err != nil {
		return errors.Join(fiber.ErrBadRequest, err)
	}
	art.ID = id

	err = c.service.UpsertArts(art)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
func (c *controller) PutVariant(ctx *fiber.Ctx) error {
	// TODO: Authorization

	id := ctx.Params("id")
	kind := ctx.Params("kind")

	variant := new(Variant)
	err := ctx.BodyParser(&variant)
	if err != nil {
		return errors.Join(fiber.ErrBadRequest, err)
	}
	variant.ArtID = id
	variant.Variation = kind

	err = c.service.UpsertVariants(variant)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *controller) GetVariantContent(ctx *fiber.Ctx) error {
	static := VariantContent{}
	err := ctx.ParamsParser(&static)
	if err != nil {
		return err
	}

	err = c.service.TakeStatics(&static)
	if err != nil {
		return err
	}

	ctx.Type("webp")
	return ctx.Send(static.Content)
}
func (c *controller) PutVariantContent(ctx *fiber.Ctx) error {
	static := VariantContent{}
	err := ctx.ParamsParser(&static)
	if err != nil {
		return err
	}
	static.Content = ctx.Body()

	err = c.service.StoreStatics(&static)
	if err != nil {
		return err
	}

	return nil
}

type artQueryFilter struct {
	category *string `query:"category"`
}

type variantQueryFilter struct {
}