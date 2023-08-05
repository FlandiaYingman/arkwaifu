package story

import (
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	arkparser2 "github.com/flandiayingman/arkwaifu/internal/pkg/arkparser"
	cols "github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"strings"
)

type Service struct {
	r *repo
}

func newService(r *repo) *Service {
	return &Service{r: r}
}

type GroupFilter struct {
	Type string
}

func (s *Service) GetStories(server ark.Server) ([]*Story, error) {
	return s.r.SelectStories(server)
}
func (s *Service) GetStory(server ark.Server, id string) (*Story, error) {
	return s.r.SelectStory(id, server)
}
func (s *Service) GetStoryGroups(server ark.Server, filter GroupFilter) ([]*Group, error) {
	if filter.Type != "" {
		return s.r.SelectStoryGroupsByType(server, filter.Type)
	}
	return s.r.SelectStoryGroups(server)
}
func (s *Service) GetStoryGroup(server ark.Server, id string) (*Group, error) {
	return s.r.SelectStoryGroup(id, server)
}

func (s *Service) GetPictureArts(server ark.Server) ([]*PictureArt, error) {
	return s.r.SelectPictureArts(server)
}
func (s *Service) GetAggregatedPictureArt(server ark.Server, id string) (*AggregatedPictureArt, error) {
	return s.r.SelectAggregatedPictureArtByID(server, id)
}

func (s *Service) GetCharacterArts(server ark.Server) ([]*CharacterArt, error) {
	return s.r.SelectCharacterArts(server)
}
func (s *Service) GetAggregatedCharacterArt(server ark.Server, id string) (*AggregatedCharacterArt, error) {
	return s.r.SelectAggregatedCharacterArtByID(server, id)
}

func (s *Service) PopulateFrom(rawTree *arkparser2.StoryTree, server ark.Server) error {
	converter := objectConverter{server: server}

	tree, err := converter.convertStoryTree(rawTree)
	if err != nil {
		return fmt.Errorf("converter error when populating story service: %w", err)
	}

	for _, group := range tree {
		// Prevent error: extended protocol limited to 65535 parameters
		// TODO: Find a more efficient way to prevent error.
		err = s.r.UpsertStoryGroups([]Group{group})
		if err != nil {
			return err
		}
	}

	return nil
}

type objectConverter struct {
	server ark.Server
}

func (c *objectConverter) convertStoryTree(tree *arkparser2.StoryTree) (Tree, error) {
	return cols.MapErr(tree.StoryGroups, c.convertStoryGroup)
}
func (c *objectConverter) convertStoryGroup(rawGroup *arkparser2.StoryGroup) (Group, error) {
	group := Group{
		Server:  c.server,
		ID:      strings.ToLower(rawGroup.ID),
		Name:    rawGroup.Name,
		Type:    "",
		Stories: nil,
		SortID:  nil, // SortID is auto-increment, no need to handle it.
	}

	var err error

	group.Type, err = c.convertGroupType(rawGroup.Type)
	if err != nil {
		return Group{}, err
	}

	group.Stories, err = cols.MapErr(rawGroup.Stories, c.convertStory)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}
func (c *objectConverter) convertGroupType(rawType string) (GroupType, error) {
	switch rawType {
	case arkparser2.TypeMain:
		return GroupTypeMainStory, nil
	case arkparser2.TypeActivity:
		return GroupTypeMajorEvent, nil
	case arkparser2.TypeMiniActivity:
		return GroupTypeMinorEvent, nil
	case arkparser2.TypeNone:
		return GroupTypeOther, nil
	default:
		return "", fmt.Errorf("unknown story group type: %v", rawType)
	}
}
func (c *objectConverter) convertStory(rawStory *arkparser2.Story) (Story, error) {
	story := Story{
		Server:        c.server,
		ID:            strings.ToLower(rawStory.ID),
		Tag:           "",
		TagText:       rawStory.TagText,
		Code:          rawStory.Code,
		Name:          rawStory.Name,
		Info:          rawStory.Info,
		GroupID:       rawStory.GroupID,
		SortID:        nil, // SortID is auto-increment, no need to handle it.
		PictureArts:   cols.Map(rawStory.Pictures, c.convertPictureArt),
		CharacterArts: cols.Map(rawStory.Characters, c.convertCharacterArt),
	}

	var err error

	story.Tag, err = c.convertStoryTagType(rawStory.TagType)
	if err != nil {
		return Story{}, err
	}

	return story, nil
}
func (c *objectConverter) convertStoryTagType(rawType string) (Tag, error) {
	switch rawType {
	case arkparser2.TagBefore:
		return TagBefore, nil
	case arkparser2.TagAfter:
		return TagAfter, nil
	case arkparser2.TagInterlude:
		return TagInterlude, nil
	default:
		return "", fmt.Errorf("unknown arkparser story group type: %v", rawType)
	}
}
func (c *objectConverter) convertPictureArt(rawPicture *arkparser2.StoryPicture) PictureArt {
	return PictureArt{
		Server:   c.server,
		ID:       strings.ToLower(rawPicture.ID),
		StoryID:  "", // This will be auto-generated by the ORM framework
		Category: rawPicture.Category,

		Title:    rawPicture.Title,
		Subtitle: rawPicture.Subtitle,

		SortID: nil, // SortID is auto-increment, no need to handle it.
	}
}
func (c *objectConverter) convertCharacterArt(rawCharacter *arkparser2.StoryCharacter) CharacterArt {
	return CharacterArt{
		Server:   c.server,
		ID:       strings.ToLower(rawCharacter.ID),
		StoryID:  "", // This will be auto-generated by the ORM framework
		Category: "character",

		Names: rawCharacter.Names,

		SortID: nil, // SortID is auto-increment, no need to handle it.
	}
}
