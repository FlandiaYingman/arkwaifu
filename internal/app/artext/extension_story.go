package artext

import (
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetArtsOfStoryGroup(ctx *fiber.Ctx) error {
	queryGroup := ctx.Query("group")
	if queryGroup == "" {
		return ctx.Next()
	}

	server, err := ark.ParseServer(ctx.Query("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	arts, err := c.service.GetArtsOfStoryGroup(server, queryGroup)
	if err != nil {
		return err
	}

	return ctx.JSON(arts)
}

func (c *Controller) GetArtsOfStory(ctx *fiber.Ctx) error {
	queryStory := ctx.Query("story")
	if queryStory == "" {
		return ctx.Next()
	}

	server, err := ark.ParseServer(ctx.Query("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	arts, err := c.service.GetArtsOfStory(server, queryStory)
	if err != nil {
		return err
	}

	return ctx.JSON(arts)
}

func (service *Service) GetArtsOfStoryGroup(server ark.Server, groupID string) ([]*art.Art, error) {
	group, err := service.story.GetStoryGroup(server, groupID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, story := range group.Stories {
		for _, art := range story.PictureArts {
			ids = append(ids, art.ID)
		}
		for _, art := range story.CharacterArts {
			ids = append(ids, art.ID)
		}
	}

	arts, err := service.art.SelectArtsByIDs(ids)
	if err != nil {
		return nil, err
	}

	return arts, err
}

func (service *Service) GetArtsOfStory(server ark.Server, storyID string) ([]*art.Art, error) {
	story, err := service.story.GetStory(server, storyID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, art := range story.PictureArts {
		ids = append(ids, art.ID)
	}
	for _, art := range story.CharacterArts {
		ids = append(ids, art.ID)
	}

	arts, err := service.art.SelectArtsByIDs(ids)
	if err != nil {
		return nil, err
	}

	return arts, err
}
