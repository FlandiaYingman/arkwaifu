package arkavg

import (
	"fmt"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type Asset struct {
	ID   string
	Kind Kind
}
type Kind string

const (
	KindImage      Kind = "images"
	KindBackground Kind = "backgrounds"
)

func GetAssets(resDir string, prefix string) ([]Asset, error) {
	var assets []Asset
	assetsImage, err := GetAssetsKind(resDir, prefix, KindImage)
	if err != nil {
		return nil, err
	}
	assets = append(assets, assetsImage...)
	assetsBackground, err := GetAssetsKind(resDir, prefix, KindBackground)
	if err != nil {
		return nil, err
	}
	assets = append(assets, assetsBackground...)
	return assets, nil
}

func GetAssetsKind(resDir string, prefix string, kind Kind) ([]Asset, error) {
	assetsPath := fmt.Sprintf("avg/%v", kind)
	assetsPath = filepath.Join(resDir, prefix, assetsPath)
	assetsDir, err := os.ReadDir(assetsPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	assets := make([]Asset, len(assetsDir))
	for i, it := range assetsDir {
		id := filepath.Base(it.Name())
		id = strings.TrimSuffix(id, filepath.Ext(id))
		assets[i] = Asset{
			ID:   id,
			Kind: kind,
		}
	}
	return assets, nil
}

func (a *Asset) Open(resDir string, prefix string) (*os.File, error) {
	assetPath := filepath.Join(resDir, prefix, "avg", string(a.Kind), pathutil.ReplaceExt(a.ID, ".png"))
	return os.Open(assetPath)
}
