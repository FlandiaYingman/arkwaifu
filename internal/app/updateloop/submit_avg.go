package updateloop

import (
	"context"
	"strings"

	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/arkavg"
)

// submitAvg submits the AVG data of the given resource version to the AVG service.
// The AVG data is parsed from the corresponding resource directory.
//
// Note that submitting AVG data is a fully overwrite operation.
func (s *Service) submitAvg(ctx context.Context, resVer ResVersion) error {
	resDir := s.ResDir(resVer)
	a, err := getAvg(resDir)
	if err != nil {
		return err
	}
	err = s.AvgService.UpdateAvg(ctx, string(resVer), a)
	if err != nil {
		return err
	}
	return nil
}

// getAvg simply wraps arkavg.GetAvg, and convert the arkavg models to the AVG service models.
func getAvg(resDir string) (avg.Avg, error) {
	r, err := arkavg.GetAvg(resDir, arkres.DefaultPrefix)
	if err != nil {
		return avg.Avg{}, err
	}
	a, err := avgFromRaw(&r, resDir)
	if err != nil {
		return avg.Avg{}, err
	}
	return a, nil
}

func avgFromRaw(ra *arkavg.Avg, resDir string) (avg.Avg, error) {
	var a = make(avg.Avg, 0, len(ra.Groups)) // note that a.Avg is []a.Group
	for _, rg := range ra.Groups {
		group, err := groupFromRaw(&rg, resDir)
		if err != nil {
			return nil, err
		}

		a = append(a, group)
	}
	return a, nil
}
func groupFromRaw(rg *arkavg.Group, resDir string) (avg.Group, error) {
	var stories = make([]avg.Story, 0, len(rg.Stories))
	for _, rs := range rg.Stories {
		story, err := storyFromRaw(&rs, resDir)
		if err != nil {
			return avg.Group{}, err
		}

		stories = append(stories, story)
	}
	return avg.Group{
		ID:      rg.ID,
		Name:    rg.Name,
		Type:    string(rg.Type),
		Stories: stories,
	}, nil
}
func storyFromRaw(rs *arkavg.Story, resDir string) (avg.Story, error) {
	rawAssets, err := arkavg.GetStoryAssets(resDir, arkres.DefaultPrefix, *rs)
	if err != nil {
		return avg.Story{}, err
	}

	var assets = make([]avg.Asset, 0, len(rawAssets))
	for _, ra := range rawAssets {
		asset := assetFromRaw(&ra)
		assets = append(assets, asset)
	}

	return avg.Story{
		ID:      rs.ID,
		Code:    rs.Code,
		Name:    rs.Name,
		Tag:     string(rs.Tag),
		GroupID: rs.GroupID,
		Assets:  assets,
	}, nil
}
func assetFromRaw(ra *arkavg.Asset) avg.Asset {
	return avg.Asset{
		Name: strings.ToLower(ra.Name),
		Kind: string(ra.Kind),
	}
}
