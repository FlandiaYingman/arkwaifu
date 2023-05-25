package story

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
	router.Get(":server/story-groups", c.GetGroups)
	router.Get(":server/story-groups/:id", c.GetGroup)
	router.Get(":server/stories", c.GetStories)
	router.Get(":server/stories/:id", c.GetStory)

	router.Get(":server/arts", c.GetArts)
	router.Get(":server/arts/:id", c.GetArts) // TODO
}

func (c *controller) GetGroups(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	queryType := ctx.Query("type")

	queryFilter := GroupFilter{
		Type: queryType,
	}

	groups, err := c.service.GetStoryGroups(server, queryFilter)
	if err != nil {
		return err
	}
	return ctx.JSON(groups)
}
func (c *controller) GetGroup(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	group, err := c.service.GetStoryGroup(server, id)
	if err != nil {
		return err
	}

	return ctx.JSON(group)
}

func (c *controller) GetStories(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	stories, err := c.service.GetStories(server)
	if err != nil {
		return err
	}

	return ctx.JSON(stories)
}
func (c *controller) GetStory(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	story, err := c.service.GetStory(server, id)
	if err != nil {
		return err
	}

	return ctx.JSON(story)
}

func (c *controller) GetArts(ctx *fiber.Ctx) error {
	server, err := ark.ParseServer(ctx.Params("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	pictureArts, err := c.service.GetPictureArts(server)
	if err != nil {
		return err
	}
	characterArts, err := c.service.GetCharacterArts(server)
	if err != nil {
		return err
	}

	arts := make([]any, 0)
	for _, art := range pictureArts {
		arts = append(arts, art)
	}
	for _, art := range characterArts {
		arts = append(arts, art)
	}

	return ctx.JSON(arts)
}
