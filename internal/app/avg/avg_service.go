package avg

import (
	"context"
)

type Service struct {
	r *Repo
}

func NewService(r *Repo) *Service {
	return &Service{r: r}
}

func (s *Service) GetVersion(ctx context.Context) (string, error) {
	return s.r.GetVersion(ctx)
}
func (s *Service) UpdateAvg(ctx context.Context, version string, avg Avg) (err error) {
	gms, sms := avgToModels(avg)
	return s.r.UpdateAvg(ctx, version, gms, sms)
}
func (s *Service) GetGroups(ctx context.Context) ([]Group, error) {
	groups, err := s.r.GetGroups(ctx)
	if err != nil {
		return nil, err
	}
	return groupsFromModels(groups), nil
}
func (s *Service) GetGroupByID(ctx context.Context, id string) (*Group, error) {
	model, err := s.r.GetGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	group := groupFromModel(*model)
	return &group, nil
}
func (s *Service) GetStories(ctx context.Context) ([]Story, error) {
	stories, err := s.r.GetStories(ctx)
	if err != nil {
		return nil, err
	}
	return storiesFromModels(stories), nil
}
func (s *Service) GetStoryByID(ctx context.Context, id string) (*Story, error) {
	model, err := s.r.GetStoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	story := storyFromModel(*model)
	return &story, nil
}
func (s *Service) Reset(ctx context.Context) error {
	return s.r.UpsertVersion(ctx, "")
}
