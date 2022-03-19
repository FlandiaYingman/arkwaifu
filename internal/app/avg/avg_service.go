package avg

import (
	"context"
)

type Service struct {
	versionRepo *VersionRepo
	groupRepo   *GroupRepo
	storyRepo   *StoryRepo
}

func NewService(versionRepo *VersionRepo, groupRepo *GroupRepo, storyRepo *StoryRepo) *Service {
	return &Service{
		versionRepo: versionRepo,
		groupRepo:   groupRepo,
		storyRepo:   storyRepo,
	}
}

func (s *Service) GetVersion(ctx context.Context) (string, error) {
	return s.versionRepo.GetVersion(ctx)
}

func (s *Service) SetAvgs(version string, groups []Group) (err error) {
	ctx := context.Background()
	err = s.versionRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = s.versionRepo.EndTx(err) }()

	err = s.groupRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = s.groupRepo.EndTx(err) }()

	err = s.storyRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = s.storyRepo.EndTx(err) }()

	err = s.versionRepo.UpsertVersion(ctx, version)
	if err != nil {
		return err
	}

	groupModels, storyModels := groupsToModels(groups)

	err = s.groupRepo.Truncate(ctx)
	if err != nil {
		return err
	}
	err = s.groupRepo.InsertGroups(ctx, groupModels)
	if err != nil {
		return err
	}

	err = s.storyRepo.Truncate(ctx)
	if err != nil {
		return err
	}
	err = s.storyRepo.InsertStories(ctx, storyModels)
	if err != nil {
		return err
	}

	return
}

func (s *Service) GetGroups(ctx context.Context) ([]Group, error) {
	groups, err := s.groupRepo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}
	return groupsFromModels(groups), nil
}

func (s *Service) GetGroupByID(ctx context.Context, id string) (*Group, error) {
	model, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	group := groupFromModel(*model)
	return &group, nil
}

func (s *Service) GetStories(ctx context.Context) ([]Story, error) {
	stories, err := s.storyRepo.GetStories(ctx)
	if err != nil {
		return nil, err
	}
	return storiesFromModels(stories), nil
}

func (s *Service) GetStoryByID(ctx context.Context, id string) (*Story, error) {
	model, err := s.storyRepo.GetStoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	story := storyFromModel(*model)
	return &story, nil
}

func (s *Service) Reset(ctx context.Context) error {
	return s.versionRepo.UpsertVersion(ctx, "")
}
