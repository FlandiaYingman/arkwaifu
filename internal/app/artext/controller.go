package artext

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Controller struct {
	service *Service
}

func newController(service *Service) *Controller {
	return &Controller{service: service}
}

type RegisterParams struct {
	fx.In
	Router fiber.Router `optional:"true"`
}

func (c *Controller) register(params RegisterParams) {
	router := params.Router
	if router == nil {
		return
	}

	// extension of story
	router.Get("arts",
		c.GetArtsOfStoryGroup,
		c.GetArtsOfStory,
		c.GetArtsExceptForStoryArts,
	)

	// extension of character
	router.Get("arts/:id/siblings", c.GetSiblingsOfCharacterArt)
}
