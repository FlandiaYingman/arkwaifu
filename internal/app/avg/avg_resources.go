package avg

import (
	"arkwaifu/internal/app/server"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return Controller{service}
}

type Group struct {
	ID        string
	Name      string
	StoryList []Story
}

type Story struct {
	ID                string
	Code              string
	Name              string
	Tag               string
	ImageResList      []string
	BackgroundResList []string
}

func (c Controller) GetGroups(ctx *fiber.Ctx) error {
	return nil

}

func (c Controller) GetGroupByID(ctx *fiber.Ctx) error {
	return nil

}

func (c Controller) GetAvgs(ctx *fiber.Ctx) error {
	return nil

}

func (c Controller) GetAvgByID(ctx *fiber.Ctx) error {
	return nil
}

func RegisterController(v0 *server.V0, c Controller) {
	v0.Get("groups", c.GetGroups)
	v0.Get("groups/:groupID", c.GetGroupByID)
	v0.Get("avgs", c.GetAvgs)
	v0.Get("avgs/:avgID", c.GetAvgByID)
}
