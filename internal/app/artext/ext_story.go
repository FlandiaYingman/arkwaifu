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

func (c *Controller) GetArtsExceptForStoryArts(ctx *fiber.Ctx) error {
	queryExceptForStoryArts := ctx.QueryBool("except-for-story-arts", false)
	if !queryExceptForStoryArts {
		return ctx.Next()
	}

	server, err := ark.ParseServer(ctx.Query("server"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	arts, err := c.service.GetArtsExceptForStoryArts(server)
	if err != nil {
		return err
	}

	return ctx.JSON(arts)
}

func (s *Service) GetArtsOfStoryGroup(server ark.Server, groupID string) ([]*art.Art, error) {
	group, err := s.story.GetStoryGroup(server, groupID)
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

	arts, err := s.art.SelectArtsByIDs(ids)
	if err != nil {
		return nil, err
	}

	return arts, err
}

func (s *Service) GetArtsOfStory(server ark.Server, storyID string) ([]*art.Art, error) {
	story, err := s.story.GetStory(server, storyID)
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

	arts, err := s.art.SelectArtsByIDs(ids)
	if err != nil {
		return nil, err
	}

	return arts, err
}

// GetArtsExceptForStoryArts gets all arts but except for the arts which are present in the story tree.
func (s *Service) GetArtsExceptForStoryArts(server ark.Server) ([]*art.Art, error) {
	pictureArts, err := s.story.GetPictureArts(server)
	if err != nil {
		return nil, err
	}
	characterArts, err := s.story.GetCharacterArts(server)
	if err != nil {
		return nil, err
	}

	storyArtsIDSet := make(map[string]bool)
	for _, pictureArt := range pictureArts {
		storyArtsIDSet[pictureArt.ID] = true
	}
	for _, characterArt := range characterArts {
		storyArtsIDSet[characterArt.ID] = true
	}

	arts, err := s.art.SelectArts()
	if err != nil {
		return nil, err
	}

	result := make([]*art.Art, 0)
	for _, art := range arts {
		if !storyArtsIDSet[art.ID] {
			result = append(result, art)
		}
	}

	return result, nil
}
