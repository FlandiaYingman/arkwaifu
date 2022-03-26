package avg

type Avg []Group
type Group struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Stories []Story `json:"stories"`
}
type Story struct {
	ID      string  `json:"id"`
	Code    string  `json:"code"`
	Name    string  `json:"name"`
	Tag     string  `json:"tag"`
	GroupID string  `json:"groupID"`
	Assets  []Asset `json:"assets"`
}
type Asset struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

func avgToModels(avg Avg) ([]groupModel, []storyModel) {
	groupModels := make([]groupModel, len(avg))
	storyModels := make([]storyModel, 0, len(avg))
	for i, group := range avg {
		groupModels[i] = groupModel{
			ID:      group.ID,
			Name:    group.Name,
			Type:    group.Type,
			Stories: nil,
		}
		storyModels = append(storyModels, storiesToModels(group, group.Stories)...)
	}
	return groupModels, storyModels
}
func storiesToModels(group Group, stories []Story) []storyModel {
	groupModels := make([]storyModel, len(stories))
	for i, story := range stories {
		assets := make([]*assetModel, len(story.Assets))
		for i, asset := range story.Assets {
			assets[i] = &assetModel{
				ID:      asset.ID,
				Kind:    asset.Kind,
				StoryID: story.ID,
			}
		}
		groupModels[i] = storyModel{
			ID:      story.ID,
			Code:    story.Code,
			Name:    story.Name,
			Tag:     story.Tag,
			Assets:  assets,
			GroupID: group.ID,
		}
	}
	return groupModels
}
func groupsFromModels(groupModels []groupModel) []Group {
	groups := make([]Group, len(groupModels))
	for i, model := range groupModels {
		groups[i] = groupFromModel(model)
	}
	return groups
}
func storiesFromModels(storyModels []storyModel) []Story {
	stories := make([]Story, len(storyModels))
	for i, model := range storyModels {
		stories[i] = storyFromModel(model)
	}
	return stories
}
func storiesFromModelsPtr(storyModels []*storyModel) []Story {
	stories := make([]Story, len(storyModels))
	for i, model := range storyModels {
		stories[i] = storyFromModel(*model)
	}
	return stories
}

func groupFromModel(model groupModel) Group {
	return Group{
		ID:      model.ID,
		Name:    model.Name,
		Type:    model.Type,
		Stories: storiesFromModelsPtr(model.Stories),
	}
}
func storyFromModel(model storyModel) Story {
	assets := make([]Asset, len(model.Assets))
	for i, asset := range model.Assets {
		assets[i] = Asset{
			ID:   asset.ID,
			Kind: asset.Kind,
		}
	}
	return Story{
		ID:      model.ID,
		Code:    model.Code,
		Name:    model.Name,
		Tag:     model.Tag,
		Assets:  assets,
		GroupID: model.GroupID,
	}
}
