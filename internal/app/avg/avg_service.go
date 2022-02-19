package avg

import (
	"context"
	"github.com/uptrace/bun"
)

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

func groupsToModels(groups []Group) ([]GroupModel, []StoryModel) {
	groupModels := make([]GroupModel, len(groups))
	storyModels := make([]StoryModel, 0, len(groups))
	for i, group := range groups {
		groupModels[i] = GroupModel{
			ID:      group.ID,
			Name:    group.Name,
			Stories: nil,
		}
		storyModels = append(storyModels, storiesToModels(group, group.StoryList)...)
	}
	return groupModels, storyModels
}

func storiesToModels(group Group, stories []Story) []StoryModel {
	groupModels := make([]StoryModel, len(stories))
	for i, story := range stories {
		images := make([]*ImageModel, len(story.ImageResList))
		for i, image := range story.ImageResList {
			images[i] = &ImageModel{
				StoryID: story.ID,
				Image:   image,
			}
		}
		backgrounds := make([]*BackgroundModel, len(story.BackgroundResList))
		for i, background := range story.BackgroundResList {
			backgrounds[i] = &BackgroundModel{
				StoryID:    story.ID,
				Background: background,
			}
		}
		groupModels[i] = StoryModel{
			BaseModel:   bun.BaseModel{},
			ID:          story.ID,
			Code:        story.Code,
			Name:        story.Name,
			Tag:         story.Tag,
			Images:      images,
			Backgrounds: backgrounds,
			GroupID:     group.ID,
		}
	}
	return groupModels
}

func groupsFromModels(groupModels []GroupModel) []Group {
	groups := make([]Group, len(groupModels))
	for i, model := range groupModels {
		groups[i] = Group{
			ID:        model.ID,
			Name:      model.Name,
			StoryList: storiesFromModelsPtr(model.Stories),
		}
	}
	return groups
}

func storiesFromModels(storyModels []StoryModel) []Story {
	stories := make([]Story, len(storyModels))
	for i, model := range storyModels {
		stories[i] = storyFromModel(model)
	}
	return stories
}

func storiesFromModelsPtr(storyModels []*StoryModel) []Story {
	stories := make([]Story, len(storyModels))
	for i, model := range storyModels {
		stories[i] = storyFromModel(*model)
	}
	return stories
}

func groupFromModel(model GroupModel) Group {
	return Group{
		ID:        model.ID,
		Name:      model.Name,
		StoryList: storiesFromModelsPtr(model.Stories),
	}
}

func storyFromModel(model StoryModel) Story {
	images := make([]string, len(model.Images))
	for i, image := range model.Images {
		images[i] = image.Image
	}
	backgrounds := make([]string, len(model.Backgrounds))
	for i, background := range model.Backgrounds {
		backgrounds[i] = background.Background
	}
	return Story{
		ID:                model.ID,
		Code:              model.Code,
		Name:              model.Name,
		Tag:               model.Tag,
		ImageResList:      images,
		BackgroundResList: backgrounds,
	}
}

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
