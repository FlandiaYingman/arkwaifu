package artext

import (
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/app/story"
)

type Service struct {
	art   *art.Service
	story *story.Service
}

func newService(serviceArt *art.Service, serviceStory *story.Service) *Service {
	return &Service{
		art:   serviceArt,
		story: serviceStory,
	}
}
