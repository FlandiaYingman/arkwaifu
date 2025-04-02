package arkscanner

import (
	"encoding/json"
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"image"
	"math"
	"os"
	"path"
	"path/filepath"
)

type (
	CharacterArt struct {
		ID             string
		Kind           string
		BodyVariations []CharacterArtBodyVariation
	}
	CharacterArtBodyVariation struct {
		BodySprite      string
		BodySpriteAlpha string
		FaceRectangle   image.Rectangle
		FaceVariations  []CharacterArtFaceVariation

		// PixelToUnits is the ratio of pixel to unit.
		// In other words, the resulting image is scaled by *dividing* this value.
		// This is because the original image is in units, and we convert it to pixels.
		PixelToUnits float64
	}
	CharacterArtFaceVariation struct {
		FaceSprite      string
		FaceSpriteAlpha string
		WholeBody       bool

		// PixelToUnits is the ratio of pixel to unit.
		// In other words, the resulting image is scaled by *dividing* this value.
		// This is because the original image is in units, and we convert it to pixels.
		PixelToUnits float64
	}
)

func (a *CharacterArt) FacePath(bodyNum int, faceNum int) string {
	// note: "*num" starts with 1, not 0
	sprite := a.BodyVariations[bodyNum-1].FaceVariations[faceNum-1].FaceSprite
	if sprite == "" {
		return ""
	}
	return path.Join(characterPrefix, a.ID, sprite)
}
func (a *CharacterArt) FacePathAlpha(bodyNum int, faceNum int) string {
	// note: "*num" starts with 1, not 0
	sprite := a.BodyVariations[bodyNum-1].FaceVariations[faceNum-1].FaceSpriteAlpha
	if sprite == "" {
		return ""
	}
	return path.Join(characterPrefix, a.ID, sprite)
}
func (a *CharacterArt) BodyPath(bodyNum int) string {
	// note: "*num" starts with 1, not 0
	sprite := a.BodyVariations[bodyNum-1].BodySprite
	if sprite == "" {
		return ""
	}
	return path.Join(characterPrefix, a.ID, sprite)
}
func (a *CharacterArt) BodyPathAlpha(bodyNum int) string {
	// note: "*num" starts with 1, not 0
	sprite := a.BodyVariations[bodyNum-1].BodySpriteAlpha
	if sprite == "" {
		return ""
	}
	return path.Join(characterPrefix, a.ID, sprite)
}

var (
	characterPrefix = "assets/torappu/dynamicassets/avg/characters"
)

func (scanner *Scanner) ScanForCharacterArts() ([]*CharacterArt, error) {
	rootCharacterArts := filepath.Join(scanner.Root, characterPrefix)
	characterEntries, err := os.ReadDir(rootCharacterArts)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	characterEntries = cols.Filter(characterEntries, func(element os.DirEntry) bool { return element.IsDir() })
	characterIDs := cols.Map(characterEntries, func(i os.DirEntry) (o string) { return pathutil.Stem(i.Name()) })
	characterArts, err := cols.MapErr(characterIDs, scanner.ScanCharacter)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return characterArts, nil
}
func (scanner *Scanner) ScanCharacter(id string) (*CharacterArt, error) {
	hubGroupArt, err := scanner.scanHubGroupOfCharacter(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	hubArt, err := scanner.scanHubOfCharacter(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if hubGroupArt == nil && hubArt == nil {
		return nil, errors.Errorf("scan character %s: neither the hub group nor the hub exist", id)
	}
	if hubGroupArt != nil && hubArt != nil {
		return nil, errors.Errorf("scan character %s: both the hub group and the hub exist", id)
	}

	if hubGroupArt != nil {
		return hubGroupArt, nil
	} else {
		return hubArt, nil
	}
}

func (scanner *Scanner) scanHubGroupOfCharacter(id string) (*CharacterArt, error) {
	hubPath := filepath.Join(scanner.Root, characterPrefix, id, "AVGCharacterSpriteHubGroup.json")
	hubJson, err := os.ReadFile(hubPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hubGroup CharacterSpriteHubGroup
	err = json.Unmarshal(hubJson, &hubGroup)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pathIDMap, err := scanner.scanPathIDMapOfCharacter(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ttMap, err := scanner.scanTypeTreeMapOfCharacter(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	art := hubGroup.toArt(id, pathIDMap, ttMap)
	return &art, nil
}
func (scanner *Scanner) scanHubOfCharacter(id string) (*CharacterArt, error) {
	hubPath := filepath.Join(scanner.Root, characterPrefix, id, "AVGCharacterSpriteHub.json")
	hubJson, err := os.ReadFile(hubPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hub CharacterSpriteHub
	err = json.Unmarshal(hubJson, &hub)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pathIDMap, err := scanner.scanPathIDMapOfCharacter(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ttMap, err := scanner.scanTypeTreeMapOfCharacter(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hubGroup := CharacterSpriteHubGroup{SpriteGroups: []CharacterSpriteHub{hub}}
	art := hubGroup.toArt(id, pathIDMap, ttMap)
	return &art, nil
}
func (scanner *Scanner) scanPathIDMapOfCharacter(id string) (map[int64]string, error) {
	mapPath := filepath.Join(scanner.Root, characterPrefix, fmt.Sprintf("%s.json", id))
	mapJson, err := os.ReadFile(mapPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pathIDMap := make(map[int64]string)
	err = json.Unmarshal(mapJson, &pathIDMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pathIDMap, nil
}
func (scanner *Scanner) scanTypeTreeMapOfCharacter(id string) (map[int64]SpriteTypeTree, error) {
	mapPath := filepath.Join(scanner.Root, characterPrefix, fmt.Sprintf("%s.typetree.json", id))
	mapJson, err := os.ReadFile(mapPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ttPathIDMap := make(map[int64]string)
	err = json.Unmarshal(mapJson, &ttPathIDMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ttMap := make(map[int64]SpriteTypeTree)
	for k, v := range ttPathIDMap {
		subMapPath := filepath.Join(scanner.Root, characterPrefix, id, v)
		subMapJson, err := os.ReadFile(subMapPath)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var tt SpriteTypeTree
		err = json.Unmarshal(subMapJson, &tt)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ttMap[k] = tt
	}

	return ttMap, nil
}

type (
	CharacterSpriteHubGroup struct {
		SpriteGroups []CharacterSpriteHub `json:"spriteGroups"`
	}
	CharacterSpriteHub struct {
		Sprites []CharacterSprite `json:"sprites"`
		FacePos struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"FacePos"`
		FaceSize struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"FaceSize"`
	}
	CharacterSprite struct {
		Sprite struct {
			MPathID int64 `json:"m_PathID"`
		} `json:"sprite"`
		AlphaTex struct {
			MPathID int64 `json:"m_PathID"`
		} `json:"alphaTex"`
		IsWholeBody int `json:"isWholeBody"`
	}
	SpriteTypeTree struct {
		PixelsToUnits float64 `json:"m_PixelsToUnits"`
	}
)

func (c *CharacterSpriteHubGroup) toArt(id string, pathIDMap map[int64]string, ttMap map[int64]SpriteTypeTree) (a CharacterArt) {
	convertSpriteHubToArt := func(i CharacterSpriteHub) (o CharacterArtBodyVariation) { return i.toArt(pathIDMap, ttMap) }
	a = CharacterArt{
		ID:             id,
		Kind:           "character",
		BodyVariations: cols.Map(c.SpriteGroups, convertSpriteHubToArt),
	}
	return
}
func (c *CharacterSpriteHub) toArt(pathIDMap map[int64]string, ttMap map[int64]SpriteTypeTree) (a CharacterArtBodyVariation) {
	convertSpriteToArt := func(i CharacterSprite) (o CharacterArtFaceVariation) { return i.toArt(pathIDMap, ttMap) }
	a = CharacterArtBodyVariation{
		BodySprite:      "",
		BodySpriteAlpha: "",
		FaceRectangle: image.Rect(
			int(math.Round(c.FacePos.X)),
			int(math.Round(c.FacePos.Y)),
			int(math.Round(c.FacePos.X+c.FaceSize.X)),
			int(math.Round(c.FacePos.Y+c.FaceSize.Y)),
		),
		FaceVariations: cols.Map(c.Sprites, convertSpriteToArt),
		PixelToUnits:   1.0,
	}
	// If the face pos is valid, then extract the last face variation as the body.
	// Otherwise, all variations contain the whole body.
	if c.FacePos.X >= 0 && c.FacePos.Y >= 0 {
		lastVariation := a.FaceVariations[len(a.FaceVariations)-1]
		a.BodySprite = lastVariation.FaceSprite
		a.BodySpriteAlpha = lastVariation.FaceSpriteAlpha
		a.PixelToUnits = lastVariation.PixelToUnits
		a.FaceVariations = a.FaceVariations[:len(a.FaceVariations)-1]
	} else {
		for i := range a.FaceVariations {
			a.FaceVariations[i].WholeBody = true
		}
	}
	return
}
func (c *CharacterSprite) toArt(pathIDMap map[int64]string, ttMap map[int64]SpriteTypeTree) (a CharacterArtFaceVariation) {
	a = CharacterArtFaceVariation{
		FaceSprite:      pathIDMap[c.Sprite.MPathID],
		FaceSpriteAlpha: pathIDMap[c.AlphaTex.MPathID],
		WholeBody:       c.IsWholeBody != 0,
		PixelToUnits:    ttMap[c.Sprite.MPathID].PixelsToUnits / 100.0,
	}
	return
}
