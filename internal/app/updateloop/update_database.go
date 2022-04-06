package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkavg"
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

	err = c.avgService.UpdateAvg(ctx, resVer, avg)
	if err != nil {
		return err
	}

	return nil
}
func (c *Controller) updateAssetDatabase(ctx context.Context) error {
	err := c.assetService.InitNames(ctx, AcceptableAssetKinds, AcceptableAssetVariants)
	if err != nil {
		return err
	}
	err = c.assetService.ScanStaticDir(ctx)
	return err
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
		Name: strings.ToLower(a.Name),
		Kind: string(a.Kind),
	}
}
