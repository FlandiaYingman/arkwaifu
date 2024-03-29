package art

import (
	"errors"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/google/uuid"
	"strings"
)

type Service struct {
	repo *repository

	users []infra.User
}

func newService(config *infra.Config, repo *repository) *Service {
	return &Service{
		repo:  repo,
		users: config.Users,
	}
}

var (
	ErrNotFound = errors.New("the resource with the identifier(s) is not found")
)

func (s *Service) SelectArts() ([]*Art, error) {
	return s.repo.SelectArts()
}
func (s *Service) SelectArtsByCategory(category Category) ([]*Art, error) {
	return s.repo.SelectArtsByCategory(string(category))
}
func (s *Service) SelectArtsByIDs(ids []string) ([]*Art, error) {
	for i, id := range ids {
		ids[i] = strings.ToLower(id)
	}
	return s.repo.SelectArtsByIDs(ids)
}
func (s *Service) SelectArtsByIDLike(like string) ([]*Art, error) {
	return s.repo.SelectArtsByIDLike(like)
}
func (s *Service) SelectArt(id string) (*Art, error) {
	return s.repo.SelectArt(strings.ToLower(id))
}
func (s *Service) SelectVariants(id string) ([]*Variant, error) {
	return s.repo.SelectVariants(strings.ToLower(id))
}
func (s *Service) SelectVariant(id string, variation Variation) (*Variant, error) {
	return s.repo.SelectVariant(strings.ToLower(id), string(variation))
}

func (s *Service) UpsertArts(arts ...*Art) error {
	for _, art := range arts {
		art.ID = strings.ToLower(art.ID)
	}
	return s.repo.UpsertArts(arts...)
}
func (s *Service) UpsertVariants(variants ...*Variant) error {
	for _, variant := range variants {
		variant.ArtID = strings.ToLower(variant.ArtID)
	}
	return s.repo.UpsertVariants(variants...)
}

func (s *Service) StoreContent(id string, variation Variation, content []byte) (err error) {
	return s.repo.StoreContent(strings.ToLower(id), string(variation), content)
}
func (s *Service) TakeContent(id string, variation Variation) (content []byte, err error) {
	return s.repo.TakeContent(strings.ToLower(id), string(variation))
}

func (s *Service) SelectArtsWhoseVariantAbsent(variation Variation) ([]*Art, error) {
	return s.repo.SelectArtsWhereVariantAbsent(string(variation))
}

func (s *Service) Authenticate(uuid uuid.UUID) *infra.User {
	for _, user := range s.users {
		if user.ID == uuid {
			return &user
		}
	}
	return nil
}
