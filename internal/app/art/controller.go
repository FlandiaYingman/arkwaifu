package art

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	router.Use("arts", c.Authenticator)

	router.Get("arts", c.GetArts)
	router.Get("arts/:id", c.GetArt)
	router.Put("arts/:id", c.PutArt)

	router.Put("arts/:id/variants/:variation", c.PutVariant)

	router.Get("arts/:id/variants/:variation/content", c.GetContent)
	router.Put("arts/:id/variants/:variation/content", c.PutContent)
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

type ContentParams struct {
	ID        string `param:"id"`
	Variation string `param:"variation"`
}

func (c *controller) GetContent(ctx *fiber.Ctx) error {
	params := ContentParams{}
	err := ctx.ParamsParser(&params)
	if err != nil {
		return err
	}

	content, err := c.service.TakeContent(params.ID, params.Variation)
	if err != nil {
		return err
	}

	ctx.Type("webp")
	return ctx.Send(content)
}
func (c *controller) PutContent(ctx *fiber.Ctx) error {
	params := ContentParams{}
	err := ctx.ParamsParser(&params)
	if err != nil {
		return err
	}

	err = c.service.StoreContent(params.ID, params.Variation, ctx.Body())
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

func (c *controller) Authenticator(ctx *fiber.Ctx) error {
	if c.SkipAuthentication(ctx) {
		return ctx.Next()
	}
	return c.Authenticate(ctx)
}
func (c *controller) SkipAuthentication(ctx *fiber.Ctx) bool {
	return ctx.Method() != fiber.MethodPut
}

func (c *controller) Authenticate(ctx *fiber.Ctx) error {
	idStr := ctx.Query("user", "")
	if idStr == "" {
		return ctx.
			Status(fiber.StatusUnauthorized).
			SendString("no user credential provided")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			SendString(fmt.Sprintf("cannot parse id of user: %s", idStr))
	}
	user := c.service.Authenticate(id)
	if user == nil {
		return ctx.
			Status(fiber.StatusUnauthorized).
			SendString(fmt.Sprintf("cannot find user with id: %s", id))
	}
	return ctx.Next()
}
