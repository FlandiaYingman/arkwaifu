package asset

import (
	"context"
	"path/filepath"

	"github.com/flandiayingman/arkwaifu/internal/pkg/util/fileutil"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
)

func (s *Service) PopulateFrom(ctx context.Context, dirPath string) error {
	err := fileutil.CopyAllContent(dirPath, s.staticDir)
	if err != nil {
		return errors.Wrapf(err, "failed to copy all content from %s to static dir %s", dirPath, s.staticDir)
	}
	err = fileutil.LowercaseAll(s.staticDir)
	if err != nil {
		return errors.Wrapf(err, "failed to lowercase all files in static dir %s", s.staticDir)
	}

	am, vm, err := scanStaticDir(s.staticDir)
	if err != nil {
		return errors.Wrapf(err, "failed to scan static dir %s", s.staticDir)
	}
	return s.repo.Update(ctx, am, vm)
}

func scanStaticDir(staticDir string) ([]modelAsset, []modelVariant, error) {
	scanPattern := filepath.Join(staticDir, "*/*/*")
	scanFiles, err := filepath.Glob(scanPattern)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to scan the static dir %s by pattern %s", staticDir, scanPattern)
	}

	assetMap := make(map[string]map[string]modelAsset)
	variantSlice := make([]modelVariant, 0, len(scanFiles))

	for _, scanFile := range scanFiles {
		scanFile, err = filepath.Rel(staticDir, scanFile)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to get relative path of %s and base %s", scanFile, staticDir)
		}
		scanFile, filename := filepath.Split(scanFile)
		scanFile, kindName := filepath.Split(filepath.Dir(scanFile))
		scanFile, variantName := filepath.Split(filepath.Dir(scanFile))
		name := pathutil.RemoveAllExt(filename)

		if _, ok := assetMap[kindName]; !ok {
			assetMap[kindName] = make(map[string]modelAsset)
		}
		if _, ok := assetMap[kindName][name]; !ok {
			assetMap[kindName][name] = modelAsset{
				Kind: kindName,
				Name: name,
			}
		}
		variantSlice = append(variantSlice, modelVariant{
			AssetKind: kindName,
			AssetName: name,
			Variant:   variantName,
			Filename:  filename,
		})
	}

	assetSlice := make([]modelAsset, 0, len(assetMap))
	for _, m := range assetMap {
		for _, asset := range m {
			assetSlice = append(assetSlice, asset)
		}
	}

	return assetSlice, variantSlice, nil
}
