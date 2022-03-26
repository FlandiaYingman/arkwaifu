package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/asset"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkavg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func (c *Controller) updateDatabase(ctx context.Context, resVer string, resDir string) error {
	rawAvg, err := arkavg.GetAvg(resDir, arkres.DefaultPrefix)
	if err != nil {
		return err
	}

	avg, err := avgFromRaw(&rawAvg, resDir)
	if err != nil {
		return err
	}

	err = c.avgService.SetAvgs(resVer, avg)
	if err != nil {
		return err
	}

	return nil
}
func (c *Controller) updateAssetDatabase(ctx context.Context, mainStaticDir string) error {
	assets, err := scanMainStaticDir(mainStaticDir)
	if err != nil {
		return err
	}
	err = c.assetService.SetAssets(ctx, assets)
	if err != nil {
		return err
	}
	return nil
}

func scanMainStaticDir(mainStaticDir string) ([]asset.Asset, error) {
	variantDirEntries, err := os.ReadDir(mainStaticDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read the main static dir %s", mainStaticDir)
	}

	var variants []asset.Asset
	for _, variantDirEntry := range variantDirEntries {
		if variantDirEntry.IsDir() {
			variantDir := filepath.Join(mainStaticDir, variantDirEntry.Name())
			err := scanVariantDir(&variants, variantDir)
			if err != nil {
				return nil, err
			}
		}
	}

	return variants, nil
}
func scanVariantDir(variants *[]asset.Asset, variantDir string) error {
	kindDirEntries, err := os.ReadDir(variantDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read the variant dir %s", variantDir)
	}

	for _, kindDirEntry := range kindDirEntries {
		if kindDirEntry.IsDir() {
			kindDir := filepath.Join(variantDir, kindDirEntry.Name())
			err := scanKindDir(variants, variantDir, kindDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func scanKindDir(variants *[]asset.Asset, variantDir, kindDir string) error {
	assetFiles, err := os.ReadDir(kindDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read the kind dir %s", kindDir)
	}

	kindName := filepath.Base(kindDir)
	variantName := filepath.Base(variantDir)

	for _, assetFile := range assetFiles {
		if assetFile.IsDir() {
			continue
		}
		assetFileName := assetFile.Name()
		assetName := pathutil.RemoveExt(assetFileName)
		*variants = append(*variants, asset.Asset{
			Kind:     kindName,
			Asset:    assetName,
			Variant:  variantName,
			FileName: assetFileName,
		})
	}
	return nil
}

func avgFromRaw(a *arkavg.Avg, resDir string) (avg.Avg, error) {
	var groups []avg.Group
	for _, rawGroup := range a.Groups {
		group, err := groupFromRaw(&rawGroup, resDir)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
func groupFromRaw(g *arkavg.Group, resDir string) (avg.Group, error) {
	var stories []avg.Story
	for _, rawStory := range g.Stories {
		story, err := storyFromRaw(&rawStory, resDir)
		if err != nil {
			return avg.Group{}, err
		}
		stories = append(stories, story)
	}
	return avg.Group{
		ID:      g.ID,
		Name:    g.Name,
		Type:    string(g.Type),
		Stories: stories,
	}, nil
}
func storyFromRaw(s *arkavg.Story, resDir string) (avg.Story, error) {
	rawAssets, err := arkavg.GetStoryAssets(resDir, arkres.DefaultPrefix, *s)
	if err != nil {
		return avg.Story{}, nil
	}

	var assetModels []avg.Asset
	for _, rawAsset := range rawAssets {
		assetModel := assetFromRaw(&rawAsset)
		assetModels = append(assetModels, assetModel)
	}

	return avg.Story{
		ID:      s.ID,
		Code:    s.Code,
		Name:    s.Name,
		Tag:     string(s.Tag),
		GroupID: s.GroupID,
		Assets:  assetModels,
	}, nil
}
func assetFromRaw(a *arkavg.Asset) avg.Asset {
	return avg.Asset{
		ID:   strings.ToLower(a.ID),
		Kind: string(a.Kind),
	}
}
