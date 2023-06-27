package art

import (
	"errors"
)

type Service struct {
	repo *repository
}

func newService(repo *repository) *Service {
	return &Service{repo}
}

var (
	ErrNotFound = errors.New("the resource with the identifier(s) is not found")
)

func (s *Service) SelectArts(category *string) ([]*Art, error) {
	if category != nil {
		return s.repo.SelectArtsByCategory(*category)
	}
	return s.repo.SelectArts()
}
func (s *Service) SelectArt(id string) (*Art, error) {
	return s.repo.SelectArt(id)
}
func (s *Service) SelectVariants(id string) ([]*Variant, error) {
	return s.repo.SelectVariants(id)
}
func (s *Service) SelectVariant(id string, kind string) (*Variant, error) {
	return s.repo.SelectVariant(id, kind)
}

func (s *Service) UpsertArts(arts ...*Art) error {
	return s.repo.UpsertArts(arts...)
}
func (s *Service) UpsertVariants(variants ...*Variant) error {
	return s.repo.UpsertVariants(variants...)
}

func (s *Service) StoreContent(id string, variation string, content []byte) (err error) {
	return s.repo.StoreContent(id, variation, content)
}
func (s *Service) TakeContent(id string, variation string) (content []byte, err error) {
	return s.repo.TakeContent(id, variation)
}

func (s *Service) SelectArtsWhoseVariantAbsent(variation string) ([]*Art, error) {
	return s.repo.SelectArtsWhereVariantAbsent(variation)
}
