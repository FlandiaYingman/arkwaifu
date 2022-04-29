package arkavg

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strconv"

	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
)

type CharAsset struct {
	Name string
	Kind Kind
	Hubs []CharAssetSpriteHub
}
type CharAssetSpriteHub struct {
	Sprites            []Asset
	SpritesAlpha       []Asset
	SpritesAlias       []string
	SpritesIsWholeBody []bool
	FacePos            struct{ X, Y float64 }
	FaceSize           struct{ X, Y float64 }
}

func (c *CharAsset) AssetName(hubNum, faceNum int) string {
	// fn and hn should begin with 1, not 0.
	return fmt.Sprintf("%s#%d$%d", c.Name, faceNum+1, hubNum+1)
}
func (c *CharAsset) AssetNameAlpha(hubNum, faceNum int) string {
	// fn and hn should begin with 1, not 0.
	return fmt.Sprintf("%s#%d$%d[alpha]", c.Name, faceNum+1, hubNum+1)
}

func ScanForCharAssets(resDir string, prefix string) ([]CharAsset, error) {
	assetDir := filepath.Join(resDir, prefix, "avg/characters")
	charAssetDirs, err := fileutil.ListDirs(assetDir)
	if err != nil {
		return nil, err
	}

	var charAssets []CharAsset
	for _, charAssetDir := range charAssetDirs {
		charAsset, err := ScanForCharAssetsByChar(filepath.Base(charAssetDir), resDir, prefix)
		if err != nil {
			return nil, err
		}
		charAssets = append(charAssets, charAsset)
	}

	return charAssets, nil
}

func ScanForCharAssetsByChar(name string, resDir, prefix string) (CharAsset, error) {
	dirPath := filepath.Join(resDir, prefix, "avg", "characters", name)
	pathIDMapPath := filepath.Join(filepath.Dir(dirPath), filepath.Base(dirPath)+".json")
	spriteHubPath := filepath.Join(dirPath, "AVGCharacterSpriteHub.json")
	spriteHubGroupPath := filepath.Join(dirPath, "AVGCharacterSpriteHubGroup.json")

	idMap, err := readPathIDMap(pathIDMapPath)
	if err != nil {
		return CharAsset{}, err
	}
	spriteHubs, err := readAllSpriteHubs(spriteHubPath, spriteHubGroupPath)
	if err != nil {
		return CharAsset{}, err
	}

	hubs := make([]CharAssetSpriteHub, 0, len(spriteHubs))
	for _, sh := range spriteHubs {
		cash := parseSpriteHub(name, sh, idMap)
		if err != nil {
			return CharAsset{}, err
		}
		hubs = append(hubs, cash)
	}
	return CharAsset{
		Name: name,
		Kind: KindCharacter,
		Hubs: hubs,
	}, nil
}

type pathIDMap map[int64]string
type spriteHub struct {
	Sprites []struct {
		Sprite struct {
			MPathID int64 `json:"m_PathID"`
		} `json:"sprite"`
		AlphaTex struct {
			MPathID int64 `json:"m_PathID"`
		} `json:"alphaTex"`
		Alias       string `json:"alias"`
		IsWholeBody int    `json:"isWholeBody"`
	} `json:"sprites"`
	FacePos struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"FacePos"`
	FaceSize struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"FaceSize"`
}

func readPathIDMap(pathIDMapPath string) (pathIDMap, error) {
	pathIDMapJson, err := os.ReadFile(pathIDMapPath)
	if err != nil {
		return nil, err
	}
	pathIDStrMap := make(map[string]string)
	pathIDIntMap := make(map[int64]string)
	err = json.Unmarshal(pathIDMapJson, &pathIDStrMap)
	if err != nil {
		return nil, err
	}
	for pathIDStr, name := range pathIDStrMap {
		pathIDInt, err := strconv.ParseInt(pathIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		pathIDIntMap[pathIDInt] = name
	}
	return pathIDIntMap, nil
}

func readAllSpriteHubs(spriteHubPath string, spriteHubGroupPath string) ([]spriteHub, error) {
	spriteHubs := make([]spriteHub, 0)
	exists, err := fileutil.Exists(spriteHubPath)
	if err != nil {
		return nil, err
	}
	if exists {
		sh, err := readSpriteHub(spriteHubPath)
		if err != nil {
			return nil, err
		}
		spriteHubs = append(spriteHubs, sh)
	}
	exists, err = fileutil.Exists(spriteHubGroupPath)
	if err != nil {
		return nil, err
	}
	if exists {
		shg, err := readSpriteHubGroup(spriteHubGroupPath)
		if err != nil {
			return nil, err
		}
		spriteHubs = append(spriteHubs, shg...)
	}
	return spriteHubs, nil
}
func readSpriteHub(spriteHubPath string) (spriteHub, error) {
	spriteHubJson, err := os.ReadFile(spriteHubPath)
	if err != nil {
		return spriteHub{}, err
	}
	var sh spriteHub
	err = json.Unmarshal(spriteHubJson, &sh)
	if err != nil {
		return spriteHub{}, err
	}
	return sh, nil
}
func readSpriteHubGroup(spriteHubGroupPath string) ([]spriteHub, error) {
	spriteHubGroupJson, err := os.ReadFile(spriteHubGroupPath)
	if err != nil {
		return nil, err
	}
	var shg struct {
		SpriteGroups []spriteHub `json:"spriteGroups"`
	}
	err = json.Unmarshal(spriteHubGroupJson, &shg)
	if err != nil {
		return nil, err
	}
	return shg.SpriteGroups, nil
}

func parseSpriteHub(name string, sh spriteHub, idm pathIDMap) CharAssetSpriteHub {
	cash := CharAssetSpriteHub{
		Sprites:            make([]Asset, len(sh.Sprites)),
		SpritesAlpha:       make([]Asset, len(sh.Sprites)),
		SpritesAlias:       make([]string, len(sh.Sprites)),
		SpritesIsWholeBody: make([]bool, len(sh.Sprites)),
	}
	cash.FacePos.X = sh.FacePos.X
	cash.FacePos.Y = sh.FacePos.Y
	cash.FaceSize.X = sh.FaceSize.X
	cash.FaceSize.Y = sh.FaceSize.Y
	for i, sprite := range sh.Sprites {
		cash.Sprites[i] = Asset{
			Name: fmt.Sprintf("%s/%s", name, idm[sprite.Sprite.MPathID]),
			Kind: KindCharacter,
		}
		cash.SpritesAlpha[i] = Asset{
			Name: fmt.Sprintf("%s/%s", name, idm[sprite.AlphaTex.MPathID]),
			Kind: KindCharacter,
		}
		cash.SpritesAlias[i] = sprite.Alias
		cash.SpritesIsWholeBody[i] = sprite.IsWholeBody != 0
	}
	return cash
}
