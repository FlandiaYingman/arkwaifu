package arkparser

import (
	"strings"
)

func (p *Parser) ParsePictures(directives []Directive) []*StoryPicture {
	pictureSlice := make([]string, 0)
	pictureMap := make(map[string]StoryPicture)

	addPicture := func(id string, category string) {
		picture := pictureMap[id]
		picture.ID = id
		picture.Category = category
		pictureMap[id] = picture
		pictureSlice = append(pictureSlice, id)
	}

	for _, directive := range directives {
		switch directive.Name {
		case "image":
			if id, ok := directive.Params["image"]; ok {
				addPicture(id, "image")
			}
		case "background":
			if id, ok := directive.Params["image"]; ok {
				addPicture(id, "background")
			}
		case "largebg", "gridbg":
			if ids, ok := directive.Params["imagegroup"]; ok {
				for _, id := range strings.Split(ids, "/") {
					addPicture(id, "background")
				}
			}
		case "showitem":
			if id, ok := directive.Params["image"]; ok {
				addPicture(id, "item")
			}
		}
	}

	result := make([]*StoryPicture, 0)
	for _, id := range pictureSlice {
		picture := pictureMap[id]
		result = append(result, &picture)
	}

	return result
}
