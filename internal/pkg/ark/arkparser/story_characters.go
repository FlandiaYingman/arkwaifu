package arkparser

import (
	"fmt"
	"github.com/wk8/go-ordered-map/v2"
	"regexp"
)

func (p *Parser) ParseCharacters(directives []Directive) []*StoryCharacter {
	stage := characterStage{
		spotlight:  "",
		characters: make(map[string]string),
		history:    make([]string, 0),
		names:      make(map[string][]string),
	}
	for _, directive := range directives {
		switch directive.Name {
		case "":
			name := directive.Params["name"]
			stage.name(name)
		case "character":
			id1 := NormalizeCharacterID(directive.Params["name"])
			id2 := NormalizeCharacterID(directive.Params["name2"])
			stage.take("1", id1)
			stage.take("2", id2)
			stage.focus(directive.Params["focus"], "1")
		case "charslot":
			if id := NormalizeCharacterID(directive.Params["name"]); id != "" {
				stage.take(directive.Params["slot"], id)
				stage.focus(directive.Params["focus"], directive.Params["slot"])
			} else {
				stage.exit()
			}
		case "dialog":
			stage.exit()
		}
	}

	result := make([]*StoryCharacter, 0)
	for _, id := range stage.history {
		result = append(result, &StoryCharacter{
			ID:    id,
			Names: unique(stage.names[id]),
		})
	}

	return result
}

func NormalizeCharacterID(id string) string {
	if id == "" {
		return ""
	}

	regex := regexp.MustCompile(`^(.*?)(?:#(\d+))?(?:\$(\d+))?$`)
	matches := regex.FindStringSubmatch(id)
	if len(matches) == 0 {
		return ""
	}

	id = matches[1]
	faceNum, bodyNum := "1", "1" // default values are 1
	if matches[2] != "" {
		faceNum = matches[2]
	}
	if matches[3] != "" {
		bodyNum = matches[3]
	}

	return fmt.Sprintf("%s#%s$%s", id, faceNum, bodyNum)
}

type characterStage struct {
	spotlight  string
	characters map[string]string

	history []string
	names   map[string][]string
}

func (s *characterStage) protagonist() string {
	return s.characters[s.spotlight]
}
func (s *characterStage) focus(slot string, defaultSlot string) {
	if slot != "" {
		s.spotlight = slot
	} else if len(s.characters) == 1 {
		for slot, _ := range s.characters {
			s.spotlight = slot
		}
	} else {
		s.spotlight = ""
	}
}
func (s *characterStage) take(slot string, id string) {
	if id != "" {
		s.characters[slot] = id
		s.history = append(s.history, id)
	} else {
		delete(s.characters, slot)
	}
}
func (s *characterStage) name(name string) {
	s.names[s.protagonist()] = append(s.names[s.protagonist()], name)
}
func (s *characterStage) exit() {
	s.spotlight = ""
	s.characters = make(map[string]string)
}

func unique[T comparable](slice []T) []T {
	m := orderedmap.New[T, struct{}](len(slice))
	for _, el := range slice {
		m.Set(el, struct{}{})
	}
	slice = nil
	for pair := m.Oldest(); pair != nil; pair = pair.Next() {
		slice = append(slice, pair.Key)
	}
	return slice
}
