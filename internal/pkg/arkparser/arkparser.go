package arkparser

import "github.com/pkg/errors"

type Parser struct {
	Root   string
	Prefix string

	pictureMap   map[string]*StoryPicture
	characterMap map[string]*StoryCharacter
}

type StoryTree struct {
	StoryGroups []*StoryGroup
}

type StoryGroup struct {
	ID      string
	Name    string
	Type    string
	Stories []*Story
}

type Story struct {
	ID      string
	GroupID string
	Code    string
	Name    string
	Info    string
	TagType string
	TagText string

	Pictures   []*StoryPicture
	Characters []*StoryCharacter
}

type StoryCharacter struct {
	ID    string
	Names []string
}

type StoryPicture struct {
	ID       string
	Category string
	Title    string
	Subtitle string
}

func (p *Parser) getPicture(id string) *StoryPicture {
	picture, ok := p.pictureMap[id]
	if !ok {
		picture = &StoryPicture{}
		p.pictureMap[id] = picture
	}
	return picture
}
func (p *Parser) setPicture(id string, picture *StoryPicture) {
	p.pictureMap[id] = picture
}
func (p *Parser) getCharacter(id string) *StoryCharacter {
	character, ok := p.characterMap[id]
	if !ok {
		character = &StoryCharacter{}
		p.characterMap[id] = character
	}
	return character
}
func (p *Parser) setCharacter(id string, character *StoryCharacter) {
	p.characterMap[id] = character
}

func (p *Parser) Parse() (*StoryTree, error) {
	jsonStoryReviewTable, err := p.ParseStoryReviewTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tree := StoryTree{}
	picIndex := make(map[string][]*StoryPicture)
	charIndex := make(map[string][]*StoryCharacter)
	for _, jsonGroup := range jsonStoryReviewTable {
		group := StoryGroup{
			ID:      jsonGroup.ID,
			Name:    jsonGroup.Name,
			Type:    jsonGroup.Type,
			Stories: nil,
		}
		for _, jsonStory := range jsonGroup.Stories {
			info, _ := p.GetInfo(jsonStory.Info)
			// If info is not found, it is likely that there is no info for the story.
			// Example: main_14_level_main_14-20_beg
			//if err != nil {
			//	return nil, errors.WithStack(err)
			//}
			directives, err := p.ParseStoryText(jsonStory.Text)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			pictures := p.ParsePictures(directives)
			characters := p.ParseCharacters(directives)
			story := Story{
				ID:         jsonStory.ID,
				GroupID:    jsonStory.ID,
				Code:       jsonStory.Code,
				Name:       jsonStory.Name,
				Info:       info,
				TagType:    jsonStory.Tag.Type,
				TagText:    jsonStory.Tag.Text,
				Pictures:   pictures,
				Characters: characters,
			}

			for _, picture := range pictures {
				picIndex[picture.ID] = append(picIndex[picture.ID], picture)
			}
			for _, character := range characters {
				charIndex[character.ID] = append(charIndex[character.ID], character)
			}

			group.Stories = append(group.Stories, &story)
		}
		tree.StoryGroups = append(tree.StoryGroups, &group)
	}

	storyReviewMetaTable, err := p.ParseStoryReviewMetaTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for pair := storyReviewMetaTable.ActArchiveResData.Pics.Oldest(); pair != nil; pair = pair.Next() {
		pic := pair.Value
		for _, picture := range picIndex[pic.AssetPath] {
			picture.Title = pic.Desc
			picture.Subtitle = pic.PicDescription
		}
	}

	return &tree, nil
}
