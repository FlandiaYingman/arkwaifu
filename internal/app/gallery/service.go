package gallery

import "github.com/flandiayingman/arkwaifu/internal/pkg/ark"

type Service struct {
	r *repository
}

func newService(r *repository) *Service {
	return &Service{r: r}
}

func (s *Service) ListGalleries(server ark.Server) ([]Gallery, error) {
	return s.r.ListGalleries(server)
}

func (s *Service) GetGalleryByID(server ark.Server, id string) (*Gallery, error) {
	return s.r.GetGalleryByID(server, id)
}

func (s *Service) ListGalleryArts(server ark.Server) ([]Art, error) {
	return s.r.ListGalleryArts(server)
}

func (s *Service) GetGalleryArtByID(server ark.Server, id string) (*Art, error) {
	return s.r.GetGalleryArtByID(server, id)
}

func (s *Service) Put(g []Gallery) (err error) {
	return s.r.Put(g)
}
