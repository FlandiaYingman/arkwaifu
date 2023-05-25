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
	for i := range arts {
		for j := range arts[i].Variants {
			if !arts[i].Variants[j].ContentPresent {
				arts[i].Variants[j].ContentPath = arts[i].Variants[j].ToStatic().PathRel()
				arts[i].Variants[j].ContentWidth = nil
				arts[i].Variants[j].ContentHeight = nil
			}
		}
	}
	return s.repo.UpsertArts(arts...)
}
func (s *Service) UpsertVariants(variants ...*Variant) error {
	for i := range variants {
		if !variants[i].ContentPresent {
			variants[i].ContentPath = variants[i].ToStatic().PathRel()
			variants[i].ContentWidth = nil
			variants[i].ContentHeight = nil
		}
	}
	return s.repo.UpsertVariants(variants...)
}

func (s *Service) StoreStatics(statics ...*VariantContent) error {
	return s.repo.StoreStatics(statics...)
}
func (s *Service) TakeStatics(statics ...*VariantContent) error {
	return s.repo.TakeStatics(statics...)
}

func (s *Service) SelectArtsWhoseVariantAbsent(variation string) ([]*Art, error) {
	return s.repo.SelectArtsWhereVariantAbsent(variation)
}
