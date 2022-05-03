package avg

import (
	"time"

	"github.com/flandiayingman/arkwaifu/internal/app/server"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service}
}
func RegisterController(v0 *server.V0, c *Controller) {
	router := v0.
		Group("/avg").
		Use(newCompress()).
		Use(newCache())
	router.Get("groups", c.GetGroups)
	router.Get("groups/:groupID", c.GetGroupByID)
	router.Get("stories", c.GetStories)
	router.Get("stories/:storyID", c.GetStoryByID)
}
func newCompress() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
func newCache() fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   1 * time.Minute,
		CacheControl: true,
	})
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
