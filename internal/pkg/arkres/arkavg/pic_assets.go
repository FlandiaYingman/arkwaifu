package arkavg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func ScanForPicAssets(resDir string, prefix string) ([]Asset, error) {
	imgs, err := ScanForPicAssetsByKind(resDir, prefix, KindImage)
	if err != nil {
		return nil, err
	}
	bkgs, err := ScanForPicAssetsByKind(resDir, prefix, KindBackground)
	if err != nil {
		return nil, err
	}
	return lo.Flatten([][]Asset{imgs, bkgs}), nil
}

func ScanForPicAssetsByKind(resDir string, prefix string, kind Kind) ([]Asset, error) {
	path := filepath.Join(resDir, prefix, "avg", string(kind))
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	assets := make([]Asset, 0)
	for _, entry := range entries {
		// There could be some JSON files (asset metadata), so filter out any non-PNG files.
		if !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}
		assets = append(assets, Asset{
			Name: entry.Name(),
			Kind: kind,
		})
	}
	return assets, nil
}
