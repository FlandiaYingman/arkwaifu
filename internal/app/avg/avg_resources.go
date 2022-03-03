package avg

import (
	"github.com/flandiayingman/arkwaifu/internal/app/server"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return Controller{service}
}

func (c *Controller) GetGroups(ctx *fiber.Ctx) error {
	groups, err := c.service.GetGroups(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(groups)
}

func (c *Controller) GetGroupByID(ctx *fiber.Ctx) error {
	groupID := ctx.Params("groupID")
	group, err := c.service.GetGroupByID(ctx.Context(), groupID)
	if err != nil {
		return err
	}
	return ctx.JSON(group)
}

func (c *Controller) GetStories(ctx *fiber.Ctx) error {
	stories, err := c.service.GetStories(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(stories)
}

func (c *Controller) GetStoryByID(ctx *fiber.Ctx) error {
	storyID := ctx.Params("storyID")
	story, err := c.service.GetStoryByID(ctx.Context(), storyID)
	if err != nil {
		return err
	}
	return ctx.JSON(story)
}

func RegisterController(v0 *server.V0, c Controller) {
	v0.Get("groups", c.GetGroups)
	v0.Get("groups/:groupID", c.GetGroupByID)
	v0.Get("stories", c.GetStories)
	v0.Get("stories/:storyID", c.GetStoryByID)
}
